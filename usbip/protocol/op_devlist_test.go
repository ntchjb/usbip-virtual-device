package protocol_test

import (
	"bytes"
	"testing"

	usbprotocol "github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
)

func TestOpDevlist(t *testing.T) {
	tests := []struct {
		name          string
		objEnc        protocol.Serializer
		objGen        func() protocol.Serializer
		expectedBytes []byte
	}{
		{
			name: "OpRepDevList",
			objEnc: &protocol.OpRepDevList{
				DeviceCount: 2,
				Devices: []protocol.DeviceInfo{
					{
						DeviceInfoTruncated: *deviceInfoTruncated,
						Interfaces: []protocol.DeviceInterface{
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
					{
						DeviceInfoTruncated: *deviceInfoTruncated,
						Interfaces: []protocol.DeviceInterface{
							{
								BInterfaceClass:    usbprotocol.CLASS_HID,
								BInterfaceSubclass: usbprotocol.SUBCLASS_HID_BOOT_INTERFACE,
								BInterfaceProtocol: usbprotocol.PROTOCOL_HID_MOUSE,
								PaddingAlignment:   0,
							},
							{
								BInterfaceClass:    usbprotocol.CLASS_AUDIO,
								BInterfaceSubclass: 0xA1,
								BInterfaceProtocol: 0xF1,
								PaddingAlignment:   0,
							},
							{
								BInterfaceClass:    usbprotocol.CLASS_AUDIO_AND_VIDEO,
								BInterfaceSubclass: 0xA2,
								BInterfaceProtocol: 0xF2,
								PaddingAlignment:   0,
							},
						},
					},
				},
			},
			objGen: func() protocol.Serializer {
				return &protocol.OpRepDevList{}
			},
			expectedBytes: appendBytes(
				[]byte{
					0x00, 0x00, 0x00, 0x02,
				},
				deviceInfoTruncatedExpectedBytes,
				[]byte{
					0x03,
					0x01,
					0x02,
					0x00,
				},
				[]byte{
					0x01,
					0xAB,
					0xFF,
					0x00,
				},
				[]byte{
					0x10,
					0xAA,
					0xFE,
					0x00,
				},
				deviceInfoTruncatedExpectedBytes,
				[]byte{
					0x03,
					0x01,
					0x02,
					0x00,
				},
				[]byte{
					0x01,
					0xA1,
					0xF1,
					0x00,
				},
				[]byte{
					0x10,
					0xA2,
					0xF2,
					0x00,
				},
			),
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
