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

	cmdQueue chan protocol.CmdSubmit
	retQueue chan protocol.RetSubmit

	processingURBsLock sync.RWMutex
	processingURBs     map[uint32]bool
	unlinkFlagsLock    sync.RWMutex
	unlinkFlags        map[uint32]bool
}

func NewWorkerPool(replyWriter io.Writer, logger *slog.Logger) WorkerPool {
	return &workerPoolImpl{
		logger:         logger,
		replyWriter:    replyWriter,
		processingURBs: make(map[uint32]bool),
		unlinkFlags:    make(map[uint32]bool),
		cmdQueue:       make(chan protocol.CmdSubmit, URB_QUEUE_SIZE),
		retQueue:       make(chan protocol.RetSubmit, URB_QUEUE_SIZE),
	}
}

func (p *workerPoolImpl) removeUnlinkFlag(seqNum uint32) bool {
	p.unlinkFlagsLock.Lock()
	defer p.unlinkFlagsLock.Unlock()

	_, ok := p.unlinkFlags[seqNum]
	delete(p.unlinkFlags, seqNum)

	return ok
}

func (p *workerPoolImpl) addUnlinkFlag(seqNum uint32) {
	p.unlinkFlagsLock.Lock()
	defer p.unlinkFlagsLock.Unlock()

	p.unlinkFlags[seqNum] = true
}

func (p *workerPoolImpl) isURBProcessing(seqNum uint32) bool {
	p.processingURBsLock.RLock()
	defer p.processingURBsLock.RUnlock()

	_, ok := p.processingURBs[seqNum]

	return ok
}

func (p *workerPoolImpl) markAsProcessing(seqNum uint32) {
	p.processingURBsLock.Lock()
	defer p.processingURBsLock.Unlock()

	p.processingURBs[seqNum] = true
}

func (p *workerPoolImpl) markAsProcessed(seqNum uint32) {
	p.processingURBsLock.Lock()
	defer p.processingURBsLock.Unlock()

	delete(p.processingURBs, seqNum)
}

func (p *workerPoolImpl) Unlink(cmd protocol.CmdUnlink) error {
	retUnlink := protocol.RetUnlink{
		CmdHeader: cmd.CmdHeader,
		Status:    -int32(syscall.ENOENT),
	}
	retUnlink.Command = protocol.RET_UNLINK

	// If given URB is currently processing by worker pool, then RetUnlink should
	// return with status -ECONNRESET.
	// Otherwise, return with status -ENOENT
	if p.isURBProcessing(cmd.UnlinkSeqNum) {
		retUnlink.Status = -int32(syscall.ECONNRESET)

		p.addUnlinkFlag(cmd.UnlinkSeqNum)
	}

	// Reply RetUnlink back to client
	if err := retUnlink.CmdHeader.Encode(p.replyWriter); err != nil {
		return fmt.Errorf("unable to encode RetUnlink header: %w", err)
	} else if err := retUnlink.Encode(p.replyWriter); err != nil {
		return fmt.Errorf("unable to encode RetUnlink: %w", err)
	}

	return nil
}

func (p *workerPoolImpl) PublishCmdSubmit(urb protocol.CmdSubmit) {
	p.logger.Debug("Received CmdSubmit", "data", urb)
	p.markAsProcessing(urb.SeqNum)
	p.cmdQueue <- urb
}

func (p *workerPoolImpl) SetDevice(device usb.Device) {
	p.device = device
}

func (p *workerPoolImpl) Start() error {
	if p.device == nil {
		return fmt.Errorf("device does not exist in this worker pool")
	}
	config := p.device.GetWorkerPoolProfile()
	// Initiate worker pool for processing CmdSubmit from queue
	for i := 0; i < config.MaximumProcWorkers; i++ {
		p.wgCmdSubmit.Add(1)
		go func() {
			defer p.wgCmdSubmit.Done()

			for urbSubmit := range p.cmdQueue {
				// Check Unlink #1: before return to caller, for processing
				if p.removeUnlinkFlag(urbSubmit.SeqNum) {
					p.logger.Debug("Unlinked URB detected, ignoring", "urbSeqNum", urbSubmit.SeqNum)
					p.markAsProcessed(urbSubmit.SeqNum)
					continue
				}
				// After this is called, further logic will ignore any upcoming unlink requests
				// for this URB and the server will reply RetSubmit back to client
				p.markAsProcessed(urbSubmit.SeqNum)

				urbRet := p.device.Process(urbSubmit)
				p.retQueue <- urbRet
			}
		}()
	}

	// Initiate worker pool for sending RetSubmit to io.Writer (which should be net.Conn)
	for i := 0; i < config.MaximumReplyWorkers; i++ {
		p.wgRetSubmit.Add(1)
		go func() {
			defer p.wgRetSubmit.Done()

			for urbRet := range p.retQueue {
				p.logger.Debug("Replying RetSubmit", "data", urbRet)
				if err := urbRet.CmdHeader.Encode(p.replyWriter); err != nil {
					p.logger.Error("unable to encode RetSubmit header to stream", "err", err)
				} else if err := urbRet.Encode(p.replyWriter); err != nil {
					p.logger.Error("unable to encode RetSubmit to stream", "err", err)
				}
			}
		}()
	}
	return nil
}

func (p *workerPoolImpl) Stop() error {
	close(p.cmdQueue)
	p.wgCmdSubmit.Wait()
	close(p.retQueue)
	p.wgRetSubmit.Wait()

	return nil
}
