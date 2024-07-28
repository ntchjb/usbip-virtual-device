package protocol

import "fmt"

func (op *OpReqImport) MarshalBinaryPreAlloc(data []byte) error {
	if len(data) < op.Length() {
		return fmt.Errorf("data too short to allocate for OpreqImport, need %d, got %d", op.Length(), len(data))
	}

	if err := op.OpHeader.MarshalBinaryPreAlloc(data[:OP_HEADER_LENGTH]); err != nil {
		return fmt.Errorf("unable to allocate OpHeader data : %w", err)
	}
	copy(data[OP_HEADER_LENGTH:], op.BusID[:])

	return nil
}

func (op *OpReqImport) MarshalBinary() (data []byte, err error) {
	dataLength := op.Length()
	data = make([]byte, dataLength)

	if err := op.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate data for OpReqImport: %w", err)
	}

	return data, nil
}

func (op *OpReqImport) UnmarshalBinaryWithLength(data []byte) (int, error) {
	var length int

	// OpHeader
	opHeaderLen, err := op.OpHeader.UnmarshalBinaryWithLength(data)
	if err != nil {
		return 0, fmt.Errorf("unable to unmarshal OpHeader: %w", err)
	}
	length += opHeaderLen

	if len(data[opHeaderLen:]) < len(op.BusID) {
		return 0, fmt.Errorf("data too short for BusID, need %d, got %d", len(op.BusID), len(data[8:]))
	}
	copy(op.BusID[:], data[opHeaderLen:])
	length += len(op.BusID)

	return length, nil
}

func (op *OpReqImport) UnmarshalBinary(data []byte) error {
	_, err := op.UnmarshalBinaryWithLength(data)

	return err
}

func (op *OpReqImport) Length() int {
	return op.OpHeader.Length() + len(op.BusID)
}

func (op *OpRepImport) MarshalBinaryPreAlloc(data []byte) error {
	if len(data) < op.Length() {
		return fmt.Errorf("data too short to allocate for OpRepImport, need %d, got %d", op.Length(), len(data))
	}

	if err := op.OpHeader.MarshalBinaryPreAlloc(data[:OP_HEADER_LENGTH]); err != nil {
		return fmt.Errorf("unable to allocate OpHeader data : %w", err)
	}

	// If the previous status field was OK (0), otherwise the reply ends with the status field.
	if op.Status == OP_STATUS_ERROR {
		return nil
	}

	if err := op.DeviceInfo.MarshalBinaryPreAlloc(data[OP_HEADER_LENGTH:]); err != nil {
		return fmt.Errorf("unable to allocate DeviceInfoTruncated: %w", err)
	}

	return nil
}

func (op *OpRepImport) MarshalBinary() (data []byte, err error) {
	dataLength := op.Length()
	data = make([]byte, dataLength)

	if err := op.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate data for OpRepImport: %w", err)
	}

	return data, nil
}

func (op *OpRepImport) UnmarshalBinaryWithLength(data []byte) (int, error) {
	var length int

	// OpHeader
	opHeaderLen, err := op.OpHeader.UnmarshalBinaryWithLength(data)
	if err != nil {
		return 0, fmt.Errorf("unable to unmarshal OpHeader: %w", err)
	}
	length += opHeaderLen

	// If the previous status field was OK (0), otherwise the reply ends with the status field.
	if op.Status == OP_STATUS_ERROR {
		return length, nil
	}

	deviceInfoLength, err := op.DeviceInfo.UnmarshalBinaryWithLength(data[opHeaderLen:])
	if err != nil {
		return 0, fmt.Errorf("unable to allocate DeviceInfoTruncated: %w", err)
	}
	length += deviceInfoLength

	return length, nil
}

func (op *OpRepImport) UnmarshalBinary(data []byte) error {
	_, err := op.UnmarshalBinaryWithLength(data)

	return err
}

func (op *OpRepImport) Length() int {
	return op.OpHeader.Length() + op.DeviceInfo.Length()
}
