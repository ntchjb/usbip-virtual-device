package stream_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
	"github.com/stretchr/testify/assert"
)

var errCranky = errors.New("cranky error")

type crankyIO struct{}

func (c *crankyIO) Read(p []byte) (n int, err error) {
	return 0, errCranky
}

func (c *crankyIO) Write(p []byte) (n int, err error) {
	return 0, errCranky
}

func TestRead(t *testing.T) {
	tests := []struct {
		name         string
		reader       io.Reader
		bufferLength int
		out          []byte
		err          error
	}{
		{
			name: "Successfully read data",
			reader: bytes.NewBuffer([]byte{
				0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A,
			}),
			bufferLength: 5,
			out:          []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			err:          nil,
		},
		{
			name:         "End of file reached",
			reader:       bytes.NewBuffer([]byte{}),
			bufferLength: 5,
			out:          nil,
			err:          io.EOF,
		},
		{
			name:         "Unknown error occurred",
			reader:       &crankyIO{},
			bufferLength: 5,
			out:          nil,
			err:          errCranky,
		},
		{
			name:         "Incomplete data read",
			reader:       bytes.NewBuffer([]byte{0x00, 0x01, 0x02}),
			bufferLength: 5,
			out:          nil,
			err:          stream.ErrIncompleteReadData,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := stream.Read(test.reader, test.bufferLength)

			assert.Equal(t, test.out, res)
			assert.ErrorIs(t, err, test.err)
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name     string
		writer   io.Writer
		sinkData func(io.Writer) []byte
		in       []byte
		err      error
	}{
		{
			name:   "Successfully write data",
			writer: bytes.NewBuffer(nil),
			sinkData: func(w io.Writer) []byte {
				return w.(*bytes.Buffer).Bytes()
			},
			in:  []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			err: nil,
		},
		{
			name:   "Error occurred",
			writer: &crankyIO{},
			sinkData: func(w io.Writer) []byte {
				return nil
			},
			in:  nil,
			err: errCranky,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := stream.Write(test.writer, test.in)

			assert.Equal(t, test.in, test.sinkData(test.writer))
			assert.ErrorIs(t, err, test.err)
		})
	}
}
