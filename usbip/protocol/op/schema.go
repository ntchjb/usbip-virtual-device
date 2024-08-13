package op

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

const (
	OP_HEADER_LENGTH             = 8
	DEVICE_INFO_TRUNCATED_LENGTH = 312
	DEVICE_INTERFACE_LENGTH      = 4
)
