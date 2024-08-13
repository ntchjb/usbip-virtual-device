package descriptor_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/ntchjb/usbip-virtual-device/usb/protocol/descriptor"
	usbipprot "github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
)

func TestStandardDeviceDescriptor(t *testing.T) {
	tests := []struct {
		name   string
		obj    usbipprot.Serializer
		bin    []byte
		newObj func() usbipprot.Serializer
		encErr error
		decErr error
	}{
		{
			name: "StandardDeviceDescriptor",
			obj: &descriptor.StandardDeviceDescriptor{
				BLength:            descriptor.STANDARD_DEVICE_DESCRIPTOR_LENGTH,
				BDescriptorType:    descriptor.DESCRIPTOR_TYPE_DEVICE,
				BCDUSB:             0x1000,
				BDeviceClass:       protocol.CLASS_AUDIO,
				BDeviceSubClass:    protocol.SUBCLASS_HID_BOOT_INTERFACE,
				BDeviceProtocol:    protocol.PROTOCOL_HID_MOUSE,
				BMaxPacketSize:     0x08,
				IDVendor:           0xAABB,
				IDProduct:          0x0001,
				BCDDevice:          0x0100,
				IManufacturer:      0x04,
				IProduct:           0x0E,
				ISerialNumber:      0x30,
				BNumConfigurations: 0x01,
			},
			bin: []byte{
				0x12,
				0x01,
				0x00, 0x10,
				0x01,
				0x01,
				0x02,
				0x08,
				0xBB, 0xAA,
				0x01, 0x00,
				0x00, 0x01,
				0x04,
				0x0E,
				0x30,
				0x01,
			},
			newObj: func() usbipprot.Serializer {
				return &descriptor.StandardDeviceDescriptor{}
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
