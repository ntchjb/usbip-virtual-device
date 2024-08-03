package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

type SetupPacket struct {
	BMRequestType SetupRequestType
	BRequest      SetupRequest
	WValue        uint16
	WIndex        uint16
	WLength       uint16
}

func (s *SetupPacket) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, SETUP_PACKET_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read SetupPacket from stream: %w", err)
	}

	s.BMRequestType = SetupRequestType(buf[0])
	s.BRequest = SetupRequest(buf[1])
	s.WValue = binary.LittleEndian.Uint16(buf[2:4])
	s.WIndex = binary.LittleEndian.Uint16(buf[4:6])
	s.WLength = binary.LittleEndian.Uint16(buf[6:8])

	return nil
}

func (s *SetupPacket) Encode(writer io.Writer) error {
	buf := make([]byte, 8)

	buf[0] = byte(s.BMRequestType)
	buf[1] = byte(s.BRequest)
	binary.LittleEndian.PutUint16(buf[2:4], s.WValue)
	binary.LittleEndian.PutUint16(buf[4:6], s.WIndex)
	binary.LittleEndian.PutUint16(buf[6:8], s.WLength)

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write SetupPacket to stream: %w", err)
	}

	return nil
}
