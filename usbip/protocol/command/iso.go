package command

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

// Decode read data from stream and store in struct.
// Note that this function does not decode CmdHeader, which should be done already during connection handling
func (c *ISOPacketDescriptor) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, ISO_PACKET_DESCRIPTOR_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read ISOPacketDescriptor from stream: %w", err)
	}

	c.Offset = binary.BigEndian.Uint32(buf[:4])
	c.ExpectedLength = binary.BigEndian.Uint32(buf[4:8])
	c.ActualLength = binary.BigEndian.Uint32(buf[8:12])
	c.Status = binary.BigEndian.Uint32(buf[12:16])

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode CmdHeader, which should be done already during connection handling
func (c *ISOPacketDescriptor) Encode(writer io.Writer) error {
	buf := make([]byte, ISO_PACKET_DESCRIPTOR_LENGTH)

	binary.BigEndian.PutUint32(buf[:4], c.Offset)
	binary.BigEndian.PutUint32(buf[4:8], c.ExpectedLength)
	binary.BigEndian.PutUint32(buf[8:12], c.ActualLength)
	binary.BigEndian.PutUint32(buf[12:16], c.Status)

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write ISOPacketDescriptor to stream: %w", err)
	}

	return nil
}
