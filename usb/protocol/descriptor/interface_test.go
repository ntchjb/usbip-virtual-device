package descriptor_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/ntchjb/usbip-virtual-device/usb/protocol/descriptor"
	usbipprot "github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
)

func TestStandardInterfaceDescriptor(t *testing.T) {
	tests := []struct {
		name   string
		obj    usbipprot.Serializer
		bin    []byte
		newObj func() usbipprot.Serializer
		encErr error
		decErr error
	}{
		{
			name: "StandardInterfaceDescriptor",
			obj: &descriptor.StandardInterfaceDescriptor{
				BLength:            descriptor.HID_DESCRIPTOR_LENGTH,
				BDescriptorType:    descriptor.DESCRIPTOR_TYPE_HID,
				BInterfaceNumber:   0x01,
				BAlternateSetting:  0x01,
				BNumEndpoints:      0x01,
				BInterfaceClass:    protocol.CLASS_HID,
				BInterfaceSubClass: protocol.SUBCLASS_HID_BOOT_INTERFACE,
				BInterfaceProtocol: protocol.PROTOCOL_HID_MOUSE,
				IInterface:         0x03,
			},
			bin: []byte{
				0x09,
				0x21,
				0x01,
				0x01,
				0x01,
				0x03,
				0x01,
				0x02,
				0x03,
			},
			newObj: func() usbipprot.Serializer {
				return &descriptor.StandardInterfaceDescriptor{}
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
