package descriptor

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

type StandardDeviceDescriptor struct {
	// Numeric expression specifying the size of this descriptor.
	BLength uint8
	// Device descriptor type (assigned by USB).
	BDescriptorType DescriptorType
	// USB HID Specification Release
	BCDUSB uint16
	// Class code (assigned by USB). Note that the HID class is defined in the Interface descriptor.
	BDeviceClass uint8
	// Subclass code (assigned by USB). These codes are qualified by the value of the bDeviceClass field.
	BDeviceSubClass uint8
	// Protocol code. These codes are qualified by the value of the bDeviceSubClass field.
	BDeviceProtocol uint8
	// Maximum packet size for endpoint zero (only 8, 16, 32, or 64 are valid).
	BMaxPacketSize uint8
	// Vendor ID (assigned by USB). For this example we’ll use 0xFFFF.
	IDVendor uint16
	// Product ID (assigned by manufacturer).
	IDProduct uint16
	// Device release number (assigned by manufacturer).
	BCDDevice uint16
	// Index of String descriptor describing manufacturer.
	IManufacturer uint8
	// Index of string descriptor describing product.
	IProduct uint8
	// Index of String descriptor describing the device’s serial number.
	ISerialNumber uint8
	// Number of possible configurations.
	BNumConfigurations uint8
}

func (s *StandardDeviceDescriptor) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, STANDARD_DEVICE_DESCRIPTOR_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read standard device descriptor from stream: %w", err)
	}

	s.BLength = buf[0]
	s.BDescriptorType = DescriptorType(buf[1])
	s.BCDUSB = binary.LittleEndian.Uint16(buf[2:4])
	s.BDeviceClass = buf[4]
	s.BDeviceSubClass = buf[5]
	s.BDeviceProtocol = buf[6]
	s.BMaxPacketSize = buf[7]
	s.IDVendor = binary.LittleEndian.Uint16(buf[8:10])
	s.IDProduct = binary.LittleEndian.Uint16(buf[10:12])
	s.BCDDevice = binary.LittleEndian.Uint16(buf[12:14])
	s.IManufacturer = buf[14]
	s.IProduct = buf[15]
	s.ISerialNumber = buf[16]
	s.BNumConfigurations = buf[17]

	return nil
}

func (s *StandardDeviceDescriptor) Encode(writer io.Writer) error {
	buf := make([]byte, STANDARD_DEVICE_DESCRIPTOR_LENGTH)

	buf[0] = s.BLength
	buf[1] = uint8(s.BDescriptorType)
	binary.LittleEndian.PutUint16(buf[2:4], s.BCDUSB)
	buf[4] = s.BDeviceClass
	buf[5] = s.BDeviceSubClass
	buf[6] = s.BDeviceProtocol
	buf[7] = s.BMaxPacketSize
	binary.LittleEndian.PutUint16(buf[8:10], s.IDVendor)
	binary.LittleEndian.PutUint16(buf[10:12], s.IDProduct)
	binary.LittleEndian.PutUint16(buf[12:14], s.BCDDevice)
	buf[14] = s.IManufacturer
	buf[15] = s.IProduct
	buf[16] = s.ISerialNumber
	buf[17] = s.BNumConfigurations

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write standard device scriptor to stream: %w", err)
	}

	return nil
}
