package descriptor

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

type StandardConfigurationDescriptor struct {
	// Size of this descriptor in bytes.
	BLength uint8
	// Configuration (assigned by USB).
	BDescriptorType DescriptorType
	// Total length of data returned for this configuration.
	// Includes the combined length of all returned descriptors (configuration, interface, endpoint, and HID) returned for this configuration.
	// This value includes the HID descriptor but none of the other HID class descriptors (report or designator).
	WTotalLength uint16
	// Number of interfaces supported by this configuration.
	BNumInterfaces uint8
	// Value to use as an argument to Set Configuration to select this configuration.
	BConfigurationValue uint8
	// Index of string descriptor describing this configuration. In this case there is none.
	IConfiguration uint8
	// Configuration characteristics
	// (D7: Bus Powered, D6: Self Powered, D5: Remote Wakeup, D4..0: Reserved (reset to 0))
	BMAttributes uint8
	// Maximum power consumption of USB device from bus in this specific configuration when the device is fully operational. Expressed in 2 mA units—
	BMaxPower uint8
}

func (s *StandardConfigurationDescriptor) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, STANDARD_CONFIGURATION_DESCRIPTOR_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read standard configuration descriptor from stream: %w", err)
	}

	s.BLength = buf[0]
	s.BDescriptorType = DescriptorType(buf[1])
	s.WTotalLength = binary.LittleEndian.Uint16(buf[2:4])
	s.BNumInterfaces = buf[4]
	s.BConfigurationValue = buf[5]
	s.IConfiguration = buf[6]
	s.BMAttributes = buf[7]
	s.BMaxPower = buf[8]

	return nil
}

func (s *StandardConfigurationDescriptor) Encode(writer io.Writer) error {
	buf := make([]byte, STANDARD_CONFIGURATION_DESCRIPTOR_LENGTH)

	buf[0] = s.BLength
	buf[1] = uint8(s.BDescriptorType)
	binary.LittleEndian.PutUint16(buf[2:4], s.WTotalLength)
	buf[4] = s.BNumInterfaces
	buf[5] = s.BConfigurationValue
	buf[6] = s.IConfiguration
	buf[7] = s.BMAttributes
	buf[8] = s.BMaxPower

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write standard configuration scriptor to stream: %w", err)
	}

	return nil
}
