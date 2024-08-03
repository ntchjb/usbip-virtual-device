package protocol

import "io"

type Serializer interface {
	Decode(reader io.Reader) error
	Encode(writer io.Writer) error
}

type Operation uint16

const (
	// Current USB/IP protocol version
	VERSION uint16 = 0x0111
)

const (
	OP_REQ_DEVLIST Operation = 0x8005
	OP_REP_DEVLIST Operation = 0x0005

	OP_REQ_IMPORT Operation = 0x8003
	OP_REP_IMPORT Operation = 0x0003
)

type OperationStatus uint32

const (
	OP_STATUS_OK    OperationStatus = 0x00000000
	OP_STATUS_ERROR OperationStatus = 0x00000001
)

type OpHeader struct {
	Version            uint16
	CommandOrReplyCode Operation
	Status             OperationStatus
}

// Retrieve the list of exported USB devices
//
// For Status field: unused, shall be set to 0
type OpReqDevList struct {
	OpHeader
}

// Reply with the list of exported USB devices
type OpRepDevList struct {
	OpHeader
	DeviceCount uint32
	Devices     []DeviceInfo
}

// Request to import (attach) a remote USB device
//
// For status field is unused, shall be set to 0
type OpReqImport struct {
	OpHeader
	// the busid of the exported device on the remote host. The possible values are taken from the message field OP_REP_DEVLIST.busid. A string closed with zero, the unused bytes shall be filled with zeros.
	BusID [32]byte
}

// Reply to import (attach) a remote USB device
type OpRepImport struct {
	OpHeader
	DeviceInfo DeviceInfoTruncated
}

// Device information for attachment, with no interface list
type DeviceInfoTruncated struct {
	// Path of the device on the host exporting the USB device, string closed with zero byte, e.g. “/sys/devices/pci0000:00/0000:00:1d.1/usb3/3-2” The unused bytes shall be filled with zero bytes.
	Path [256]byte
	// Bus ID of the exported device, string closed with zero byte, e.g. “3-2”. The unused bytes shall be filled with zero bytes.
	BusID [32]byte
	// Bus number of the device
	BusNum uint32
	// Device number of given bus
	DevNum uint32
	// Speed of this device
	Speed uint32
	// Vendor ID
	IDVendor uint16
	// Product ID
	IDProduct uint16
	// Device release number (assigned by manufacturer).
	BCDDevice uint16
	// Class code (assigned by USB). Note that the HID class is defined in the Interface descriptor.
	BDeviceClass uint8
	// Subclass code (assigned by USB). These codes are qualified by the value of the bDeviceClass field.
	BDeviceSubclass uint8
	// Protocol code. These codes are qualified by the value of the bDeviceSubClass field.
	BDeviceProtocol uint8
	// Value to use as an argument to Set Configuration to select this configuration.
	BConfigurationValue uint8
	// Number of possible configurations.
	BNumConfigurations uint8
	// Number of interfaces supported by this configuration.
	BNumInterfaces uint8
}

// Device information for attachment
type DeviceInfo struct {
	DeviceInfoTruncated
	Interfaces []DeviceInterface
}

// Device interface
type DeviceInterface struct {
	// Class code (assigned by the USB-IF).
	BInterfaceClass uint8
	// Subclass code (assigned by the USB-IF).
	BInterfaceSubclass uint8
	// Protocol code (assigned by the USB).
	BInterfaceProtocol uint8
	// padding byte for alignment, shall be set to zero
	PaddingAlignment uint8
}

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
	OP_HEADER_LENGTH                = 8
	DEVICE_INFO_TRUNCATED_LENGTH    = 312
	DEVICE_INTERFACE_LENGTH         = 4
	CMD_HEADER_LENGTH               = 20
	ISO_PACKET_DESCRIPTOR_LENGTH    = 16
	CMD_SUBMIT_STATIC_FIELDS_LENGTH = 28
	RET_SUBMIT_STATIC_FIELDS_LENGTH = 28
	CMD_UNLINK_STATIC_FIELDS_LENGTH = 28
	RET_UNLINK_STATIC_FIELDS_LENGTH = 28
)
