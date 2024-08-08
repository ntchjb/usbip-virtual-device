package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ntchjb/usbip-virtual-device/usb"
	"github.com/ntchjb/usbip-virtual-device/usbip"
)

func main() {
	gracefulStop := make(chan os.Signal, 1)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	mouse := NewGenericHIDMouseDevice(logger)
	deviceRegistrar := usb.NewDeviceRegistrar(usb.DeviceRegistrarConfig{
		BusNum:         1,
		MaxDeviceCount: 10,
	})
	if err := deviceRegistrar.Register(mouse); err != nil {
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
}
