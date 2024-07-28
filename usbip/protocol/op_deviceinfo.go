package protocol

import (
	"encoding/binary"
	"fmt"
)

func (op *DeviceInfoTruncated) UnmarshalBinaryWithLength(data []byte) (int, error) {
	if len(data) < DEVICE_INFO_TRUNCATED_LENGTH {
		return 0, fmt.Errorf("data too short for DeviceInfoTruncated, need %d, got %d", DEVICE_INFO_TRUNCATED_LENGTH, len(data))
	}
	copy(op.Path[:], data[:256])
	copy(op.BusID[:], data[256:288])
	op.BusNum = binary.BigEndian.Uint32(data[288:292])
	op.DevNum = binary.BigEndian.Uint32(data[292:296])
	op.Speed = binary.BigEndian.Uint32(data[296:300])
	op.IDVendor = binary.BigEndian.Uint16(data[300:302])
	op.IDProduct = binary.BigEndian.Uint16(data[302:304])
	op.BCDDevice = binary.BigEndian.Uint16(data[304:306])
	op.BDeviceClass = data[306]
	op.BDeviceSubclass = data[307]
	op.BDeviceProtocol = data[308]
	op.BConfigurationValue = data[309]
	op.BNumConfigurations = data[310]
	op.BNumInterfaces = data[311]

	return DEVICE_INFO_TRUNCATED_LENGTH, nil
}

func (op *DeviceInfoTruncated) UnmarshalBinary(data []byte) error {
	_, err := op.UnmarshalBinaryWithLength(data)

	return err
}

func (op *DeviceInfoTruncated) MarshalBinary() (data []byte, err error) {
	data = make([]byte, op.Length())

	if err := op.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate DeviceInfoTruncated: %w", err)
	}

	return data, nil
}

func (op *DeviceInfoTruncated) MarshalBinaryPreAlloc(data []byte) error {
	if len(data) < op.Length() {
		return fmt.Errorf("data too short to allocate for DeviceInfoTruncated, need %d, got %d", op.Length(), len(data))
	}

	copy(data[:256], op.Path[:])
	copy(data[256:288], op.BusID[:])
	binary.BigEndian.PutUint32(data[288:292], op.BusNum)
	binary.BigEndian.PutUint32(data[292:296], op.DevNum)
	binary.BigEndian.PutUint32(data[296:300], op.Speed)
	binary.BigEndian.PutUint16(data[300:302], op.IDVendor)
	binary.BigEndian.PutUint16(data[302:304], op.IDProduct)
	binary.BigEndian.PutUint16(data[304:306], op.BCDDevice)
	data[306] = op.BDeviceClass
	data[307] = op.BDeviceSubclass
	data[308] = op.BDeviceProtocol
	data[309] = op.BConfigurationValue
	data[310] = op.BNumConfigurations
	data[311] = op.BNumInterfaces

	return nil
}

func (op *DeviceInfoTruncated) Length() int {
	return DEVICE_INFO_TRUNCATED_LENGTH
}

func (op *DeviceInfo) UnmarshalBinaryWithLength(data []byte) (int, error) {
	var length int
	deviceInfoTruncatedLength, err := op.DeviceInfoTruncated.UnmarshalBinaryWithLength(data)
	if err != nil {
		return 0, fmt.Errorf("unable to unmarshal DeviceInfoTruncated: %w", err)
	}
	length += deviceInfoTruncatedLength
	for j := 0; j < int(op.BNumInterfaces); j++ {
		startIdx := deviceInfoTruncatedLength + j*4
		if len(data[startIdx:]) < DEVICE_INTERFACE_LENGTH {
			return 0, fmt.Errorf("data too short for DeviceInterface #%d, need %d, got %d", j+1, DEVICE_INTERFACE_LENGTH, len(data[startIdx:]))
		}
		intf := DeviceInterface{
			BInterfaceClass:    data[startIdx],
			BInterfaceSubclass: data[startIdx+1],
			BInterfaceProtocol: data[startIdx+2],
			PaddingAlignment:   data[startIdx+3],
		}
		op.Interfaces = append(op.Interfaces, intf)
		length += DEVICE_INTERFACE_LENGTH
	}

	return length, nil
}

func (op *DeviceInfo) UnmarshalBinary(data []byte) error {
	_, err := op.UnmarshalBinaryWithLength(data)

	return err
}

func (op *DeviceInfo) MarshalBinary() (data []byte, err error) {
	data = make([]byte, op.Length())

	if err := op.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate DeviceInfo data: %w", err)
	}

	return data, nil
}

func (op *DeviceInfo) MarshalBinaryPreAlloc(data []byte) error {
	if len(op.Interfaces) != int(op.BNumInterfaces) {
		return fmt.Errorf("device have mismatch number of interfaces; specified %d, array length is %d", op.BNumInterfaces, len(op.Interfaces))
	}
	if len(data) < op.Length() {
		return fmt.Errorf("data too short to allocate for DeviceInfo DeviceInfo, need %d, got %d", op.Length(), len(data))
	}

	if err := op.DeviceInfoTruncated.MarshalBinaryPreAlloc(data[:DEVICE_INTERFACE_LENGTH]); err != nil {
		return fmt.Errorf("unable to marshal DeviceInfoTruncated: %w", err)
	}
	for i, intf := range op.Interfaces {
		startIdx := DEVICE_INFO_TRUNCATED_LENGTH + i*DEVICE_INTERFACE_LENGTH
		data[startIdx] = intf.BInterfaceClass
		data[startIdx+1] = intf.BInterfaceSubclass
		data[startIdx+2] = intf.BInterfaceProtocol
		data[startIdx+3] = intf.PaddingAlignment
	}

	return nil
}

func (op *DeviceInfo) Length() int {
	return DEVICE_INFO_TRUNCATED_LENGTH + int(op.BNumInterfaces)*DEVICE_INTERFACE_LENGTH
}
