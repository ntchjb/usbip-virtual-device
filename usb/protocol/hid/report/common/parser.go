package common

import "encoding/binary"

func ByteToBool(b byte) bool {
	return b != 0
}

func ParseUint(item []byte) uint32 {
	if len(item) == 4 {
		return binary.LittleEndian.Uint32(item[:4])
	} else if len(item) == 2 {
		return uint32(binary.LittleEndian.Uint16(item[:2]))
	} else {
		return uint32(item[0])
	}
}

func ParseInt(item []byte) int32 {
	if len(item) == 4 {
		return int32(binary.LittleEndian.Uint32(item[:4]))
	} else if len(item) == 2 {
		return int32(int16(binary.LittleEndian.Uint16(item[:2])))
	} else {
		return int32(int8(item[0]))
	}
}

// ParseNibbleInt converts a byte to int4
// considering 2s-complement format
// 1. Get the MSB digit, dudcted from 0 to be either 0b1111_1111 or 0b0000_0000
// 2. AND with 0b1111_0000 to get 4 most significant bits
// 3. OR with given nibble, so that we replace 4 most significant bits with either 1s or 0s based on 4th digit
func ParseNibbleInt(nibble byte) int8 {
	nibble = nibble & 0b0000_1111
	return int8(((0 - nibble>>3) & 0b1111_0000) | nibble)
}
