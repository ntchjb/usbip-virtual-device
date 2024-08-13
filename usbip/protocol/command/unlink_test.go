package command_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol/command"
	"github.com/stretchr/testify/assert"
)

func TestCmdUnlink(t *testing.T) {
	tests := []struct {
		name          string
		objEnc        protocol.Serializer
		objGen        func() protocol.Serializer
		expectedBytes []byte
	}{
		{
			name: "CmdUnlink",
			objEnc: &command.CmdUnlink{
				UnlinkSeqNum: 0x01020304,
				Padding:      [24]byte{},
			},
			objGen: func() protocol.Serializer {
				return &command.CmdUnlink{}
			},
			expectedBytes: []byte{
				0x01, 0x02, 0x03, 0x04,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
		{
			name: "RetUnlink",
			objEnc: &command.RetUnlink{
				Status:  0x00000001,
				Padding: [24]byte{},
			},
			objGen: func() protocol.Serializer {
				return &command.RetUnlink{}
			},
			expectedBytes: []byte{
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
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
