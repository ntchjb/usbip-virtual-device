package handler

import (
	"fmt"
	"io"
	"log/slog"
	"sync"
	"syscall"

	"github.com/ntchjb/usbip-virtual-device/usb"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
)

const (
	URB_QUEUE_SIZE = 1024

	URB_STATUS_UNLINKING  uint8 = 0
	URB_STATUS_PROCESSING uint8 = 1
	URB_STATUS_REPLYING   uint8 = 2
)

// WorkerPool is a pool of workers processing CmdSubmit requests and reply with RetSubmit to client
type WorkerPool interface {
	// Start worker pool
	Start() error
	// Stop worker pool
	Stop() error
	// SetDevice assignes device to worker pool as data receiver and processor
	SetDevice(device usb.Device)
	// Delegated functions from Queue
	// Mark URB as unlinked by given sequence number
	Unlink(header protocol.CmdUnlink) error
	// Publish CmdSubmit to worker pool to be further processed
	PublishCmdSubmit(urb protocol.CmdSubmit)
}

type workerPoolImpl struct {
	wgCmdSubmit sync.WaitGroup
	wgRetSubmit sync.WaitGroup
	logger      *slog.Logger
	device      usb.Device
	replyWriter io.Writer
	conf        usb.WorkerPoolProfile

	cmdQueue    chan protocol.CmdSubmit
	retQueue    chan protocol.RetSubmit
	unlinkQueue chan protocol.RetUnlink

	processingURBsLock sync.RWMutex
	processingURBs     map[uint32]uint8
}

func NewWorkerPool(replyWriter io.Writer, logger *slog.Logger) WorkerPool {
	return &workerPoolImpl{
		logger:         logger,
		replyWriter:    replyWriter,
		processingURBs: make(map[uint32]uint8),
		cmdQueue:       make(chan protocol.CmdSubmit, URB_QUEUE_SIZE),
		retQueue:       make(chan protocol.RetSubmit, URB_QUEUE_SIZE),
		unlinkQueue:    make(chan protocol.RetUnlink, URB_QUEUE_SIZE),
	}
}

func (p *workerPoolImpl) markAsProcessing(seqNum uint32) bool {
	p.processingURBsLock.Lock()
	defer p.processingURBsLock.Unlock()

	if _, ok := p.processingURBs[seqNum]; !ok {
		p.processingURBs[seqNum] = URB_STATUS_PROCESSING
		return true
	} else {
		return false
	}
}

func (p *workerPoolImpl) markAsUnlink(seqNum uint32) bool {
	p.processingURBsLock.Lock()
	defer p.processingURBsLock.Unlock()

	if _, ok := p.processingURBs[seqNum]; ok {
		p.processingURBs[seqNum] = URB_STATUS_UNLINKING
		return true
	} else {
		return false
	}
}

// markAsReplying marks URB as processed, waiting for replies.
// This function deletes status only if its current status is not PROCESSING.
// That means URBs that got unlinked before marking as REPLYING will be ignored
func (p *workerPoolImpl) markAsReplying(seqNum uint32) bool {
	var res bool
	p.processingURBsLock.Lock()
	defer p.processingURBsLock.Unlock()

	if status, ok := p.processingURBs[seqNum]; ok && status == URB_STATUS_PROCESSING {
		p.processingURBs[seqNum] = URB_STATUS_REPLYING
		res = true
	} else {
		delete(p.processingURBs, seqNum)
		res = false
	}

	return res
}

// markAsReplying marks URB as processed, waiting for replies.
// This function deletes status once it's called.
func (p *workerPoolImpl) markAsReplied(seqNum uint32) bool {
	var res bool
	p.processingURBsLock.Lock()
	defer p.processingURBsLock.Unlock()

	if status, ok := p.processingURBs[seqNum]; ok && status == URB_STATUS_REPLYING {
		res = true
	} else {
		res = false
	}
	delete(p.processingURBs, seqNum)

	return res
}

