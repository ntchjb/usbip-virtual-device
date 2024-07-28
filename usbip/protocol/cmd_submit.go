package protocol

import (
	"encoding/binary"
	"fmt"
)

func (c *CmdSubmit) MarshalBinaryPreAlloc(data []byte) error {
	if c.Direction == DIR_OUT && len(c.TransferBuffer) != int(c.TransferBufferLength) {
		return fmt.Errorf("actual transfer buffer length does not match TransferBufferLength for DIR_OUT CmdSubmit; expected %d, actual %d", c.TransferBufferLength, len(c.TransferBuffer))
	}
	if (c.NumberOfPackets == 0x00000000 || c.NumberOfPackets == 0xffffffff) && len(c.ISOPacketDescriptors) > 0 {
		return fmt.Errorf("this is non-ISO transfer, but contains ISO packet descriptors")
	} else if c.NumberOfPackets != 0x00000000 && c.NumberOfPackets != 0xffffffff && int(c.NumberOfPackets) != len(c.ISOPacketDescriptors) {
		return fmt.Errorf("number of packets does not match with actual ISO packet descriptor length, expected %d, actual %d", c.NumberOfPackets, len(c.ISOPacketDescriptors))
	}

	if len(data) < c.Length() {
		return fmt.Errorf("data too short to allocate for CmdSubmit, need %d, got %d", c.Length(), len(data))
	}

	headerLength := c.CmdHeader.Length()
	if err := c.CmdHeader.MarshalBinaryPreAlloc(data[:headerLength]); err != nil {
		return fmt.Errorf("unable to allocate CmdHeader: %w", err)
	}

	binary.BigEndian.PutUint32(data[headerLength:headerLength+4], c.TransferFlags)
	binary.BigEndian.PutUint32(data[headerLength+4:headerLength+8], c.TransferBufferLength)
	binary.BigEndian.PutUint32(data[headerLength+8:headerLength+12], c.StartFrame)
	binary.BigEndian.PutUint32(data[headerLength+12:headerLength+16], c.NumberOfPackets)
	binary.BigEndian.PutUint32(data[headerLength+16:headerLength+20], c.Interval)
	binary.BigEndian.PutUint64(data[headerLength+20:headerLength+28], c.Setup)

	if c.TransferBufferLength > 0 {
		copy(data[headerLength+28:headerLength+28+int(c.TransferBufferLength)], c.TransferBuffer)
	}

	startIdx := headerLength + 28 + int(c.TransferBufferLength)
	for _, descriptor := range c.ISOPacketDescriptors {
		if err := descriptor.MarshalBinaryPreAlloc(data[startIdx : startIdx+descriptor.Length()]); err != nil {
			return fmt.Errorf("unable to allocate ISOPacketDescriptor: %w", err)
		}
		startIdx += descriptor.Length()
	}

	return nil
}

func (c *CmdSubmit) MarshalBinary() (data []byte, err error) {
	data = make([]byte, c.Length())

	if err := c.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate CmdSubmit: %w", err)
	}

	return data, nil
}

