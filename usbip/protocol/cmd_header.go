package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

func (c *CmdHeader) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, CMD_HEADER_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read CmdHeader from stream: %w", err)
	}

	c.Command = Command(binary.BigEndian.Uint32(buf[:4]))
	c.SeqNum = binary.BigEndian.Uint32(buf[4:8])
	c.DevID = binary.BigEndian.Uint32(buf[8:12])
	c.Direction = Direction(binary.BigEndian.Uint32(buf[12:16]))
	c.EndpointNumber = binary.BigEndian.Uint32(buf[16:20])

	return nil
}

func (c *CmdHeader) Encode(writer io.Writer) error {
	buf := make([]byte, CMD_HEADER_LENGTH)

	binary.BigEndian.PutUint32(buf[:4], uint32(c.Command))
	binary.BigEndian.PutUint32(buf[4:8], c.SeqNum)
	binary.BigEndian.PutUint32(buf[8:12], c.DevID)
	binary.BigEndian.PutUint32(buf[12:16], uint32(c.Direction))
	binary.BigEndian.PutUint32(buf[16:20], c.EndpointNumber)

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write CmdHeader to stream: %w", err)
	}

	return nil
}
