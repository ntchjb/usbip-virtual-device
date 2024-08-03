package handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"syscall"

	"github.com/ntchjb/usbip-virtual-device/usb"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
)

// WorkerPool is a pool of workers processing CmdSubmit requests and reply with RetSubmit to client
type WorkerPool interface {
	// Start worker pool
	Start() error
	// Stop worker pool
	Stop() error
	// Set URB processor after selected a USB device via OpImport operation
	SetProcessor(processor usb.URBProcessor)

	// Delegated functions from Queue
	// Mark URB as unlinked by given sequence number
	Unlink(header protocol.CmdUnlink) error
	// Publish CmdSubmit to worker pool to be further processed
	PublishCmdSubmit(urb protocol.CmdSubmit)
}

type WorkerPoolConfig struct {
	MaximumProcWorkers  int
	MaximumReplyWorkers int
}

type workerPoolImpl struct {
	config      WorkerPoolConfig
	wg          sync.WaitGroup
	queue       Queue
	logger      *slog.Logger
	processor   usb.URBProcessor
	replyWriter io.Writer

	processingURBsLock sync.RWMutex
	processingURBs     map[uint32]bool
	unlinkFlagsLock    sync.RWMutex
	unlinkFlags        map[uint32]bool
}

func NewWorkerPool(config WorkerPoolConfig, queue Queue, replyWriter io.Writer, logger *slog.Logger) WorkerPool {
	return &workerPoolImpl{
		queue:          queue,
		config:         config,
		logger:         logger,
		replyWriter:    replyWriter,
		processingURBs: make(map[uint32]bool),
		unlinkFlags:    make(map[uint32]bool),
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

func (p *workerPoolImpl) markAsDone(seqNum uint32) {
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
	if err := retUnlink.Encode(p.replyWriter); err != nil {
		return fmt.Errorf("unable to encode RetUnlink: %w", err)
	}

	return nil
}

func (p *workerPoolImpl) PublishCmdSubmit(urb protocol.CmdSubmit) {
	p.markAsProcessing(urb.SeqNum)
	p.queue.PublishCmdSubmit(urb)
}

func (p *workerPoolImpl) SetProcessor(processor usb.URBProcessor) {
	p.processor = processor
}

func (p *workerPoolImpl) Start() error {
	// Initiate worker pool for processing CmdSubmit from queue
	for i := 0; i < p.config.MaximumProcWorkers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()

			for {
				urbSubmit, err := p.queue.ConsumeCmdSubmit()
				if err != nil {
					if !errors.Is(err, ErrQueueClosed) {
						p.logger.Error("unable to dequeue CmdSubmit", "err", err)
					}
					break
				}

				// Check Unlink #1: before return to caller, for processing
				if p.removeUnlinkFlag(urbSubmit.SeqNum) {
					p.logger.Debug("Unlinked URB detected, ignoring", "urbSeqNum", urbSubmit.SeqNum)
					p.markAsDone(urbSubmit.SeqNum)
					continue
				}

				urbRet := p.processor.Process(urbSubmit)
				p.queue.PublishRetSubmit(urbRet)
			}
		}()
	}

	// Initiate worker pool for sending RetSubmit to io.Writer (which should be net.Conn)
	for i := 0; i < p.config.MaximumReplyWorkers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()

			for {
				urbRet, err := p.queue.ConsumeRetSubmit()
				if err != nil {
					if !errors.Is(err, ErrQueueClosed) {
						p.logger.Error("unable to dequeue RetSubmit", "err", err)
					}
					break
				}

				// After markAsDone is called, further logic will ignore any upcoming unlink requests
				// for this URB and the server will reply RetSubmit back to client
				p.markAsDone(urbRet.SeqNum)

				// Check Unlink #2: before return to caller, for sending to client
				if p.removeUnlinkFlag(urbRet.SeqNum) {
					p.logger.Debug("Unlinked URB detected, ignoring", "urbSeqNum", urbRet.SeqNum)
					p.markAsDone(urbRet.SeqNum)
					continue
				}

				if err := urbRet.Encode(p.replyWriter); err != nil {
					p.logger.Error("unable to encode RetSubmit to stream", "err", err)
				}
			}
		}()
	}
	return nil
}

func (p *workerPoolImpl) Stop() error {
	if err := p.queue.Close(); err != nil {
		return fmt.Errorf("unable to close URB queue: %w", err)
	}
	p.wg.Wait()

	return nil
}
