package protocol_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol"
	usbipprot "github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
)

func TestSetupPacketSerializer(t *testing.T) {
	tests := []struct {
		name   string
		obj    usbipprot.Serializer
		bin    []byte
		newObj func() usbipprot.Serializer
		encErr error
		decErr error
	}{
		{
			name: "SetupPacket",
			obj: &protocol.SetupPacket{
				BMRequestType: 0xFF,
				BRequest:      0xFF,
				WValue:        0x1234,
				WIndex:        0x1234,
				WLength:       0x1234,
			},
			bin: []byte{
				0xFF,
				0xFF,
				0x34, 0x12,
				0x34, 0x12,
				0x34, 0x12,
			},
			newObj: func() usbipprot.Serializer {
				return &protocol.SetupPacket{}
			},
			encErr: nil,
			decErr: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			writer := new(bytes.Buffer)
			err := test.obj.Encode(writer)

			assert.ErrorIs(t, err, test.encErr)
			assert.Equal(t, test.bin, writer.Bytes())

			newObj := test.newObj()
			err = newObj.Decode(writer)

			assert.ErrorIs(t, err, test.decErr)
			assert.Equal(t, test.obj, newObj)
		})
	}
}
