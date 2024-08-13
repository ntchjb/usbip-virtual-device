package common

import "encoding/binary"

func ByteToBool(b byte) bool {
	return b != 0
}

func ParseUint(item []byte) uint32 {
	if len(item) == 1 {
		return uint32(item[0])
	} else if len(item) > 1 {
		return binary.LittleEndian.Uint32(item[0:2])
	}

	return 0
}

func ParseInt(item []byte) int32 {
	if len(item) >= 4 {
		return int32(binary.LittleEndian.Uint32(item[:4]))
	} else if len(item) >= 2 {
		return int32(int16(binary.LittleEndian.Uint16(item[:2])))
	} else {
		return int32(int8(item[0]))
	}
}
