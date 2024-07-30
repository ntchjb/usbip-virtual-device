package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

// Decode read data from stream and store in struct.
// Note that this function does not decode CmdHeader, which should be done already during connection handling
func (c *CmdUnlink) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, CMD_UNLINK_STATIC_FIELDS_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read CmdUnlink from stream: %w", err)
	}

	c.UnlinkSeqNum = binary.BigEndian.Uint32(buf[:4])
	copy(c.Padding[:], buf[4:28])

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode CmdHeader, which should be done already during connection handling
func (c *CmdUnlink) Encode(writer io.Writer) error {
	buf := make([]byte, CMD_UNLINK_STATIC_FIELDS_LENGTH)

	binary.BigEndian.PutUint32(buf[:4], c.UnlinkSeqNum)
	copy(buf[4:28], c.Padding[:])

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write CmdUnlink to stream: %w", err)
	}

	return nil
}

// Decode read data from stream and store in struct.
// Note that this function does not decode CmdHeader, which should be done already during connection handling
func (c *RetUnlink) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, RET_UNLINK_STATIC_FIELDS_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read RetUnlink from stream: %w", err)
	}

	c.Status = binary.BigEndian.Uint32(buf[:4])
	copy(c.Padding[:], buf[4:28])

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode CmdHeader, which should be done already during connection handling
func (c *RetUnlink) Encode(writer io.Writer) error {
	buf := make([]byte, RET_UNLINK_STATIC_FIELDS_LENGTH)
	binary.BigEndian.PutUint32(buf[:4], c.Status)
	copy(buf[4:28], c.Padding[:])

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write RetUnlink to stream: %w", err)
	}

	return nil
}