func (p *workerPoolImpl) Unlink(cmd protocol.CmdUnlink) error {
	p.logger.Debug("Unlink request received", "data", cmd)
	retUnlink := protocol.RetUnlink{
		CmdHeader: protocol.CmdHeader{
			Command: protocol.RET_UNLINK,
			SeqNum:  cmd.SeqNum,
		},
		Status: 0,
	}

	// If given URB is currently processing by worker pool, then RetUnlink should
	// return with status -ECONNRESET.
	// Otherwise, return with status 0
	if p.markAsUnlink(cmd.UnlinkSeqNum) {
		retUnlink.Status = -int32(syscall.ECONNRESET)
	} else {
		p.logger.Debug("Unlink is ignored, does not receive CmdSubmit yet", "seqNum", cmd.SeqNum, "unlinkSeqNum", cmd.UnlinkSeqNum)
	}

	p.unlinkQueue <- retUnlink

	return nil
}

func (p *workerPoolImpl) PublishCmdSubmit(urb protocol.CmdSubmit) {
	p.logger.Debug("Received CmdSubmit", "data", urb)
	if !p.markAsProcessing(urb.SeqNum) {
		p.logger.Error("Found duplicated URB, ignoring", "urb", urb)
		return
	}
	p.cmdQueue <- urb
}

func (p *workerPoolImpl) SetDevice(device usb.Device) {
	p.device = device
}

func (p *workerPoolImpl) Start() error {
	if p.device == nil {
		return fmt.Errorf("device does not exist in this worker pool")
	}
	p.conf = p.device.GetWorkerPoolProfile()
	// Initiate worker pool for processing CmdSubmit from queue
	for i := 0; i < p.conf.MaximumProcWorkers; i++ {
		p.wgCmdSubmit.Add(1)
		go func() {
			defer p.wgCmdSubmit.Done()

			for urbSubmit := range p.cmdQueue {
				if !p.markAsReplying(urbSubmit.SeqNum) {
					p.logger.Debug("Unlinked URB detected before processing it, ignoring", "urbSeqNum", urbSubmit.SeqNum)
					continue
				}

				urbRet := p.device.Process(urbSubmit)
				p.retQueue <- urbRet
			}
		}()
	}

	// Initiate worker pool for sending RetSubmit to io.Writer (which should be net.Conn)
	for i := 0; i < p.conf.MaximumReplyWorkers; i++ {
		p.wgRetSubmit.Add(1)
		go func() {
			defer p.wgRetSubmit.Done()

			for urbRet := range p.retQueue {
				if !p.markAsReplied(urbRet.SeqNum) {
					p.logger.Debug("Unlinked URB detected, ignoring", "urbSeqNum", urbRet.SeqNum)
					continue
				}

				p.logger.Debug("Replying RetSubmit", "data", urbRet)
				if err := urbRet.CmdHeader.Encode(p.replyWriter); err != nil {
					p.logger.Error("unable to encode RetSubmit header to stream", "err", err, "seqNum", urbRet.SeqNum)
				} else if err := urbRet.Encode(p.replyWriter); err != nil {
					p.logger.Error("unable to encode RetSubmit to stream", "err", err, "seqNum", urbRet.SeqNum)
				}
			}
		}()
	}

	// Initialize worker pool for sending RetUnlink to io.Writer (which should be net.Conn)
	for i := 0; i < p.conf.MaximumUnlinkReplyWorkers; i++ {
		p.wgRetSubmit.Add(1)
		go func() {
			defer p.wgRetSubmit.Done()

			for urbRet := range p.unlinkQueue {
				p.logger.Debug("Replying RetUnlink", "data", urbRet)
				if err := urbRet.CmdHeader.Encode(p.replyWriter); err != nil {
					p.logger.Error("unable to encode RetUnlink header to stream", "err", err, "seqNum", urbRet.SeqNum)
				} else if err := urbRet.Encode(p.replyWriter); err != nil {
					p.logger.Error("unable to encode RetUnlink to stream", "err", err, "seqNum", urbRet.SeqNum)
				}
				p.logger.Debug("Unlink Replied", "urb", urbRet)
			}
		}()
	}

	return nil
}

func (p *workerPoolImpl) Stop() error {
	close(p.cmdQueue)
	p.wgCmdSubmit.Wait()
	close(p.retQueue)
	close(p.unlinkQueue)
	p.wgRetSubmit.Wait()

	p.conf = usb.WorkerPoolProfile{}

	return nil
}