func (c *CmdSubmit) UnmarshalBinaryWithLength(data []byte) (int, error) {
	var length int

	if len(data) < CMD_HEADER_LENGTH {
		return 0, fmt.Errorf("data too short for CmdHeader, need %d, got %d", c.Length(), len(data))
	}

	cmdHeaderLength, err := c.CmdHeader.UnmarshalBinaryWithLength(data)
	if err != nil {
		return 0, fmt.Errorf("unable to unmarshal CmdHeader: %w", err)
	}
	length += cmdHeaderLength

	if len(data[length:]) < CMD_SUBMIT_STATIC_FIELDS_LENGTH {
		return 0, fmt.Errorf("data too short for static fields of CmdSubmit, need %d, got %d", CMD_SUBMIT_STATIC_FIELDS_LENGTH, len(data[length:]))
	}
	c.TransferFlags = binary.BigEndian.Uint32(data[length : length+4])
	c.TransferBufferLength = binary.BigEndian.Uint32(data[length+4 : length+8])
	c.StartFrame = binary.BigEndian.Uint32(data[length+8 : length+12])
	c.NumberOfPackets = binary.BigEndian.Uint32(data[length+12 : length+16])
	c.Interval = binary.BigEndian.Uint32(data[length+16 : length+20])
	c.Setup = binary.BigEndian.Uint64(data[length+20 : length+28])
	length += CMD_SUBMIT_STATIC_FIELDS_LENGTH

	if c.TransferBufferLength > 0 && c.Direction == DIR_OUT {
		if len(data[length:]) < int(c.TransferBufferLength) {
			return 0, fmt.Errorf("data too short for TransferBuffer, need %d, got %d", c.TransferBufferLength, len(data[length:]))
		}
		c.TransferBuffer = make([]byte, c.TransferBufferLength)
		copy(c.TransferBuffer, data[length:length+int(c.TransferBufferLength)])
		length += int(c.TransferBufferLength)
	}

	if c.NumberOfPackets != 0x00000000 && c.NumberOfPackets != 0xffffffff {
		if len(data[length:]) < int(c.NumberOfPackets)*ISO_PACKET_DESCRIPTOR_LENGTH {
			return 0, fmt.Errorf("data too short for ISOPacketDescriptors, need %d, got %d", int(c.NumberOfPackets)*ISO_PACKET_DESCRIPTOR_LENGTH, len(data[length:]))
		}
		c.ISOPacketDescriptors = make([]ISOPacketDescriptor, c.NumberOfPackets)
		for i := range c.ISOPacketDescriptors {
			isoPacketDescriptorLength, err := c.ISOPacketDescriptors[i].UnmarshalBinaryWithLength(data[length:])
			if err != nil {
				return 0, fmt.Errorf("unable to unmarshal ISOPacketDescriptor #%d: %w", i, err)
			}
			length += isoPacketDescriptorLength
		}
	}

	return length, nil
}

func (c *CmdSubmit) UnmarshalBinary(data []byte) error {
	_, err := c.UnmarshalBinaryWithLength(data)

	return err
}

func (c *CmdSubmit) Length() int {
	length := CMD_HEADER_LENGTH + CMD_SUBMIT_STATIC_FIELDS_LENGTH
	if c.Direction == DIR_OUT {
		length += len(c.TransferBuffer)
	}
	length += ISO_PACKET_DESCRIPTOR_LENGTH * len(c.ISOPacketDescriptors)
	return length
}

func (c *RetSubmit) MarshalBinaryPreAlloc(data []byte) error {
	if c.Direction == DIR_IN && len(c.TransferBuffer) != int(c.ActualLength) {
		return fmt.Errorf("actual transfer buffer length does not match ActualLength for DIR_IN RetSubmit; expected %d, actual %d", c.ActualLength, len(c.TransferBuffer))
	}
	if (c.NumberOfPackets == 0x00000000 || c.NumberOfPackets == 0xffffffff) && len(c.ISOPacketDescriptors) > 0 {
		return fmt.Errorf("this is non-ISO transfer, but contains ISO packet descriptors")
	} else if c.NumberOfPackets != 0x00000000 && c.NumberOfPackets != 0xffffffff && int(c.NumberOfPackets) != len(c.ISOPacketDescriptors) {
		return fmt.Errorf("number of packets does not match with actual ISO packet descriptor length, expected %d, actual %d", c.NumberOfPackets, len(c.ISOPacketDescriptors))
	}

	if len(data) < c.Length() {
		return fmt.Errorf("data too short to allocate for CmdSubmit, need %d, got %d", c.Length(), len(data))
	}

	headerLength := c.CmdHeader.Length()
	if err := c.CmdHeader.MarshalBinaryPreAlloc(data[:headerLength]); err != nil {
		return fmt.Errorf("unable to allocate CmdHeader: %w", err)
	}

	binary.BigEndian.PutUint32(data[headerLength:headerLength+4], c.Status)
	binary.BigEndian.PutUint32(data[headerLength+4:headerLength+8], c.ActualLength)
	binary.BigEndian.PutUint32(data[headerLength+8:headerLength+12], c.StartFrame)
	binary.BigEndian.PutUint32(data[headerLength+12:headerLength+16], c.NumberOfPackets)
	binary.BigEndian.PutUint32(data[headerLength+16:headerLength+20], c.ErrorCount)
	binary.BigEndian.PutUint64(data[headerLength+20:headerLength+28], c.Padding)

	if c.ActualLength > 0 {
		copy(data[headerLength+28:headerLength+28+int(c.ActualLength)], c.TransferBuffer)
	}

	startIdx := headerLength + 28 + int(c.ActualLength)
	for _, descriptor := range c.ISOPacketDescriptors {
		if err := descriptor.MarshalBinaryPreAlloc(data[startIdx : startIdx+descriptor.Length()]); err != nil {
			return fmt.Errorf("unable to allocate ISOPacketDescriptor: %w", err)
		}
		startIdx += descriptor.Length()
	}

	return nil
}

