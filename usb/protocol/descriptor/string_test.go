package descriptor_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol/descriptor"
	usbipprot "github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
)

func TestStringDescriptor(t *testing.T) {
	tests := []struct {
		name   string
		obj    usbipprot.Serializer
		bin    []byte
		newObj func() usbipprot.Serializer
		encErr error
		decErr error
	}{
		{
			name: "StringDescriptor",
			obj: &descriptor.StringDescriptor{
				BLength:         0x0A,
				BDescriptorType: descriptor.DESCRIPTOR_TYPE_STRING,
				Content: []uint16{
					0x0001, 0x1001, 0x1234, 0x2345,
				},
			},
			bin: []byte{
				0x0a,
				0x03,
				0x01, 0x00,
				0x01, 0x10,
				0x34, 0x12,
				0x45, 0x23,
			},
			newObj: func() usbipprot.Serializer {
				return &descriptor.StringDescriptor{}
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
