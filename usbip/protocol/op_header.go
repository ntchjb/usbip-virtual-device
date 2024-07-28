package protocol

import (
	"encoding/binary"
	"fmt"
)

func (op *OpHeader) UnmarshalBinaryWithLength(data []byte) (int, error) {
	if len(data) < OP_HEADER_LENGTH {
		return 0, fmt.Errorf("unable to unmarshal OpHeader: need data length at least %d, got %d", OP_HEADER_LENGTH, len(data))
	}

	op.Version = binary.BigEndian.Uint16(data[:2])
	op.CommandOrReplyCode = Operation(binary.BigEndian.Uint16(data[2:4]))
	op.Status = OperationStatus(binary.BigEndian.Uint32(data[4:8]))

	return OP_HEADER_LENGTH, nil
}

func (op *OpHeader) UnmarshalBinary(data []byte) error {
	_, err := op.UnmarshalBinaryWithLength(data)

	return err
}

func (op *OpHeader) MarshalBinary() (data []byte, err error) {
	data = make([]byte, op.Length())

	if err := op.MarshalBinaryPreAlloc(data); err != nil {
		return nil, fmt.Errorf("unable to allocate data to OpHeader: %w", err)
	}

	return data, nil
}

func (op *OpHeader) MarshalBinaryPreAlloc(data []byte) (err error) {
	if len(data) < op.Length() {
		return fmt.Errorf("data too short to allocate OpHeader; expected %d, got %d", op.Length(), len(data))
	}
	binary.BigEndian.PutUint16(data[:2], op.Version)
	binary.BigEndian.PutUint16(data[2:4], uint16(op.CommandOrReplyCode))
	binary.BigEndian.PutUint32(data[4:8], uint32(op.Status))

	return nil
}

func (op *OpHeader) Length() int {
	return OP_HEADER_LENGTH
}
