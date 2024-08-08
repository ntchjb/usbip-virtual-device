package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

type HIDDescriptor struct {
	// Numeric expression that is the total size of the HID descriptor.
	BLength uint8
	// Constant name specifying type of HID descriptor.
	BDescriptorType DescriptorType
	// Numeric expression identifying the HIDClass Specification release.
	BCDHID uint16
	// Numeric expression identifying country code of the localized hardware.
	BCountryCode uint8
	// Numeric expression specifying the number of class descriptors (always at least one i.e. Report descriptor.)
	BNumDescriptors uint8
	// Constant name identifying type of class descriptor. See Section 7.1.2: Set_Descriptor Request for a table of class descriptor constants.
	BClassDescriptorType DescriptorType
	// Numeric expression that is the total size of the Report descriptor.
	WDescriptorLength uint16
	// Constant name specifying type of optional descriptor.
	BOptionalDescriptorType DescriptorType
	// Numeric expression that is the total size of the optional descriptor.
	BOptionalDescriptorLength uint16
}

func (h *HIDDescriptor) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, HID_DESCRIPTOR_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read HID descriptor from stream: %w", err)
	}

	h.BLength = buf[0]
	h.BDescriptorType = DescriptorType(buf[1])
	h.BCDHID = binary.LittleEndian.Uint16(buf[2:4])
	h.BCountryCode = buf[4]
	h.BNumDescriptors = buf[5]
	h.BClassDescriptorType = DescriptorType(buf[6])
	h.WDescriptorLength = binary.LittleEndian.Uint16(buf[7:9])

	if h.BLength > HID_DESCRIPTOR_LENGTH {
		buf, err := stream.Read(reader, 3)
		if err != nil {
			return fmt.Errorf("unable to read optional HID descriptor type from stream: %w", err)
		}
		h.BOptionalDescriptorType = DescriptorType(buf[0])
		h.BOptionalDescriptorLength = binary.LittleEndian.Uint16(buf[1:3])
	}

	return nil
}

func (h *HIDDescriptor) Encode(writer io.Writer) error {
	buf := make([]byte, HID_DESCRIPTOR_LENGTH)

	buf[0] = h.BLength
	buf[1] = byte(h.BDescriptorType)
	binary.LittleEndian.PutUint16(buf[2:4], h.BCDHID)
	buf[4] = h.BCountryCode
	buf[5] = h.BNumDescriptors
	buf[6] = byte(h.BClassDescriptorType)
	binary.LittleEndian.PutUint16(buf[7:9], h.WDescriptorLength)

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write HID descriptor to stream: %w", err)
	}

	if h.BOptionalDescriptorType > 0 {
		buf := make([]byte, 3)
		buf[0] = byte(h.BOptionalDescriptorType)
		binary.LittleEndian.PutUint16(buf[1:3], h.BOptionalDescriptorLength)

		if err := stream.Write(writer, buf); err != nil {
			return fmt.Errorf("unable to write HID descriptor optional types to stream: %w", err)
		}
	}

	return nil
}
