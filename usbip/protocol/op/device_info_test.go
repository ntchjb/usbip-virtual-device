package op_test

import (
	"bytes"
	"testing"

	usbprotocol "github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol/op"
	"github.com/stretchr/testify/assert"
)

var (
	deviceInfoPath = [256]byte{
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0xBB, 0xBB, 0xBB, 0xBB, 0xBB, 0xBB,
	}
	deviceInfoBusID = usbprotocol.BusID{
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0xBB, 0xBB,
	}

	deviceInfoTruncated = &op.DeviceInfoTruncated{
		Path:                deviceInfoPath,
		BusID:               deviceInfoBusID,
		BusNum:              1,
		DevNum:              5,
		Speed:               usbprotocol.SPEED_USB2_HIGH,
		IDVendor:            0xABCD,
		IDProduct:           0xDCBA,
		BCDDevice:           127,
		BDeviceClass:        usbprotocol.CLASS_HID,
		BDeviceSubclass:     usbprotocol.SUBCLASS_HID_BOOT_INTERFACE,
		BDeviceProtocol:     usbprotocol.PROTOCOL_HID_KEYBOARD,
		BConfigurationValue: 1,
		BNumConfigurations:  3,
		BNumInterfaces:      3,
	}
	deviceInfoTruncatedExpectedBytes = appendBytes(
		deviceInfoPath[:],
		deviceInfoBusID[:],
		[]byte{
			0x00, 0x00, 0x00, 0x01,
			0x00, 0x00, 0x00, 0x05,
			0x00, 0x00, 0x00, 0x02,
			0xAB, 0xCD,
			0xDC, 0xBA,
			0x00, 0x7F,
			0x03,
			0x01,
			0x01,
			0x01,
			0x03,
			0x03,
		})
)

func appendBytes(arrBytes ...[]byte) []byte {
	res := []byte{}
	for _, bytes := range arrBytes {
		res = append(res, bytes...)
	}

	return res
}

func TestOpDeviceInfo(t *testing.T) {
	tests := []struct {
		name          string
		objEnc        protocol.Serializer
		objGen        func() protocol.Serializer
		expectedBytes []byte
	}{
		{
			name:   "DeviceInfoTruncated",
			objEnc: deviceInfoTruncated,
			objGen: func() protocol.Serializer {
				return &op.DeviceInfoTruncated{}
			},
			expectedBytes: deviceInfoTruncatedExpectedBytes,
		},
		{
			name: "DeviceInterface",
			objEnc: &op.DeviceInterface{
				BInterfaceClass:    usbprotocol.CLASS_HID,
				BInterfaceSubclass: usbprotocol.SUBCLASS_HID_BOOT_INTERFACE,
				BInterfaceProtocol: usbprotocol.PROTOCOL_HID_MOUSE,
				PaddingAlignment:   0,
			},
			objGen: func() protocol.Serializer {
				return &op.DeviceInterface{}
			},
			expectedBytes: []byte{
				0x03,
				0x01,
				0x02,
				0x00,
			},
		},
		{
			name: "DeviceInfo",
			objEnc: &op.DeviceInfo{
				DeviceInfoTruncated: *deviceInfoTruncated,
				Interfaces: []op.DeviceInterface{
					{
						BInterfaceClass:    usbprotocol.CLASS_HID,
						BInterfaceSubclass: usbprotocol.SUBCLASS_HID_BOOT_INTERFACE,
						BInterfaceProtocol: usbprotocol.PROTOCOL_HID_MOUSE,
						PaddingAlignment:   0,
					},
					{
						BInterfaceClass:    usbprotocol.CLASS_AUDIO,
						BInterfaceSubclass: 0xAB,
						BInterfaceProtocol: 0xFF,
						PaddingAlignment:   0,
					},
					{
						BInterfaceClass:    usbprotocol.CLASS_AUDIO_AND_VIDEO,
						BInterfaceSubclass: 0xAA,
						BInterfaceProtocol: 0xFE,
						PaddingAlignment:   0,
					},
				},
			},
			objGen: func() protocol.Serializer {
				return &op.DeviceInfo{}
			},
			expectedBytes: appendBytes(deviceInfoTruncatedExpectedBytes, []byte{
				0x03,
				0x01,
				0x02,
				0x00,
			}, []byte{
				0x01,
				0xAB,
				0xFF,
				0x00,
			}, []byte{
				0x10,
				0xAA,
				0xFE,
				0x00,
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			buf := new(bytes.Buffer)
			encodeErr := test.objEnc.Encode(buf)

			assert.NoError(t, encodeErr)
			assert.Equal(t, test.expectedBytes, buf.Bytes())

			newObj := test.objGen()
			decodeErr := newObj.Decode(buf)

			assert.NoError(t, decodeErr)
			assert.Equal(t, test.objEnc, newObj)
		})
	}
}
