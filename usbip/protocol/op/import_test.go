package op_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol/op"
	"github.com/stretchr/testify/assert"
)

func TestOpImport(t *testing.T) {
	tests := []struct {
		name          string
		objEnc        protocol.Serializer
		objGen        func() protocol.Serializer
		expectedBytes []byte
	}{
		{
			name: "OpReqImport",
			objEnc: &op.OpReqImport{
				BusID: deviceInfoBusID,
			},
			objGen: func() protocol.Serializer {
				return &op.OpReqImport{}
			},
			expectedBytes: deviceInfoBusID[:],
		},
		{
			name: "OpRepImport",
			objEnc: &op.OpRepImport{
				DeviceInfo: *deviceInfoTruncated,
			},
			objGen: func() protocol.Serializer {
				return &op.OpRepImport{}
			},
			expectedBytes: deviceInfoTruncatedExpectedBytes,
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
