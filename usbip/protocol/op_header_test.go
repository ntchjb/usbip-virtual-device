package protocol_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
)

func TestOpHeader(t *testing.T) {
	tests := []struct {
		name          string
		objEnc        protocol.Serializer
		objGen        func() protocol.Serializer
		expectedBytes []byte
	}{
		{
			name: "OpHeader",
			objEnc: &protocol.OpHeader{
				Version:            protocol.VERSION,
				CommandOrReplyCode: protocol.OP_REQ_DEVLIST,
				Status:             protocol.OP_STATUS_ERROR,
			},
			objGen: func() protocol.Serializer {
				return &protocol.OpHeader{}
			},
			expectedBytes: []byte{
				0x01, 0x11,
				0x80, 0x05,
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
