package hid_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid"
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
			obj: &hid.HIDDescriptor{
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
				return &hid.HIDDescriptor{}
			},
			encErr: nil,
			decErr: nil,
		},
		{
			name: "StandardEndpointDescriptor - With Optional",
			obj: &hid.HIDDescriptor{
				BLength:              protocol.HID_DESCRIPTOR_LENGTH + 3,
				BDescriptorType:      protocol.DESCRIPTOR_TYPE_HID,
				BCDHID:               0x0110,
				BCountryCode:         0x01,
				BNumDescriptors:      0x03,
				BClassDescriptorType: 0x22,
				WDescriptorLength:    0x003F,
				OptionalDescriptorTypes: []hid.OptionalHIDDescriptorTypes{
					{
						BOptionalDescriptorType:   protocol.DESCRIPTOR_TYPE_HID_REPORT,
						BOptionalDescriptorLength: 0x0041,
					},
					{
						BOptionalDescriptorType:   protocol.DESCRIPTOR_TYPE_HID_REPORT,
						BOptionalDescriptorLength: 0x0042,
					},
				},
			},
			bin: []byte{
				0x0c,
				0x21,
				0x10, 0x01,
				0x01,
				0x03,
				0x22,
				0x3F, 0x00,
				0x22,
				0x41, 0x00,
				0x22,
				0x42, 0x00,
			},
			newObj: func() usbipprot.Serializer {
				return &hid.HIDDescriptor{}
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
