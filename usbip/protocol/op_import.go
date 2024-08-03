package protocol

import (
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

// Decode read data from stream and store in struct.
// Note that this function does not decode OpHeader, which should be done already during connection handling
func (op *OpReqImport) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, len(op.BusID))
	if err != nil {
		return fmt.Errorf("unable to read BusID from stream: %w", err)
	}

	copy(op.BusID[:], buf)

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode OpHeader, which should be done already during connection handling
func (op *OpReqImport) Encode(writer io.Writer) error {
	if err := stream.Write(writer, op.BusID[:]); err != nil {
		return fmt.Errorf("unable to write BusID to stream: %w", err)
	}

	return nil
}

// Decode read data from stream and store in struct.
// Note that this function does not decode OpHeader, which should be done already during connection handling
func (op *OpRepImport) Decode(reader io.Reader) error {
	if err := op.DeviceInfo.Decode(reader); err != nil {
		return fmt.Errorf("unable to decode DeviceInfoTruncated: %w", err)
	}
	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode OpHeader, which should be done already during connection handling
func (op *OpRepImport) Encode(writer io.Writer) error {
	if err := op.DeviceInfo.Encode(writer); err != nil {
		return fmt.Errorf("unable to encode DeviceInfoTruncated: %w", err)
	}
	return nil
}
