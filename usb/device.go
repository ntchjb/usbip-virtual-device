package usb

import (
	usbprotocol "github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol/command"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol/op"
)

type WorkerPoolProfile struct {
	// Maximum number of goroutines to process incoming URBs in parallel.
	// Set this to 1 to process all incoming URBs in sequence.
	MaximumProcWorkers int
	// Maximum number of goroutines to reply return data in parallel.
	// Set this to 1 to reply return data one-by-one in sequence.
	MaximumReplyWorkers int
	// Maximum number of goroutines to reply unlink data in parallel.
	// Set this to 1 to reply return data one-by-one in sequence.
	MaximumUnlinkReplyWorkers int
}

// Device represents a USB device and its logic
type Device interface {
	// SetBusID assigns bus ID to this device. BusID should be assigned by device registrar
	SetBusID(busNum, devNum uint)
	// GetBusID returns bus ID of this device
	GetBusID() usbprotocol.BusID
	// GetDeviceInfo returns device information used by OpDevList
	GetDeviceInfo() op.DeviceInfo
	// GetURBProcessor returns an instance of processor of this device, used by handler's worker pool
	Process(data command.CmdSubmit) command.RetSubmit
	// GetWorkerPoolProfile indicates how worker pool behave for this device, such as, set worker count to 1 to process URB requests in sequences
	GetWorkerPoolProfile() WorkerPoolProfile
}
