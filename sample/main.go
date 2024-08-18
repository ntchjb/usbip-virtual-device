package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ntchjb/usbip-virtual-device/sample/echo"
	"github.com/ntchjb/usbip-virtual-device/sample/mouse"
	"github.com/ntchjb/usbip-virtual-device/usb"
	"github.com/ntchjb/usbip-virtual-device/usbip"
)

func main() {
	gracefulStop := make(chan os.Signal, 1)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	device1 := mouse.NewGenericHIDMouseDevice(logger)
	device2 := echo.NewHIDEchoDevice(logger)
	deviceRegistrar := usb.NewDeviceRegistrar(usb.DeviceRegistrarConfig{
		BusNum:         1,
		MaxDeviceCount: 10,
	})
	if err := deviceRegistrar.Register(device1); err != nil {
		panic(err)
	}
	if err := deviceRegistrar.Register(device2); err != nil {
		panic(err)
	}

	server := usbip.NewUSBIPServer(usbip.USBIPServerConfig{
		ListenAddress:        "127.0.0.1:3240",
		TCPConnectionTimeout: 60 * time.Second,
		MaxTCPConnection:     10,
	}, deviceRegistrar, logger)

	if err := server.Open(); err != nil {
		panic(err)
	}

	logger.Info("Server is up")

	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	<-gracefulStop

	if err := server.Close(); err != nil {
		panic(err)
	}

	if err := deviceRegistrar.Close(); err != nil {
		panic(err)
	}
}
