package handler_test

import (
	"bytes"
	"log/slog"
	"sync"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb"
	"github.com/ntchjb/usbip-virtual-device/usbip/handler"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var urbQueueCmdSubmits = []protocol.CmdSubmit{
	{
		CmdHeader: protocol.CmdHeader{
			Command:        protocol.CMD_SUBMIT,
			SeqNum:         1,
			DevID:          0x00010001,
			Direction:      protocol.DIR_OUT,
			EndpointNumber: 1,
		},
		TransferFlags:        0,
		TransferBufferLength: 5,
		NumberOfPackets:      0xffffffff,
		Interval:             0x00000064,
		TransferBuffer: []byte{
			0x01, 0x02, 0x03, 0x04, 0x05,
		},
	},
	{
		CmdHeader: protocol.CmdHeader{
			Command:        protocol.CMD_SUBMIT,
			SeqNum:         2,
			DevID:          0x00010001,
			Direction:      protocol.DIR_OUT,
			EndpointNumber: 1,
		},
		TransferFlags:        0,
		TransferBufferLength: 5,
		NumberOfPackets:      0xffffffff,
		Interval:             0x00000064,
		TransferBuffer: []byte{
			0x01, 0x02, 0x03, 0x04, 0x05,
		},
	},
	{
		CmdHeader: protocol.CmdHeader{
			Command:        protocol.CMD_SUBMIT,
			SeqNum:         3,
			DevID:          0x00010001,
			Direction:      protocol.DIR_OUT,
			EndpointNumber: 1,
		},
		TransferFlags:        0,
		TransferBufferLength: 5,
		NumberOfPackets:      0xffffffff,
		Interval:             0x00000064,
		TransferBuffer: []byte{
			0x01, 0x02, 0x03, 0x04, 0x05,
		},
	},
}

var urbQueueRetSubmits = []protocol.RetSubmit{
	{
		CmdHeader: protocol.CmdHeader{
			Command:        protocol.CMD_SUBMIT,
			SeqNum:         1,
			DevID:          0x00010001,
			Direction:      protocol.DIR_OUT,
			EndpointNumber: 1,
		},
		Status:          0x00000001,
		ActualLength:    5,
		StartFrame:      0,
		NumberOfPackets: 0xffffffff,
		ErrorCount:      0,
		Padding:         0,
		TransferBuffer:  []byte{0x01, 0x02, 0x03, 0x04, 0x05},
	},
	{
		CmdHeader: protocol.CmdHeader{
			Command:        protocol.CMD_SUBMIT,
			SeqNum:         2,
			DevID:          0x00010001,
			Direction:      protocol.DIR_OUT,
			EndpointNumber: 1,
		},
		Status:          0x00000001,
		ActualLength:    5,
		StartFrame:      0,
		NumberOfPackets: 0xffffffff,
		ErrorCount:      0,
		Padding:         0,
		TransferBuffer:  []byte{0x01, 0x02, 0x03, 0x04, 0x05},
	},
	{
		CmdHeader: protocol.CmdHeader{
			Command:        protocol.CMD_SUBMIT,
			SeqNum:         3,
			DevID:          0x00010001,
			Direction:      protocol.DIR_OUT,
			EndpointNumber: 1,
		},
		Status:          0x00000001,
		ActualLength:    5,
		StartFrame:      0,
		NumberOfPackets: 0xffffffff,
		ErrorCount:      0,
		Padding:         0,
		TransferBuffer:  []byte{0x01, 0x02, 0x03, 0x04, 0x05},
	},
}

type Buffer struct {
	b bytes.Buffer
	m sync.Mutex
}

func (b *Buffer) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Read(p)
}
func (b *Buffer) Write(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Write(p)
}
func (b *Buffer) Bytes() []byte {
	b.m.Lock()
	defer b.m.Unlock()
	return b.b.Bytes()
}

func TestWorkerPool(t *testing.T) {
	ctrl := gomock.NewController(t)

	replies := new(Buffer)
	logger := slog.Default()
	wp := handler.NewWorkerPool(replies, logger)
	device := usb.NewMockDevice(ctrl)

	device.EXPECT().GetWorkerPoolProfile().Return(usb.WorkerPoolProfile{
		MaximumProcWorkers:  1,
		MaximumReplyWorkers: 1,
	})
	device.EXPECT().Process(urbQueueCmdSubmits[0]).Return(urbQueueRetSubmits[0]).Times(1)
	device.EXPECT().Process(urbQueueCmdSubmits[1]).Return(urbQueueRetSubmits[1]).AnyTimes()
	device.EXPECT().Process(urbQueueCmdSubmits[2]).Return(urbQueueRetSubmits[2]).Times(1)

	wp.SetDevice(device)
	startErr := wp.Start()

	assert.NoError(t, startErr)

	wp.PublishCmdSubmit(urbQueueCmdSubmits[0])
	wp.PublishCmdSubmit(urbQueueCmdSubmits[1])
	errUnlink := wp.Unlink(protocol.CmdUnlink{
		CmdHeader: protocol.CmdHeader{
			Command:        protocol.CMD_UNLINK,
			SeqNum:         4,
			DevID:          0x00010001,
			Direction:      protocol.DIR_OUT,
			EndpointNumber: 1,
		},
		UnlinkSeqNum: 2,
	})
	wp.PublishCmdSubmit(urbQueueCmdSubmits[2])

	assert.NoError(t, errUnlink)

	wp.Stop()

	assert.Equal(t, []byte{
		// protocol.RetUnlink
		0x00, 0x00, 0x00, 0x04,
		0x00, 0x00, 0x00, 0x04,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xff, 0xff, 0xff, 0x98, // ECONNRESET
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,

		// protocol.RetSubmit
		0x00, 0x00, 0x00, 0x01, // Command
		0x00, 0x00, 0x00, 0x01, // SeqNum
		0x00, 0x01, 0x00, 0x01, // DevID
		0x00, 0x00, 0x00, 0x00, // Direction
		0x00, 0x00, 0x00, 0x01, // EndpointNumber

		0x00, 0x00, 0x00, 0x01, // Status
		0x00, 0x00, 0x00, 0x05, // ActualLength
		0x00, 0x00, 0x00, 0x00, // StartFrame
		0xff, 0xff, 0xff, 0xff, // NumberOfPackets
		0x00, 0x00, 0x00, 0x00, // ErrorCount
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Padding
		0x01, 0x02, 0x03, 0x04, 0x05, // TransferBuffer

		// protocol.RetSubmit for SeqNum 2 should not be here because it's unlinked

		// protocol.RetSubmit
		0x00, 0x00, 0x00, 0x01, // Command
		0x00, 0x00, 0x00, 0x03, // SeqNum
		0x00, 0x01, 0x00, 0x01, // DevID
		0x00, 0x00, 0x00, 0x00, // Direction
		0x00, 0x00, 0x00, 0x01, // EndpointNumber

		0x00, 0x00, 0x00, 0x01, // Status
		0x00, 0x00, 0x00, 0x05, // ActualLength
		0x00, 0x00, 0x00, 0x00, // StartFrame
		0xff, 0xff, 0xff, 0xff, // NumberOfPackets
		0x00, 0x00, 0x00, 0x00, // ErrorCount
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // Padding
		0x01, 0x02, 0x03, 0x04, 0x05, // TransferBuffer
	}, replies.Bytes())
}
