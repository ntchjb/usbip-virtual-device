package command

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

// Decode read data from stream and store in struct.
// Note that this function does not decode CmdHeader, which should be done already during connection handling
func (c *CmdSubmit) Decode(reader io.Reader) error {
	staticFieldBuf, err := stream.Read(reader, CMD_SUBMIT_STATIC_FIELDS_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read CmdSubmit static fields from stream: %w", err)
	}

	c.TransferFlags = binary.BigEndian.Uint32(staticFieldBuf[:4])
	c.TransferBufferLength = binary.BigEndian.Uint32(staticFieldBuf[4:8])
	c.StartFrame = binary.BigEndian.Uint32(staticFieldBuf[8:12])
	c.NumberOfPackets = binary.BigEndian.Uint32(staticFieldBuf[12:16])
	c.Interval = binary.BigEndian.Uint32(staticFieldBuf[16:20])
	copy(c.Setup[:], staticFieldBuf[20:28])

	if c.TransferBufferLength > 0 && c.Direction == DIR_OUT {
		transferBuf, err := stream.Read(reader, int(c.TransferBufferLength))
		if err != nil {
			return fmt.Errorf("unable to read TransferBuffer from stream: %w", err)
		}
		c.TransferBuffer = transferBuf
	}

	if c.NumberOfPackets != 0x00000000 && c.NumberOfPackets != 0xffffffff {
		c.ISOPacketDescriptors = make([]ISOPacketDescriptor, c.NumberOfPackets)
		for i := range c.ISOPacketDescriptors {
			if err := c.ISOPacketDescriptors[i].Decode(reader); err != nil {
				return fmt.Errorf("unable to decode ISOPacketDescriptor #%d: %w", i, err)
			}
		}
	}

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode CmdHeader, which should be done already during connection handling
func (c *CmdSubmit) Encode(writer io.Writer) error {
	if c.Direction == DIR_OUT && len(c.TransferBuffer) != int(c.TransferBufferLength) {
		return fmt.Errorf("actual transfer buffer length does not match TransferBufferLength for DIR_OUT CmdSubmit; expected %d, actual %d", c.TransferBufferLength, len(c.TransferBuffer))
	}
	if (c.NumberOfPackets == 0x00000000 || c.NumberOfPackets == 0xffffffff) && len(c.ISOPacketDescriptors) > 0 {
		return fmt.Errorf("this is non-ISO transfer, but contains ISO packet descriptors")
	} else if c.NumberOfPackets != 0x00000000 && c.NumberOfPackets != 0xffffffff && int(c.NumberOfPackets) != len(c.ISOPacketDescriptors) {
		return fmt.Errorf("number of packets does not match with actual ISO packet descriptor length, expected %d, actual %d", c.NumberOfPackets, len(c.ISOPacketDescriptors))
	}

	buf := make([]byte, CMD_SUBMIT_STATIC_FIELDS_LENGTH)
	binary.BigEndian.PutUint32(buf[:4], c.TransferFlags)
	binary.BigEndian.PutUint32(buf[4:8], c.TransferBufferLength)
	binary.BigEndian.PutUint32(buf[8:12], c.StartFrame)
	binary.BigEndian.PutUint32(buf[12:16], c.NumberOfPackets)
	binary.BigEndian.PutUint32(buf[16:20], c.Interval)
	copy(buf[20:28], c.Setup[:])

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write CmdSubmit static fields to stream: %w", err)
	}

	if len(c.TransferBuffer) > 0 {
		if err := stream.Write(writer, c.TransferBuffer); err != nil {
			return fmt.Errorf("unable to write TransferBuffer to stream: %w", err)
		}
	}

	for _, descriptor := range c.ISOPacketDescriptors {
		if err := descriptor.Encode(writer); err != nil {
			return fmt.Errorf("unable to write ISOPacketDescriptor to stream: %w", err)
		}
	}

	return nil
}

// Decode read data from stream and store in struct.
// Note that this function does not decode CmdHeader, which should be done already during connection handling
func (c *RetSubmit) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, RET_SUBMIT_STATIC_FIELDS_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read RetSubmit from stream: %w", err)
	}

	c.Status = binary.BigEndian.Uint32(buf[:4])
	c.ActualLength = binary.BigEndian.Uint32(buf[4:8])
	c.StartFrame = binary.BigEndian.Uint32(buf[8:12])
	c.NumberOfPackets = binary.BigEndian.Uint32(buf[12:16])
	c.ErrorCount = binary.BigEndian.Uint32(buf[16:20])
	c.Padding = binary.BigEndian.Uint64(buf[20:28])

	if c.ActualLength > 0 && c.Direction == DIR_IN {
		buf, err := stream.Read(reader, int(c.ActualLength))
		if err != nil {
			return fmt.Errorf("unable to read TransferBuffer from stream: %w", err)
		}
		c.TransferBuffer = buf
	}

	if c.NumberOfPackets != 0x00000000 && c.NumberOfPackets != 0xffffffff {
		c.ISOPacketDescriptors = make([]ISOPacketDescriptor, c.NumberOfPackets)
		for i := range c.ISOPacketDescriptors {
			if err := c.ISOPacketDescriptors[i].Decode(reader); err != nil {
				return fmt.Errorf("unable to decode ISOPacketDescriptor #%d: %w", i, err)
			}
		}
	}

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode CmdHeader, which should be done already during connection handling
func (c *RetSubmit) Encode(writer io.Writer) error {
	if c.Direction == DIR_IN && len(c.TransferBuffer) != int(c.ActualLength) {
		return fmt.Errorf("actual transfer buffer length does not match ActualLength for DIR_IN RetSubmit; expected %d, actual %d", c.ActualLength, len(c.TransferBuffer))
	}
	if (c.NumberOfPackets == 0x00000000 || c.NumberOfPackets == 0xffffffff) && len(c.ISOPacketDescriptors) > 0 {
		return fmt.Errorf("this is non-ISO transfer, but contains ISO packet descriptors")
	} else if c.NumberOfPackets != 0x00000000 && c.NumberOfPackets != 0xffffffff && int(c.NumberOfPackets) != len(c.ISOPacketDescriptors) {
		return fmt.Errorf("number of packets does not match with actual ISO packet descriptor length, expected %d, actual %d", c.NumberOfPackets, len(c.ISOPacketDescriptors))
	}

	buf := make([]byte, RET_SUBMIT_STATIC_FIELDS_LENGTH)

	binary.BigEndian.PutUint32(buf[:4], c.Status)
	binary.BigEndian.PutUint32(buf[4:8], c.ActualLength)
	binary.BigEndian.PutUint32(buf[8:12], c.StartFrame)
	binary.BigEndian.PutUint32(buf[12:16], c.NumberOfPackets)
	binary.BigEndian.PutUint32(buf[16:20], c.ErrorCount)
	binary.BigEndian.PutUint64(buf[20:28], c.Padding)

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write RetSubmit static field to stream: %w", err)
	}

	if c.ActualLength > 0 {
		if err := stream.Write(writer, c.TransferBuffer); err != nil {
			return fmt.Errorf("unable to write TransferBuffer to stream: %w", err)
		}
	}

	for _, descriptor := range c.ISOPacketDescriptors {
		if err := descriptor.Encode(writer); err != nil {
			return fmt.Errorf("unable to encode ISOPacketDescriptor: %w", err)
		}
	}

	return nil
}
