package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

// Decode read data from stream and store in struct.
// Note that this function does not decode OpHeader, which should be done already during connection handling
func (op *DeviceInfoTruncated) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, DEVICE_INFO_TRUNCATED_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read stream for DeviceInfoTruncated: %w", err)
	}

	copy(op.Path[:], buf[:256])
	copy(op.BusID[:], buf[256:288])
	op.BusNum = binary.BigEndian.Uint32(buf[288:292])
	op.DevNum = binary.BigEndian.Uint32(buf[292:296])
	op.Speed = binary.BigEndian.Uint32(buf[296:300])
	op.IDVendor = binary.BigEndian.Uint16(buf[300:302])
	op.IDProduct = binary.BigEndian.Uint16(buf[302:304])
	op.BCDDevice = binary.BigEndian.Uint16(buf[304:306])
	op.BDeviceClass = buf[306]
	op.BDeviceSubclass = buf[307]
	op.BDeviceProtocol = buf[308]
	op.BConfigurationValue = buf[309]
	op.BNumConfigurations = buf[310]
	op.BNumInterfaces = buf[311]

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode OpHeader, which should be done already during connection handling
func (op *DeviceInfoTruncated) Encode(writer io.Writer) error {
	buf := make([]byte, DEVICE_INFO_TRUNCATED_LENGTH)

	copy(buf[:256], op.Path[:])
	copy(buf[256:288], op.BusID[:])
	binary.BigEndian.PutUint32(buf[288:292], op.BusNum)
	binary.BigEndian.PutUint32(buf[292:296], op.DevNum)
	binary.BigEndian.PutUint32(buf[296:300], op.Speed)
	binary.BigEndian.PutUint16(buf[300:302], op.IDVendor)
	binary.BigEndian.PutUint16(buf[302:304], op.IDProduct)
	binary.BigEndian.PutUint16(buf[304:306], op.BCDDevice)
	buf[306] = op.BDeviceClass
	buf[307] = op.BDeviceSubclass
	buf[308] = op.BDeviceProtocol
	buf[309] = op.BConfigurationValue
	buf[310] = op.BNumConfigurations
	buf[311] = op.BNumInterfaces

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write DeviceInfoTruncated to stream: %w", err)
	}

	return nil
}

// Decode read data from stream and store in struct.
// Note that this function does not decode OpHeader, which should be done already during connection handling
func (op *DeviceInterface) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, DEVICE_INTERFACE_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read DeviceInterface from stream: %w", err)
	}

	op.BInterfaceClass = buf[0]
	op.BInterfaceSubclass = buf[1]
	op.BInterfaceProtocol = buf[2]
	op.PaddingAlignment = buf[3]

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode OpHeader, which should be done already during connection handling
func (op *DeviceInterface) Encode(writer io.Writer) error {
	buf := make([]byte, DEVICE_INTERFACE_LENGTH)
	buf[0] = op.BInterfaceClass
	buf[1] = op.BInterfaceSubclass
	buf[2] = op.BInterfaceProtocol
	buf[3] = op.PaddingAlignment

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write DeviceInterface to stream: %w", err)
	}

	return nil
}

// Decode read data from stream and store in struct.
// Note that this function does not decode OpHeader, which should be done already during connection handling
func (op *DeviceInfo) Decode(reader io.Reader) error {
	if err := op.DeviceInfoTruncated.Decode(reader); err != nil {
		return fmt.Errorf("unable to decode DeviceInfoTruncated: %w", err)
	}

	op.Interfaces = make([]DeviceInterface, op.BNumInterfaces)
	for i := 0; i < int(op.BNumInterfaces); i++ {
		if err := op.Interfaces[i].Decode(reader); err != nil {
			return fmt.Errorf("unable to decode DeviceInterface: %w", err)
		}
	}

	return nil
}

// Encode writes data from struct to stream.
// Note that this function does not encode OpHeader, which should be done already during connection handling
func (op *DeviceInfo) Encode(writer io.Writer) error {
	if int(op.BNumInterfaces) != len(op.Interfaces) {
		return fmt.Errorf("expected number of interfaces does not match the actual, expected %d, actual %d", op.BNumInterfaces, len(op.Interfaces))
	}
	if err := op.DeviceInfoTruncated.Encode(writer); err != nil {
		return fmt.Errorf("unable to encode DeviceInfoTruncated: %w", err)
	}

	for _, intf := range op.Interfaces {
		if err := intf.Encode(writer); err != nil {
			return fmt.Errorf("unable to encode DeviceInterface: %w", err)
		}
	}

	return nil
}
