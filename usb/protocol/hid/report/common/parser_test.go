package common_test

import (
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report/common"
	"github.com/stretchr/testify/assert"
)

func TestByteToBool(t *testing.T) {
	res1 := common.ByteToBool(0)
	res2 := common.ByteToBool(1)
	res3 := common.ByteToBool(255)

	assert.False(t, res1)
	assert.True(t, res2)
	assert.True(t, res3)
}

func TestParseInt(t *testing.T) {
	tests := []struct {
		name     string
		item     []byte
		expected int32
	}{
		{
			name:     "Int8 positive",
			item:     []byte{0x7F},
			expected: 127,
		},
		{
			name:     "Int8 negative",
			item:     []byte{0x80},
			expected: -128,
		},
		{
			name:     "Int8 zero",
			item:     []byte{0x00},
			expected: 0,
		},
		{
			name:     "Int16 positive",
			item:     []byte{0xFF, 0x7F},
			expected: 32767,
		},
		{
			name:     "Int16 negative",
			item:     []byte{0x00, 0x80},
			expected: -32768,
		},
		{
			name:     "Int16 zero",
			item:     []byte{0x00, 0x00},
			expected: 0,
		},
		{
			name:     "Int32 positive",
			item:     []byte{0xFF, 0xFF, 0xFF, 0x7F},
			expected: 2147483647,
		},
		{
			name:     "Int32 negative",
			item:     []byte{0x00, 0x00, 0x00, 0x80},
			expected: -2147483648,
		},
		{
			name:     "Int32 zero",
			item:     []byte{0x00, 0x00, 0x00, 0x00},
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := common.ParseInt(test.item)

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestParseUint(t *testing.T) {
	tests := []struct {
		name     string
		item     []byte
		expected uint32
	}{
		{
			name:     "Uint8 positive",
			item:     []byte{0xFF},
			expected: 255,
		},
		{
			name:     "Uint8 zero",
			item:     []byte{0x00},
			expected: 0,
		},
		{
			name:     "Uint16 positive",
			item:     []byte{0xFF, 0xFF},
			expected: 65535,
		},
		{
			name:     "Uint16 zero",
			item:     []byte{0x00, 0x00},
			expected: 0,
		},
		{
			name:     "Uint32 positive",
			item:     []byte{0xFF, 0xFF, 0xFF, 0xFF},
			expected: 4294967295,
		},
		{
			name:     "Uint32 zero",
			item:     []byte{0x00, 0x00, 0x00, 0x00},
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := common.ParseUint(test.item)

			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestParseNibbleInt(t *testing.T) {
	tests := []struct {
		name     string
		item     byte
		expected int8
	}{
		{
			name:     "minimum",
			item:     0b1010_1000,
			expected: -8,
		},
		{
			name:     "maximum",
			item:     0b0000_0111,
			expected: 7,
		},
		{
			name:     "zero",
			item:     0b0000_0000,
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			actual := common.ParseNibbleInt(test.item)

			assert.Equal(t, test.expected, actual)
		})
	}
}
