package op_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol/op"
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
			objEnc: &op.OpHeader{
				Version:            op.VERSION,
				CommandOrReplyCode: op.OP_REQ_DEVLIST,
				Status:             op.OP_STATUS_ERROR,
			},
			objGen: func() protocol.Serializer {
				return &op.OpHeader{}
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
