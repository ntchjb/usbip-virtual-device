package protocol_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
)

func TestCmdISO(t *testing.T) {
	tests := []struct {
		name          string
		objEnc        protocol.Serializer
		objGen        func() protocol.Serializer
		expectedBytes []byte
	}{
		{
			name: "CmdISO",
			objEnc: &protocol.ISOPacketDescriptor{
				Offset:         0x01020304,
				ExpectedLength: 0x01010101,
				ActualLength:   0x01020304,
				Status:         0x00000001,
			},
			objGen: func() protocol.Serializer {
				return &protocol.ISOPacketDescriptor{}
			},
			expectedBytes: []byte{
				0x01, 0x02, 0x03, 0x04,
				0x01, 0x01, 0x01, 0x01,
				0x01, 0x02, 0x03, 0x04,
				0x00, 0x00, 0x00, 0x01,
			},
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
