package protocol_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol"
	usbipprot "github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
)

func TestStandardHIDDescriptor(t *testing.T) {
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
			obj: &protocol.HIDDescriptor{
				BLength:              protocol.HID_DESCRIPTOR_LENGTH,
				BDescriptorType:      protocol.DESCRIPTOR_TYPE_HID,
				BCDHID:               0x0110,
				BCountryCode:         0x01,
				BNumDescriptors:      0x01,
				BClassDescriptorType: 0x22,
				WDescriptorLength:    0x003F,
			},
			bin: []byte{
				0x09,
				0x21,
				0x10, 0x01,
				0x01,
				0x01,
				0x22,
				0x3F, 0x00,
			},
			newObj: func() usbipprot.Serializer {
				return &protocol.HIDDescriptor{}
			},
			encErr: nil,
			decErr: nil,
		},
		{
			name: "StandardEndpointDescriptor - With Optional",
			obj: &protocol.HIDDescriptor{
				BLength:                   protocol.HID_DESCRIPTOR_LENGTH + 3,
				BDescriptorType:           protocol.DESCRIPTOR_TYPE_HID,
				BCDHID:                    0x0110,
				BCountryCode:              0x01,
				BNumDescriptors:           0x01,
				BClassDescriptorType:      0x22,
				WDescriptorLength:         0x003F,
				BOptionalDescriptorType:   protocol.DESCRIPTOR_TYPE_HID,
				BOptionalDescriptorLength: protocol.HID_DESCRIPTOR_LENGTH,
			},
			bin: []byte{
				0x0c,
				0x21,
				0x10, 0x01,
				0x01,
				0x01,
				0x22,
				0x3F, 0x00,
				0x21,
				0x09, 0x00,
			},
			newObj: func() usbipprot.Serializer {
				return &protocol.HIDDescriptor{}
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
