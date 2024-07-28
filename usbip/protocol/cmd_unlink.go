package protocol

import (
	"encoding/binary"
	"fmt"
)

func (c *CmdUnlink) MarshalBinaryPreAlloc(data []byte) error {
	if len(data) < c.Length() {
		return fmt.Errorf("data too short to allocate for CmdUnlink, need %d, got %d", c.Length(), len(data))
	}

	headerLength := c.CmdHeader.Length()
	if err := c.CmdHeader.MarshalBinaryPreAlloc(data[:headerLength]); err != nil {
		return fmt.Errorf("unable to allocate CmdHeader: %w", err)
	}

	binary.BigEndian.PutUint32(data[headerLength:headerLength+4], c.UnlinkSeqNum)
	copy(data[headerLength+4:headerLength+28], c.Padding[:])

	return nil
}

func (c *CmdUnlink) MarshalBinary() (data []byte, err error) {
	data = make([]byte, c.Length())

	if err := c.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate CmdUnlink: %w", err)
	}

	return data, nil
}

func (c *CmdUnlink) UnmarshalBinaryWithLength(data []byte) (int, error) {
	var length int

	if len(data) < CMD_HEADER_LENGTH {
		return 0, fmt.Errorf("data too short for CmdHeader, need %d, got %d", c.Length(), len(data))
	}

	cmdHeaderLength, err := c.CmdHeader.UnmarshalBinaryWithLength(data)
	if err != nil {
		return 0, fmt.Errorf("unable to unmarshal CmdHeader: %w", err)
	}
	length += cmdHeaderLength

	if len(data[length:]) < CMD_UNLINK_STATIC_FIELDS_LENGTH {
		return 0, fmt.Errorf("data too short for static fields of CmdUnlink, need %d, got %d", CMD_UNLINK_STATIC_FIELDS_LENGTH, len(data))
	}

	c.UnlinkSeqNum = binary.BigEndian.Uint32(data[length : length+4])
	copy(c.Padding[:], data[length+4:length+28])
	length += CMD_UNLINK_STATIC_FIELDS_LENGTH

	return length, nil
}

func (c *CmdUnlink) UnmarshalBinary(data []byte) error {
	_, err := c.UnmarshalBinaryWithLength(data)

	return err
}

func (c *CmdUnlink) Length() int {
	return CMD_HEADER_LENGTH + CMD_UNLINK_STATIC_FIELDS_LENGTH
}

func (c *RetUnlink) MarshalBinaryPreAlloc(data []byte) error {
	if len(data) < c.Length() {
		return fmt.Errorf("data too short to allocate for RetUnlink, need %d, got %d", c.Length(), len(data))
	}

	headerLength := c.CmdHeader.Length()
	if err := c.CmdHeader.MarshalBinaryPreAlloc(data[:headerLength]); err != nil {
		return fmt.Errorf("unable to allocate CmdHeader: %w", err)
	}

	binary.BigEndian.PutUint32(data[headerLength:headerLength+4], c.Status)
	copy(data[headerLength+4:headerLength+28], c.Padding[:])

	return nil
}

func (c *RetUnlink) MarshalBinary() (data []byte, err error) {
	data = make([]byte, c.Length())

	if err := c.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate RetUnlink: %w", err)
	}

	return data, nil
}

func (c *RetUnlink) UnmarshalBinaryWithLength(data []byte) (int, error) {
	var length int

	if len(data) < CMD_HEADER_LENGTH {
		return 0, fmt.Errorf("data too short for CmdHeader, need %d, got %d", c.Length(), len(data))
	}

	cmdHeaderLength, err := c.CmdHeader.UnmarshalBinaryWithLength(data)
	if err != nil {
		return 0, fmt.Errorf("unable to unmarshal CmdHeader: %w", err)
	}
	length += cmdHeaderLength

	if len(data[length:]) < RET_UNLINK_STATIC_FIELDS_LENGTH {
		return 0, fmt.Errorf("data too short for static fields of RetUnlink, need %d, got %d", RET_UNLINK_STATIC_FIELDS_LENGTH, len(data))
	}

	c.Status = binary.BigEndian.Uint32(data[length : length+4])
	copy(c.Padding[:], data[length+4:length+28])
	length += RET_UNLINK_STATIC_FIELDS_LENGTH

	return length, nil
}

func (c *RetUnlink) UnmarshalBinary(data []byte) error {
	_, err := c.UnmarshalBinaryWithLength(data)

	return err
}

func (c *RetUnlink) Length() int {
	return CMD_HEADER_LENGTH + RET_UNLINK_STATIC_FIELDS_LENGTH
}
