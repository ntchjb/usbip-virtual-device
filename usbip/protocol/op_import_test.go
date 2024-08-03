package protocol_test

import (
	"bytes"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
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
			objEnc: &protocol.OpReqImport{
				BusID: deviceInfoBusID,
			},
			objGen: func() protocol.Serializer {
				return &protocol.OpReqImport{}
			},
			expectedBytes: deviceInfoBusID[:],
		},
		{
			name: "OpRepImport",
			objEnc: &protocol.OpRepImport{
				DeviceInfo: *deviceInfoTruncated,
			},
			objGen: func() protocol.Serializer {
				return &protocol.OpRepImport{}
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
