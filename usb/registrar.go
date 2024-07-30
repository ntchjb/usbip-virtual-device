package usb

import (
	"errors"
)

var (
	ErrDeviceNotFound            = errors.New("USB device not found")
	ErrMaximumDeviceCountReached = errors.New("maximum number of registered device reached")
)

const (
	MAX_DEVICE_COUNT = 1_000_000
)

type DeviceRegistrar interface {
	Register(device Device) error
	GetDevice(busID BusID) (Device, error)
	GetAvailableDevices() []Device
}

type deviceRegistrarImpl struct {
	devices map[BusID]Device
}

func NewDeviceRegistrar() DeviceRegistrar {
	return &deviceRegistrarImpl{
		devices: make(map[BusID]Device),
	}
}

func (r *deviceRegistrarImpl) Register(device Device) error {
	if len(r.devices) >= MAX_DEVICE_COUNT {
		return ErrMaximumDeviceCountReached
	}
	r.devices[device.GetBusID()] = device

	return nil
}

func (r *deviceRegistrarImpl) GetDevice(busID BusID) (Device, error) {
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
