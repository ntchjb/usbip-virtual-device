package command

type Command uint32

const (
	CMD_SUBMIT Command = 0x0000_0001
	RET_SUBMIT Command = 0x0000_0003
	CMD_UNLINK Command = 0x0000_0002
	RET_UNLINK Command = 0x0000_0004
)

type Direction uint32

const (
	DIR_OUT Direction = 0x0000_0000
	DIR_IN  Direction = 0x0000_0001
)

// Command header for all 4 commands
type CmdHeader struct {
	Command Command
	// sequential number that identifies requests and corresponding responses; incremented per connection
	SeqNum uint32
	// specifies a remote USB device uniquely instead of busnum and devnum; for client (request), this value is ((busnum << 16) | devnum); for server (response), this shall be set to 0
	DevID uint32
	// 0: USBIP_DIR_OUT
	//
	// 1: USBIP_DIR_IN
	//
	// only used by client, for server this shall be 0
	Direction Direction
	// endpoint number only used by client, for server this shall be 0; for UNLINK, this shall be 0
	EndpointNumber uint32
}

// ISOPacketDescriptor for ISO endpoint transfer
type ISOPacketDescriptor struct {
	Offset         uint32
	ExpectedLength uint32
	ActualLength   uint32
	Status         uint32
}

// Submit an URB
//
// 'command' shall be 0x00000001
type CmdSubmit struct {
	CmdHeader
	// possible values depend on the USBIP_URB transfer_flags. Refer to include/uapi/linux/usbip.h and USB Request Block (URB). Refer to usbip_pack_cmd_submit() and tweak_transfer_flags() in drivers/usb/usbip/ usbip_common.c.
	TransferFlags uint32
	// use URB transfer_buffer_length
	TransferBufferLength uint32
	// use URB start_frame; initial frame for ISO transfer; shall be set to 0 if not ISO transfer
	StartFrame uint32
	// number of ISO packets; shall be set to 0xffffffff if not ISO transfer
	NumberOfPackets uint32
	// maximum time for the request on the server-side host controller
	Interval uint32
	// data bytes for USB setup, filled with zeros if not used
	Setup [8]byte
	// TransferBuffer has variable length
	// direction OUT -> Length = len(TransferBufferLength),
	// direction IN -> Length 0
	// For ISO transfers the padding between each ISO packets is not transmitted.
	TransferBuffer       []byte
	ISOPacketDescriptors []ISOPacketDescriptor
}

// Reply for submitting an URB
//
// 'command' shall be 0x00000003
type RetSubmit struct {
	CmdHeader
	// zero for successful URB transaction, otherwise some kind of error happened.
	Status uint32
	// number of URB data bytes; use URB actual_length
	ActualLength uint32
	// use URB start_frame; initial frame for ISO transfer; shall be set to 0 if not ISO transfer
	StartFrame uint32
	// number of ISO packets; shall be set to 0xffffffff if not ISO transfer
	NumberOfPackets uint32
	ErrorCount      uint32
	// padding, shall be set to 0
	Padding uint64
	// direction IN -> Length = len(ActualLength),
	// direction OUT -> Length 0
	//
	// For ISO transfers the padding between each ISO packets is not transmitted.
	TransferBuffer       []byte
	ISOPacketDescriptors []ISOPacketDescriptor
}

// Unlink an URB
//
// 'command' shall be 0x00000002
type CmdUnlink struct {
	CmdHeader
	// UNLINK sequence number, of the SUBMIT request to unlink
	UnlinkSeqNum uint32
	// padding should be all zero
	Padding [24]byte
}

// Reply for URB unlink
//
// 'command' shall be 0x00000004
type RetUnlink struct {
	CmdHeader
	// This is similar to the status of USBIP_RET_SUBMIT (share the same memory offset). When UNLINK is successful, status is -ECONNRESET; when USBIP_CMD_UNLINK is after USBIP_RET_SUBMIT status is 0
	Status int32
	// padding, shall be set to 0
	Padding [24]byte
}

const (
	CMD_HEADER_LENGTH               = 20
	ISO_PACKET_DESCRIPTOR_LENGTH    = 16
	CMD_SUBMIT_STATIC_FIELDS_LENGTH = 28
	RET_SUBMIT_STATIC_FIELDS_LENGTH = 28
	CMD_UNLINK_STATIC_FIELDS_LENGTH = 28
	RET_UNLINK_STATIC_FIELDS_LENGTH = 28
)
