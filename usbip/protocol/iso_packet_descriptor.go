package protocol

import (
	"encoding/binary"
	"fmt"
)

func (c *ISOPacketDescriptor) MarshalBinaryPreAlloc(data []byte) error {
	if len(data) < c.Length() {
		return fmt.Errorf("data too short to allocate for ISOPacketDescriptor, need %d, got %d", c.Length(), len(data))
	}

	binary.BigEndian.PutUint32(data[:4], c.Offset)
	binary.BigEndian.PutUint32(data[4:8], c.ExpectedLength)
	binary.BigEndian.PutUint32(data[8:12], c.ActualLength)
	binary.BigEndian.PutUint32(data[12:16], c.Status)

	return nil
}

func (c *ISOPacketDescriptor) MarshalBinary() (data []byte, err error) {
	data = make([]byte, c.Length())

	if err := c.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate ISOPacketDescriptor: %w", err)
	}

	return data, nil
}

func (c *ISOPacketDescriptor) UnmarshalBinaryWithLength(data []byte) (int, error) {
	var length int

	if len(data) < ISO_PACKET_DESCRIPTOR_LENGTH {
		return 0, fmt.Errorf("data too short for ISOPacketDescriptor, need %d, got %d", ISO_PACKET_DESCRIPTOR_LENGTH, len(data))
	}

	c.Offset = binary.BigEndian.Uint32(data[:4])
	c.ExpectedLength = binary.BigEndian.Uint32(data[4:8])
	c.ActualLength = binary.BigEndian.Uint32(data[8:12])
	c.Status = binary.BigEndian.Uint32(data[12:16])

	length += ISO_PACKET_DESCRIPTOR_LENGTH

	return length, nil
}

func (c *ISOPacketDescriptor) UnmarshalBinary(data []byte) error {
	_, err := c.UnmarshalBinaryWithLength(data)

	return err
}

func (c *ISOPacketDescriptor) Length() int {
	return ISO_PACKET_DESCRIPTOR_LENGTH
}
