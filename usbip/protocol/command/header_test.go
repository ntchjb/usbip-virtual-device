package command_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol/command"
	"github.com/stretchr/testify/assert"
)

func TestCmdHeader(t *testing.T) {
	tests := []struct {
		name          string
		objEnc        protocol.Serializer
		objGen        func() protocol.Serializer
		expectedBytes []byte
	}{
		{
			name: "CmdHeader",
			objEnc: &command.CmdHeader{
				Command:        command.RET_SUBMIT,
				SeqNum:         0x12345678,
				DevID:          0x0001000A,
				Direction:      command.DIR_IN,
				EndpointNumber: 0x0000000A,
			},
			objGen: func() protocol.Serializer {
				return &command.CmdHeader{}
			},
			expectedBytes: []byte{
				0x00, 0x00, 0x00, 0x03,
				0x12, 0x34, 0x56, 0x78,
				0x00, 0x01, 0x00, 0x0A,
				0x00, 0x00, 0x00, 0x01,
				0x00, 0x00, 0x00, 0x0A,
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
