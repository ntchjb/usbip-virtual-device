package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

type StringDescriptor struct {
	BLength         uint8
	BDescriptorType DescriptorType
	Content         []uint16
}

func (s *StringDescriptor) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, 1)
	if err != nil {
		return fmt.Errorf("unable to read string descriptor bLength from stream: %w", err)
	}
	s.BLength = buf[0]

	buf, err = stream.Read(reader, int(s.BLength))
	if err != nil {
		return fmt.Errorf("unable to read string descriptor for stream: %w", err)
	}
	s.BDescriptorType = DescriptorType(buf[0])

	s.Content = make([]uint16, (len(buf)-1)/2)

	for i := 1; i < len(buf); i += 2 {
		s.Content[i/2] = binary.LittleEndian.Uint16(buf[i : i+2])
	}

	return nil
}

func (s *StringDescriptor) Encode(writer io.Writer) error {
	buf := make([]byte, 2+len(s.Content)*2)
	buf[0] = s.BLength
	buf[1] = byte(s.BDescriptorType)
	startIdx := 2
	for _, con := range s.Content {
		binary.LittleEndian.PutUint16(buf[startIdx:startIdx+2], con)
		startIdx += 2
	}

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write string descriptor to stream: %w", err)
	}

	return nil
}
