package protocol

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/ntchjb/usbip-virtual-device/usbip/stream"
)

type StandardEndpointDescriptor struct {
	// Size of this descriptor in bytes.
	BLength uint8
	// Endpoint descriptor type (assigned by USB).
	BDescriptorType DescriptorType
	// The address of the endpoint on the USB device described by this descriptor.
	// (D0..3: The endpoint number, D4..6: Reserved, reset to zero, D7: Direction (0: OUT endpoint, 1: IN endpoint, ignored: control endpoint))
	BEndpointAddress uint8
	// This field describes the endpoint’s attributes when it is configured using the bConfigurationValue.
	// (D0..1: Transfer type (00: Control, 01: Isochronous, 10: Bulk, 11: Interrupt))
	BMAttributes uint8
	// Maximum packet size this endpoint is capable of sending or receiving when this configuration is selected.
	// For interrupt endpoints, this value is used to reserve the bus time in the schedule, required for the per frame data payloads.
	// Smaller data payloads may be sent, but will terminate the transfer and thus require intervention to restart.
	WMaxPacketSize uint16
	// Interval for polling endpoint for data transfers. Expressed in milliseconds.
	BInterval uint8
}

func (s *StandardEndpointDescriptor) Decode(reader io.Reader) error {
	buf, err := stream.Read(reader, STANDARD_ENDPOINT_DESCRIPTOR_LENGTH)
	if err != nil {
		return fmt.Errorf("unable to read standard endpoint descriptor from stream: %w", err)
	}

	s.BLength = buf[0]
	s.BDescriptorType = DescriptorType(buf[1])
	s.BEndpointAddress = buf[2]
	s.BMAttributes = buf[3]
	s.WMaxPacketSize = binary.LittleEndian.Uint16(buf[4:6])
	s.BInterval = buf[6]

	return nil
}

func (s *StandardEndpointDescriptor) Encode(writer io.Writer) error {
	buf := make([]byte, STANDARD_ENDPOINT_DESCRIPTOR_LENGTH)

	buf[0] = s.BLength
	buf[1] = byte(s.BDescriptorType)
	buf[2] = s.BEndpointAddress
	buf[3] = s.BMAttributes
	binary.LittleEndian.PutUint16(buf[4:6], s.WMaxPacketSize)
	buf[6] = s.BInterval

	if err := stream.Write(writer, buf); err != nil {
		return fmt.Errorf("unable to write standard endpoint descriptor to stream: %w", err)
	}

	return nil
}
