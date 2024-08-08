package protocol_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol"
	usbipprot "github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
)

func TestStandardConfigurationDescriptor(t *testing.T) {
	tests := []struct {
		name   string
		obj    usbipprot.Serializer
		bin    []byte
		newObj func() usbipprot.Serializer
		encErr error
		decErr error
	}{
		{
			name: "StandardConfigurationDescriptorEncode",
			obj: &protocol.StandardConfigurationDescriptor{
				BLength:             protocol.STANDARD_CONFIGURATION_DESCRIPTOR_LENGTH,
				BDescriptorType:     protocol.DESCRIPTOR_TYPE_CONFIGURATION,
				WTotalLength:        0x003B,
				BNumInterfaces:      0x02,
				BConfigurationValue: 0x01,
				IConfiguration:      0x00,
				BMAttributes:        0b10100000,
				BMaxPower:           0x32,
			},
			bin: []byte{
				0x09,
				0x02,
				0x3B, 0x00,
				0x02,
				0x01,
				0x00,
				0b10100000,
				0x32,
			},
			newObj: func() usbipprot.Serializer {
				return &protocol.StandardConfigurationDescriptor{}
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
