package handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"

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
	Publisher
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
}

func NewWorkerPool(config WorkerPoolConfig, queue Queue, replyWriter io.Writer, logger *slog.Logger) WorkerPool {
	return &workerPoolImpl{
		queue:       queue,
		config:      config,
		logger:      logger,
		replyWriter: replyWriter,
	}
}

func (p *workerPoolImpl) Unlink(seqNum uint32) {
	p.queue.Unlink(seqNum)
}

func (p *workerPoolImpl) PublishCmdSubmit(urb protocol.CmdSubmit) {
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
					if errors.Is(err, ErrQueueClosed) {
						break
					} else if errors.Is(err, ErrAlreadyUnlinked) {
						p.logger.Debug("Unlinked URB detected, ignoring", "urbSeqNum", urbSubmit.SeqNum)
						continue
					} else {
						p.logger.Error("unable to dequeue CmdSubmit", "err", err)
						continue
					}
				}

				urbRet := p.processor.ProcessSubmit(urbSubmit)
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
					if errors.Is(err, ErrQueueClosed) {
						break
					} else if errors.Is(err, ErrAlreadyUnlinked) {
						p.logger.Debug("Unlinked URB detected, ignoring", "urbSeqNum", urbRet.SeqNum)
						continue
					} else {
						p.logger.Error("unable to dequeue CmdSubmit", "err", err)
						continue
					}
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
		return fmt.Errorf("unable to close UDB queue: %w", err)
	}
	p.wg.Wait()

	return nil
}
