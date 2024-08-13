package op

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

// Decode read data from stream and store in struct.
// Note that this function does not decode OpHeader, which should be done already during connection handling
func (op *OpHeader) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, OP_HEADER_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read data from stream: %w", err)
	}

	op.Version = binary.BigEndian.Uint16(buf[:2])
	op.CommandOrReplyCode = Operation(binary.BigEndian.Uint16(buf[2:4]))
	op.Status = OperationStatus(binary.BigEndian.Uint32(buf[4:8]))

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode OpHeader, which should be done already during connection handling
func (op *OpHeader) Encode(writer io.Writer) error {
	buf := make([]byte, OP_HEADER_LENGTH)

	binary.BigEndian.PutUint16(buf[:2], op.Version)
	binary.BigEndian.PutUint16(buf[2:4], uint16(op.CommandOrReplyCode))
	binary.BigEndian.PutUint32(buf[4:8], uint32(op.Status))

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write data to stream: %w", err)
	}

	return nil
}
