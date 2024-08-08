package usb

import (
	usbprotocol "github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
)

type WorkerPoolProfile struct {
	MaximumProcWorkers  int
	MaximumReplyWorkers int
}

// Device represents a USB device and its logic
type Device interface {
	// SetBusID assigns bus ID to this device. BusID should be assigned by device registrar
	SetBusID(busNum, devNum uint)
	// GetBusID returns bus ID of this device
	GetBusID() usbprotocol.BusID
	// GetDeviceInfo returns device information used by OpDevList
	GetDeviceInfo() protocol.DeviceInfo
	// GetURBProcessor returns an instance of processor of this device, used by handler's worker pool
	Process(data protocol.CmdSubmit) protocol.RetSubmit
	// GetWorkerPoolProfile indicates how worker pool behave for this device, such as, set worker count to 1 to process URB requests in sequences
	GetWorkerPoolProfile() WorkerPoolProfile
}