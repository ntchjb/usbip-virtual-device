package usb

import "github.com/ntchjb/usbip-virtual-device/usbip/protocol"

type BusID [32]byte

type URBProcessor interface {
	ProcessSubmit(data protocol.CmdSubmit) protocol.RetSubmit
}

type Device interface {
	GetBusID() BusID
	GetDeviceInfo() protocol.DeviceInfo
	GetURBProcessor() URBProcessor
}
