package protocol

import (
	"encoding/binary"
	"fmt"
)

func (c *CmdHeader) MarshalBinaryPreAlloc(data []byte) error {
	if len(data) < c.Length() {
		return fmt.Errorf("data too short to allocate for CmdHeader, need %d, got %d", c.Length(), len(data))
	}

	binary.BigEndian.PutUint32(data[:4], c.Command)
	binary.BigEndian.PutUint32(data[4:8], c.SeqNum)
	binary.BigEndian.PutUint32(data[8:12], c.DevID)
	binary.BigEndian.PutUint32(data[12:16], uint32(c.Direction))
	binary.BigEndian.PutUint32(data[16:20], c.EndpointNumber)

	return nil
}

func (c *CmdHeader) MarshalBinary() (data []byte, err error) {
	data = make([]byte, c.Length())

	if err := c.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate CmdHeader: %w", err)
	}

	return data, nil
}

func (c *CmdHeader) UnmarshalBinaryWithLength(data []byte) (int, error) {
	var length int

	if len(data) < CMD_HEADER_LENGTH {
		return 0, fmt.Errorf("data too short for CmdHeader, need %d, got %d", CMD_HEADER_LENGTH, len(data))
	}

	c.Command = binary.BigEndian.Uint32(data[:4])
	c.SeqNum = binary.BigEndian.Uint32(data[4:8])
	c.DevID = binary.BigEndian.Uint32(data[8:12])
	c.Direction = Direction(binary.BigEndian.Uint32(data[12:16]))
	c.EndpointNumber = binary.BigEndian.Uint32(data[16:20])

	length += CMD_HEADER_LENGTH

	return length, nil
}

func (c *CmdHeader) UnmarshalBinary(data []byte) error {
	_, err := c.UnmarshalBinaryWithLength(data)

	return err
}

func (c *CmdHeader) Length() int {
	return CMD_HEADER_LENGTH
}
