package handler

import (
	"errors"
	"sync"

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
	Publisher
	// Get CmdSubmit from queue, used by worker pool
	ConsumeCmdSubmit() (protocol.CmdSubmit, error)
	// Publish RetSubmit to queue, used by worker pool after processed CmdSubmit
	PublishRetSubmit(urb protocol.RetSubmit)
	// Get RetSubmit from queue, used by worker pool, preparing for sending to io.Writer (net.Conn)
	ConsumeRetSubmit() (protocol.RetSubmit, error)
	// Close the queue
	Close() error
}

// Publisher is queue functions used by request handler, but delegated by worker
type Publisher interface {
	// Mark URB as unlinked by given sequence number
	Unlink(seqNum uint32)
	// Publish CmdSubmit to worker pool to be further processed
	PublishCmdSubmit(urb protocol.CmdSubmit)
}

type queueImpl struct {
	cmdQueue        chan protocol.CmdSubmit
	retQueue        chan protocol.RetSubmit
	unlinkFlags     map[uint32]bool
	unlinkFlagsLock sync.RWMutex
}

func NewURBQueue() Queue {
	return &queueImpl{
		cmdQueue:    make(chan protocol.CmdSubmit, QUEUE_SIZE),
		retQueue:    make(chan protocol.RetSubmit, QUEUE_SIZE),
		unlinkFlags: make(map[uint32]bool),
	}
}

// Unlink marks URB with given sequence number as "unlinked"
func (u *queueImpl) Unlink(seqNum uint32) {
	u.unlinkFlagsLock.Lock()
	defer u.unlinkFlagsLock.Unlock()

	u.unlinkFlags[seqNum] = true
}

func (u *queueImpl) popUnlinkFlag(seqNum uint32) bool {
	u.unlinkFlagsLock.Lock()
	defer u.unlinkFlagsLock.Unlock()

	_, ok := u.unlinkFlags[seqNum]
	delete(u.unlinkFlags, seqNum)

	return ok
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

	// Check Unlink #1: before return to caller, for processing
	if u.popUnlinkFlag(urb.SeqNum) {
		return urb, ErrAlreadyUnlinked
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

	// Check Unlink #2: before return to caller, for sending to client
	if u.popUnlinkFlag(urb.SeqNum) {
		return urb, ErrAlreadyUnlinked
	}

	return urb, nil
}

func (u *queueImpl) Close() error {
	close(u.cmdQueue)
	close(u.retQueue)

	return nil
}
