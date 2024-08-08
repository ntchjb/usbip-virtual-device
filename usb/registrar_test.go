package usb_test

import (
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb"
	"github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestRegistrar(t *testing.T) {
	ctrl := gomock.NewController(t)

	device1 := usb.NewMockDevice(ctrl)
	device2 := usb.NewMockDevice(ctrl)
	device3 := usb.NewMockDevice(ctrl)
	device1.EXPECT().SetBusID(uint(1), uint(1)).Return()
	device1.EXPECT().GetBusID().Return(protocol.BusID{0x01, 0x02, 0x03})
	device2.EXPECT().SetBusID(uint(1), uint(2)).Return()
	device2.EXPECT().GetBusID().Return(protocol.BusID{0x01, 0x02, 0x04})

	config := usb.DeviceRegistrarConfig{
		BusNum:         1,
		MaxDeviceCount: 2,
	}
	registrar := usb.NewDeviceRegistrar(config)

	err1 := registrar.Register(device1)
	err2 := registrar.Register(device2)
	err3 := registrar.Register(device3)

	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.ErrorIs(t, err3, usb.ErrMaximumDeviceCountReached)

	devices := registrar.GetAvailableDevices()
	assert.Equal(t, []usb.Device{device1, device2}, devices)

	actualDevice1, err := registrar.GetDevice(protocol.BusID{0x01, 0x02, 0x03})
	assert.NoError(t, err)
	assert.Equal(t, device1, actualDevice1)

	actualDevice2, err := registrar.GetDevice(protocol.BusID{0x01, 0x02, 0x04})
	assert.NoError(t, err)
	assert.Equal(t, device1, actualDevice2)

	actualDevice3, err := registrar.GetDevice(protocol.BusID{0x01, 0x02, 0x05})
	assert.ErrorIs(t, err, usb.ErrDeviceNotFound)
	assert.Equal(t, nil, actualDevice3)
}
