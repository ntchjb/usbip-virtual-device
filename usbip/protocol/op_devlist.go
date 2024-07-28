package protocol

import (
	"encoding/binary"
	"fmt"
)

func (op *OpRepDevList) MarshalBinaryPreAlloc(data []byte) error {
	if len(data) < op.Length() {
		return fmt.Errorf("data too short to allocate for OpRepDevList, need %d, got %d", op.Length(), len(data))
	}

	if err := op.OpHeader.MarshalBinaryPreAlloc(data[:OP_HEADER_LENGTH]); err != nil {
		return fmt.Errorf("unable to allocate OpHeader: %w", err)
	}
	binary.BigEndian.PutUint32(data[OP_HEADER_LENGTH:OP_HEADER_LENGTH+4], op.DeviceCount)

	var prevBNumInterfaces uint8
	for i, device := range op.Devices {
		startIdx := OP_HEADER_LENGTH + 4 + i*DEVICE_INFO_TRUNCATED_LENGTH + int(prevBNumInterfaces)*DEVICE_INTERFACE_LENGTH
		if err := device.MarshalBinaryPreAlloc(data[startIdx : startIdx+DEVICE_INFO_TRUNCATED_LENGTH+int(device.BNumInterfaces)*DEVICE_INTERFACE_LENGTH]); err != nil {
			return fmt.Errorf("unable to allocate DeviceInfo #%d: %w", i, err)
		}
		prevBNumInterfaces = device.BNumInterfaces
	}

	return nil
}

func (op *OpRepDevList) MarshalBinary() (data []byte, err error) {
	dataLength := op.Length()
	data = make([]byte, dataLength)

	if err := op.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate data for OpRepDevList: %w", err)
	}

	return data, nil
}

func (op *OpRepDevList) UnmarshalBinaryWithLength(data []byte) (int, error) {
	var length int

	// OpHeader
	opHeaderLen, err := op.OpHeader.UnmarshalBinaryWithLength(data)
	if err != nil {
		return 0, fmt.Errorf("unable to unmarshal OpHeader: %w", err)
	}
	length += opHeaderLen

	// DeviceCount
	if len(data[opHeaderLen:]) < 4 {
		return 0, fmt.Errorf("data too short for DeviceCount, need 4, got %d", len(data[8:]))
	}
	op.DeviceCount = binary.BigEndian.Uint32(data[opHeaderLen:])
	length += 4

	// Devices
	for i, startIdx := 0, length; i < int(op.DeviceCount); i++ {
		device := DeviceInfo{}
		deviceInfoLength, err := device.UnmarshalBinaryWithLength(data[startIdx:])
		if err != nil {
			return 0, fmt.Errorf("unable to unmarshal DeviceInfo #%d: %w", i, err)
		}
		startIdx += deviceInfoLength
		length += deviceInfoLength
		op.Devices = append(op.Devices, device)
	}

	return length, nil
}

func (op *OpRepDevList) UnmarshalBinary(data []byte) error {
	_, err := op.UnmarshalBinaryWithLength(data)

	return err
}

func (op *OpRepDevList) Length() int {
	dataLength := op.OpHeader.Length()
	for _, device := range op.Devices {
		dataLength += DEVICE_INFO_TRUNCATED_LENGTH + int(device.BNumInterfaces)*DEVICE_INTERFACE_LENGTH
	}

	return dataLength
}
