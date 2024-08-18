package usb

import (
	"errors"
	"fmt"

	usbprotocol "github.com/ntchjb/usbip-virtual-device/usb/protocol"
)

var (
	ErrDeviceNotFound            = errors.New("USB device not found")
	ErrMaximumDeviceCountReached = errors.New("maximum number of registered device reached")
)

type DeviceRegistrar interface {
	// Register a USB device and assign new BusID/Path to it
	Register(device Device) error
	// Get a USB device by bus ID
	GetDevice(busID usbprotocol.BusID) (Device, error)
	// Get all registered devices
	GetAvailableDevices() []Device
	// Close all registered devices
	Close() error
}

type DeviceRegistrarConfig struct {
	// BusNum is used for generating BusID
	BusNum         uint
	MaxDeviceCount int
}

type deviceRegistrarImpl struct {
	devices map[usbprotocol.BusID]Device
	config  DeviceRegistrarConfig
	// currentDevNum is used for generating BusID
	currentDevNum uint
}

func NewDeviceRegistrar(config DeviceRegistrarConfig) DeviceRegistrar {
	return &deviceRegistrarImpl{
		devices: make(map[usbprotocol.BusID]Device),
		config:  config,
	}
}

func (r *deviceRegistrarImpl) createNewBusID() (uint, uint) {
	r.currentDevNum++

	return r.config.BusNum, r.currentDevNum
}

func (r *deviceRegistrarImpl) Register(device Device) error {
	if len(r.devices) >= r.config.MaxDeviceCount {
		return ErrMaximumDeviceCountReached
	}
	busNum, devNum := r.createNewBusID()
	device.SetBusID(busNum, devNum)
	r.devices[device.GetBusID()] = device

	return nil
}

func (r *deviceRegistrarImpl) GetDevice(busID usbprotocol.BusID) (Device, error) {
	if device, ok := r.devices[busID]; !ok {
		return nil, ErrDeviceNotFound
	} else {
		return device, nil
	}
}

func (r *deviceRegistrarImpl) GetAvailableDevices() []Device {
	if len(r.devices) == 0 {
		return nil
	}
	devices := make([]Device, len(r.devices))
	i := 0
	for _, device := range r.devices {
		devices[i] = device
		i++
	}

	return devices
}

func (r *deviceRegistrarImpl) Close() (err error) {
	for _, device := range r.devices {
		if deviceErr := device.Close(); deviceErr != nil {
			err = fmt.Errorf("unable to close device %s: %w", device.GetBusID(), err)
		}
	}

	return
}