func (c *RetSubmit) MarshalBinary() (data []byte, err error) {
	data = make([]byte, c.Length())

	if err := c.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate RetSubmit: %w", err)
	}

	return data, nil
}

func (c *RetSubmit) UnmarshalBinaryWithLength(data []byte) (int, error) {
	var length int

	if len(data) < CMD_HEADER_LENGTH {
		return 0, fmt.Errorf("data too short for CmdHeader, need %d, got %d", c.Length(), len(data))
	}

	cmdHeaderLength, err := c.CmdHeader.UnmarshalBinaryWithLength(data)
	if err != nil {
		return 0, fmt.Errorf("unable to unmarshal CmdHeader: %w", err)
	}
	length += cmdHeaderLength

	if len(data[length:]) < RET_SUBMIT_STATIC_FIELDS_LENGTH {
		return 0, fmt.Errorf("data too short for static fields of RetSubmit, need %d, got %d", RET_SUBMIT_STATIC_FIELDS_LENGTH, len(data[length:]))
	}
	c.Status = binary.BigEndian.Uint32(data[length : length+4])
	c.ActualLength = binary.BigEndian.Uint32(data[length+4 : length+8])
	c.StartFrame = binary.BigEndian.Uint32(data[length+8 : length+12])
	c.NumberOfPackets = binary.BigEndian.Uint32(data[length+12 : length+16])
	c.ErrorCount = binary.BigEndian.Uint32(data[length+16 : length+20])
	c.Padding = binary.BigEndian.Uint64(data[length+20 : length+28])
	length += RET_SUBMIT_STATIC_FIELDS_LENGTH

	if c.ActualLength > 0 && c.Direction == DIR_IN {
		if len(data[length:]) < int(c.ActualLength) {
			return 0, fmt.Errorf("data too short for TransferBuffer, need %d, got %d", c.ActualLength, len(data[length:]))
		}
		c.TransferBuffer = make([]byte, c.ActualLength)
		copy(c.TransferBuffer, data[length:length+int(c.ActualLength)])
		length += int(c.ActualLength)
	}

	if c.NumberOfPackets != 0x00000000 && c.NumberOfPackets != 0xffffffff {
		if len(data[length:]) < int(c.NumberOfPackets)*ISO_PACKET_DESCRIPTOR_LENGTH {
			return 0, fmt.Errorf("data too short for ISOPacketDescriptors, need %d, got %d", int(c.NumberOfPackets)*ISO_PACKET_DESCRIPTOR_LENGTH, len(data[length:]))
		}
		c.ISOPacketDescriptors = make([]ISOPacketDescriptor, c.NumberOfPackets)
		for i := range c.ISOPacketDescriptors {
			isoPacketDescriptorLength, err := c.ISOPacketDescriptors[i].UnmarshalBinaryWithLength(data[length:])
			if err != nil {
				return 0, fmt.Errorf("unable to unmarshal ISOPacketDescriptor #%d: %w", i, err)
			}
			length += isoPacketDescriptorLength
		}
	}

	return length, nil
}

func (c *RetSubmit) UnmarshalBinary(data []byte) error {
	_, err := c.UnmarshalBinaryWithLength(data)

	return err
}

func (c *RetSubmit) Length() int {
	length := CMD_HEADER_LENGTH + RET_SUBMIT_STATIC_FIELDS_LENGTH
	if c.Direction == DIR_IN {
		length += len(c.TransferBuffer)
	}
	length += ISO_PACKET_DESCRIPTOR_LENGTH * len(c.ISOPacketDescriptors)
	return length
}
