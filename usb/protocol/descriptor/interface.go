package descriptor

import (
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

type StandardInterfaceDescriptor struct {
	// Size of this descriptor in bytes.
	BLength uint8
	// Interface descriptor type (assigned by USB).
	BDescriptorType DescriptorType
	// Number of interface. Zero-based value identifying the index in the array of concurrent interfaces supported by this configuration.
	BInterfaceNumber uint8
	// Value used to select alternate setting for the interface identified in the prior field.
	BAlternateSetting uint8
	// Number of endpoints used by this interface (excluding endpoint zero). If this value is zero, this interface only uses endpoint zero.
	BNumEndpoints uint8
	// Class code (HID code: 3).
	BInterfaceClass uint8
	// Subclass code (HID: 0: No subclass, 1: Boot Interface subclass)
	BInterfaceSubClass uint8
	// Protocol code (HID: 0: None, 1: Keyboard, 2: Mouse)
	BInterfaceProtocol uint8
	// Index of string descriptor describing this interface.
	IInterface uint8
}

func (s *StandardInterfaceDescriptor) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, STANDARD_INTERFACE_DESCRIPTOR_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read standard interface descriptor from stream: %w", err)
	}

	s.BLength = buf[0]
	s.BDescriptorType = DescriptorType(buf[1])
	s.BInterfaceNumber = buf[2]
	s.BAlternateSetting = buf[3]
	s.BNumEndpoints = buf[4]
	s.BInterfaceClass = buf[5]
	s.BInterfaceSubClass = buf[6]
	s.BInterfaceProtocol = buf[7]
	s.IInterface = buf[8]

	return nil
}

func (s *StandardInterfaceDescriptor) Encode(writer io.Writer) error {
	buf := make([]byte, STANDARD_INTERFACE_DESCRIPTOR_LENGTH)

	buf[0] = s.BLength
	buf[1] = uint8(s.BDescriptorType)
	buf[2] = s.BInterfaceNumber
	buf[3] = s.BAlternateSetting
	buf[4] = s.BNumEndpoints
	buf[5] = s.BInterfaceClass
	buf[6] = s.BInterfaceSubClass
	buf[7] = s.BInterfaceProtocol
	buf[8] = s.IInterface

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write standard interface descriptor to stream: %w", err)
	}

	return nil
}
