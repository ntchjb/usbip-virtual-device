package handler

import (
	"errors"

	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
)

const (
	QUEUE_SIZE = 1024
)

var (
	ErrAlreadyUnlinked = errors.New("URB was already unlinked")
	ErrQueueClosed     = errors.New("queue is closed")
)

type Queue interface {
	// Get CmdSubmit from queue, used by worker pool
	ConsumeCmdSubmit() (protocol.CmdSubmit, error)
	// Publish RetSubmit to queue, used by worker pool after processed CmdSubmit
	PublishRetSubmit(urb protocol.RetSubmit)
	// Get RetSubmit from queue, used by worker pool, preparing for sending to io.Writer (net.Conn)
	ConsumeRetSubmit() (protocol.RetSubmit, error)
	// Publish CmdSubmit to worker pool to be further processed
	PublishCmdSubmit(urb protocol.CmdSubmit)
	// Close the queue
	Close() error
}

type queueImpl struct {
	cmdQueue chan protocol.CmdSubmit
	retQueue chan protocol.RetSubmit
}

func NewURBQueue() Queue {
	return &queueImpl{
		cmdQueue: make(chan protocol.CmdSubmit, QUEUE_SIZE),
		retQueue: make(chan protocol.RetSubmit, QUEUE_SIZE),
	}
}

// QueueSubimt adds CmdSubmit to queue, waiting for data processing
func (u *queueImpl) PublishCmdSubmit(urb protocol.CmdSubmit) {
	u.cmdQueue <- urb
}

// ConsumeCmdSubmit return CmdSubmit from queue to be processed,
// returns ErrAlreadyUnlinked if URB is already unlinked
func (u *queueImpl) ConsumeCmdSubmit() (protocol.CmdSubmit, error) {
	urb, ok := <-u.cmdQueue
	if !ok {
		return urb, ErrQueueClosed
	}

	return urb, nil
}

// PublishRetSubmit adds RetSubmit to queue, waiting for sending to client
func (u *queueImpl) PublishRetSubmit(urb protocol.RetSubmit) {
	u.retQueue <- urb
}

// ConsumeRetSubmit returns RetSubmit (if not unlinked) from queue to be sent to client,
// returns ErrAlreadyUnlinked if URB is already unlinked
func (u *queueImpl) ConsumeRetSubmit() (protocol.RetSubmit, error) {
	urb, ok := <-u.retQueue
	if !ok {
		return urb, ErrQueueClosed
	}

	return urb, nil
}

func (u *queueImpl) Close() error {
	close(u.cmdQueue)
	close(u.retQueue)

	return nil
}
