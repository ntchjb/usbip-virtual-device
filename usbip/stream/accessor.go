package stream

import (
	"errors"
	"fmt"
	"io"
)

var (
	ErrIncompleteReadData  error = errors.New("incomplete data read from stream")
	ErrIncompleteWriteData error = errors.New("incomplete data write to stream")
)

// Read reads data from `reader` and use following conditions
//
// - If data is read partially (not fully filled `buf` array), then ErrIncompleteReadData is returned
// - If data is read but no data left, then io.EOF is returned
// - If error occurred, return wrapped error
// - Else, return buffer
func Read(reader io.Reader, bufferLength int) ([]byte, error) {
	buf := make([]byte, bufferLength)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("unable to read data from stream: %w", err)
	}
	if n == 0 {
		return nil, io.EOF
	}
	if n != len(buf) {
		return nil, ErrIncompleteReadData
	}

	return buf, nil
}

// Write writes data to given stream.
//
// - If error occurred, return wrapped error
// - If data is written partially, return incomplete data error
// - Else, return nil
func Write(writer io.Writer, buf []byte) error {
	n, err := writer.Write(buf)
	if err != nil {
		return fmt.Errorf("unable to write data to stream: %w", err)
	}
	if n < len(buf) {
		// This should be unreacheable because writer should return error if n < len(buf)
		return ErrIncompleteWriteData
	}

	return nil
}
