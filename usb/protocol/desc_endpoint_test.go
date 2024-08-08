package protocol_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol"
	usbipprot "github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
)

func TestStandardEndpointDescriptor(t *testing.T) {
	tests := []struct {
		name   string
		obj    usbipprot.Serializer
		bin    []byte
		newObj func() usbipprot.Serializer
		encErr error
		decErr error
	}{
		{
			name: "StandardEndpointDescriptor",
			obj: &protocol.StandardEndpointDescriptor{
				BLength:          protocol.STANDARD_ENDPOINT_DESCRIPTOR_LENGTH,
				BDescriptorType:  protocol.DESCRIPTOR_TYPE_ENDPOINT,
				BEndpointAddress: 0b10000001,
				BMAttributes:     0b00000011,
				WMaxPacketSize:   0x0008,
				BInterval:        0x0A,
			},
			bin: []byte{
				0x07,
				0x05,
				0b10000001,
				0b00000011,
				0x08, 0x00,
				0x0A,
			},
			newObj: func() usbipprot.Serializer {
				return &protocol.StandardEndpointDescriptor{}
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
