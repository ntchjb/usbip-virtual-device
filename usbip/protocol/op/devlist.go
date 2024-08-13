package op

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

// Decode read data from stream and store in struct.
// Note that this function does not decode OpHeader, which should be done already during connection handling
func (op *OpRepDevList) Decode(reader io.Reader) error {
	deviceCountBuf, err := stream.Read(reader, 4)
	if err != nil {
		return fmt.Errorf("unable to read device count buf from stream: %w", err)
	}
	op.DeviceCount = binary.BigEndian.Uint32(deviceCountBuf)
	op.Devices = make([]DeviceInfo, op.DeviceCount)
	for i := 0; i < int(op.DeviceCount); i++ {
		if err := op.Devices[i].Decode(reader); err != nil {
			return fmt.Errorf("unable to decode device info: %w", err)
		}
	}

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode OpHeader, which should be done already during connection handling
func (op *OpRepDevList) Encode(writer io.Writer) error {
	deviceCountBuf := make([]byte, 4)
	if int(op.DeviceCount) != len(op.Devices) {
		return fmt.Errorf("expected device count does not match actual device count, expected %d, got %d", op.DeviceCount, len(op.Devices))
	}
	binary.BigEndian.PutUint32(deviceCountBuf, op.DeviceCount)
	if err := stream.Write(writer, deviceCountBuf); err != nil {
		return fmt.Errorf("unable to write deviceCount to stream: %w", err)
	}

	for _, device := range op.Devices {
		if err := device.Encode(writer); err != nil {
			return fmt.Errorf("unable to encode DeviceInfo: %w", err)
		}
	}

	return nil
}
