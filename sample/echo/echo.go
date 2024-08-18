package echo

import (
	"bytes"
	"fmt"
	"log/slog"
	"strings"
	"unicode/utf16"

	"github.com/ntchjb/usbip-virtual-device/usb"
	usbprotocol "github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/ntchjb/usbip-virtual-device/usb/protocol/descriptor"
	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid"
	"github.com/ntchjb/usbip-virtual-device/usb/protocol/hid/report"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol/command"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol/op"
)

var (
	echoHIDReport = []byte{
		byte(report.HID_REPORT_TAG_USAGE_PAGE) | 0x02, 0xA0, 0xFF, // Usage Page (0xFFA0)
		byte(report.HID_REPORT_TAG_USAGE) | 0x01, 0x01, // Usage (0x01)
		byte(report.HID_REPORT_TAG_COLLECTION) | 0x01, 0x01, // Collection (Application)
		byte(report.HID_REPORT_TAG_USAGE) | 0x01, 0x03, // Usage (0x03)
		byte(report.HID_REPORT_TAG_LOGICAL_MINIMUM) | 0x01, 0x00, // Logical Minimum (0x00)
		byte(report.HID_REPORT_TAG_LOGICAL_MAXIMUM) | 0x02, 0xFF, 0x00, // Logical Maximum (0xFF)
		byte(report.HID_REPORT_TAG_REPORT_SIZE) | 0x01, 0x08, // Report Size (8)
		byte(report.HID_REPORT_TAG_REPORT_COUNT) | 0x01, 0x40, // Report Count (64)
		byte(report.HID_REPORT_TAG_INPUT) | 0x01, 0b0000_1000, // Input (Data,Array,Absolute,Wrap)
		byte(report.HID_REPORT_TAG_USAGE) | 0x01, 0x04, // Usage (0x04)
		byte(report.HID_REPORT_TAG_LOGICAL_MINIMUM) | 0x01, 0x00, // Logical Minimum (0x00)
		byte(report.HID_REPORT_TAG_LOGICAL_MAXIMUM) | 0x02, 0xFF, 0x00, // Logical Maximum (0xFF)
		byte(report.HID_REPORT_TAG_REPORT_SIZE) | 0x01, 0x08, // Report Size (8)
		byte(report.HID_REPORT_TAG_REPORT_COUNT) | 0x01, 0x40, // Report Count (64)
		byte(report.HID_REPORT_TAG_OUTPUT) | 0x01, 0b0000_1000, // Output (ata,Array,Absolute,Wrap)
		byte(report.HID_REPORT_TAG_END_COLLECTION), // End Collection
	}
)

type genericHIDEchoDevice struct {
	deviceInfo op.DeviceInfo
	logger     *slog.Logger

	// This echo content does not guarantee to echo string in the same order as received strings
	// because we use WorkerPool with more than 1 goroutines, so URB commands are received un-ordered
	echoContent chan string
}

func NewHIDEchoDevice(logger *slog.Logger) usb.Device {
	return &genericHIDEchoDevice{
		logger: logger,
		deviceInfo: op.DeviceInfo{
			DeviceInfoTruncated: op.DeviceInfoTruncated{
				Speed:               usbprotocol.SPEED_USB2_HIGH,
				IDVendor:            0xecc0,
				IDProduct:           0x0001,
				BCDDevice:           1,
				BDeviceClass:        usbprotocol.CLASS_BASEDON_INTERFACE,
				BDeviceSubclass:     usbprotocol.SUBCLASS_NONE,
				BConfigurationValue: 1,
				BNumConfigurations:  1,
				BNumInterfaces:      1,
			},
			Interfaces: []op.DeviceInterface{
				{
					BInterfaceClass:    usbprotocol.CLASS_HID,
					BInterfaceSubclass: usbprotocol.SUBCLASS_NONE,
					BInterfaceProtocol: usbprotocol.PROTOCOL_NONE,
				},
			},
		},
		echoContent: make(chan string, 128),
	}
}

func (g *genericHIDEchoDevice) GetWorkerPoolProfile() usb.WorkerPoolProfile {
	return usb.WorkerPoolProfile{
		MaximumProcWorkers:        8,
		MaximumReplyWorkers:       8,
		MaximumUnlinkReplyWorkers: 8,
	}
}

func (g *genericHIDEchoDevice) GetBusID() usbprotocol.BusID {
	return g.deviceInfo.BusID
}

func (g *genericHIDEchoDevice) SetBusID(busNum, devNum uint) {
	busIDString := fmt.Sprintf("%d-%d", busNum, devNum)
	var busID usbprotocol.BusID
	var path [256]byte
	copy(busID[:], []byte(busIDString))
	copy(path[:], []byte("/sys/devices/pci0000:00/0000:00:1d.1/usb3/"+busIDString))
	g.deviceInfo.BusID = busID
	g.deviceInfo.BusNum = uint32(busNum)
	g.deviceInfo.DevNum = uint32(devNum)
	g.deviceInfo.Path = path
}

func (g *genericHIDEchoDevice) GetDeviceInfo() op.DeviceInfo {
	return g.deviceInfo
}

func (g *genericHIDEchoDevice) Process(data command.CmdSubmit) command.RetSubmit {
	switch data.EndpointNumber {
	case usbprotocol.ENDPOINT_CONTROL:
		{
			var setupPacket usbprotocol.SetupPacket
			retTransferBuffer := make([]byte, data.TransferBufferLength)
			if err := setupPacket.Decode(bytes.NewBuffer(data.Setup[:])); err != nil {
				g.logger.Error("unable to decode SetupPacket", "err", err)
				return g.createErrorRetSubmit(data.CmdHeader)
			}

			retData, err := g.processControlMsg(setupPacket)
			if err != nil {
				g.logger.Error("unable to process control message", "err", err)
				return g.createErrorRetSubmit(data.CmdHeader)
			}

			copy(retTransferBuffer, retData)
			return g.createSuccessRetSubmit(data.CmdHeader, retTransferBuffer)
		}
	case usbprotocol.ENDPOINT_DEV_TO_HOST:
		{
			retData, err := g.releaseEchoString()
			if err != nil {
				g.logger.Error("unable to process data message", "err", err)
				return g.createErrorRetSubmit(data.CmdHeader)
			}
			if len(retData) >= int(data.TransferBufferLength) {
				retData = retData[:data.TransferBufferLength]
			}
			return g.createSuccessRetSubmit(data.CmdHeader, retData)
		}
	case usbprotocol.ENDPOINT_HOST_TO_DEV:
		{
			err := g.queueEchoString(data)
			if err != nil {
				g.logger.Error("unable to process data message", "err", err)
				return g.createErrorRetSubmit(data.CmdHeader)
			}
			return g.createSuccessRetSubmit(data.CmdHeader, nil)
		}
	default:
		g.logger.Error("unknown endpoint number", "endpoint", data.EndpointNumber)

		return g.createErrorRetSubmit(data.CmdHeader)
	}
}

func (g *genericHIDEchoDevice) createErrorRetSubmit(header command.CmdHeader) command.RetSubmit {
	return command.RetSubmit{
		CmdHeader: command.CmdHeader{
			Command: command.RET_SUBMIT,
			SeqNum:  header.SeqNum,
		},
		Status: 99,
	}
}

func (g *genericHIDEchoDevice) createSuccessRetSubmit(header command.CmdHeader, returnData []byte) command.RetSubmit {
	return command.RetSubmit{
		CmdHeader: command.CmdHeader{
			Command: command.RET_SUBMIT,
			SeqNum:  header.SeqNum,
		},
		Status:          0,
		ActualLength:    uint32(len(returnData)),
		StartFrame:      0,
		NumberOfPackets: 0,
		ErrorCount:      0,
		Padding:         0,
		TransferBuffer:  returnData,
	}
}

func (g *genericHIDEchoDevice) processControlMsg(setup usbprotocol.SetupPacket) ([]byte, error) {
	g.logger.Debug("Received control message SetupPacket", "setup", setup)
	switch setup.BMRequestType.Recipient() {
	case usbprotocol.SETUP_RECIPIENT_DEVICE:
		return g.processControlDeviceMsg(setup)
	case usbprotocol.SETUP_RECIPIENT_INTERFACE:
		return g.processControlInterfaceMsg(setup)
	default:
		return nil, fmt.Errorf("unknown or unimplemented SetupPacket's BMRequestType: %x", setup.BMRequestType)
	}
}

// processControlDeviceMsg processes requests that are sent to Device, which should be all standard requests
func (g *genericHIDEchoDevice) processControlDeviceMsg(setup usbprotocol.SetupPacket) ([]byte, error) {
	g.logger.Debug("Processing control device msg")
	switch setup.BRequest {
	case usbprotocol.REQUEST_GET_DESCRIPTOR:
		descriptorType, index := descriptor.GetDescriptorTypeAndIndex(setup.WValue)
		return g.getDescriptor(descriptorType, index)
	case usbprotocol.REQUEST_GET_STATUS:
		// It's a self-powered device, in little-endian format
		return []byte{0x01, 0x00}, nil
	case usbprotocol.REQUEST_SET_CONFIGURATION:
		// Configuration only have one and cannot be changed, so no-op
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown or unimplemented SetupPacket's BRequest: %d", setup.BRequest)
	}
}

// processControlInterfaceMsg processes requests that are sent to Interface, which should be HID requests
func (g *genericHIDEchoDevice) processControlInterfaceMsg(setup usbprotocol.SetupPacket) ([]byte, error) {
	g.logger.Debug("Processing control interface msg")
	switch setup.BRequest {
	case usbprotocol.REQUEST_GET_DESCRIPTOR:
		descriptorType, index := descriptor.GetDescriptorTypeAndIndex(setup.WValue)
		return g.getDescriptor(descriptorType, index)
	case usbprotocol.REQUEST_HID_SET_IDLE:
		// this is software device and don't have idle mechanism, so no-op
		return nil, nil
	case usbprotocol.REQUEST_HID_SET_PROTOCOL:
		// we always use vendor-specific protocol, so no-op
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown BRequest type: %d", setup.BRequest)
	}
}

func (g *genericHIDEchoDevice) getDescriptor(descriptorType descriptor.DescriptorType, index uint8) ([]byte, error) {
	switch descriptorType {
	case descriptor.DESCRIPTOR_TYPE_DEVICE:
		desc := g.getDeviceDescriptor()
		buf := new(bytes.Buffer)
		if err := desc.Encode(buf); err != nil {
			return nil, fmt.Errorf("unable to encode standard device descriptor: %w", err)
		}
		return buf.Bytes(), nil
	case descriptor.DESCRIPTOR_TYPE_CONFIGURATION:
		hidDesc := g.getHIDDescriptor(uint16(len(echoHIDReport)))
		intfDesc := g.getInterfaceDescriptor()
		endpointDescs := g.getEndpointDescriptor()
		configDescBuf := new(bytes.Buffer)
		configDescDetailBuf := new(bytes.Buffer)

		if err := intfDesc.Encode(configDescDetailBuf); err != nil {
			return nil, fmt.Errorf("unable to encode interface descriptor: %w", err)
		}
		if err := hidDesc.Encode(configDescDetailBuf); err != nil {
			return nil, fmt.Errorf("unable to encode HID descriptor: %w", err)
		}
		for _, endpointDesc := range endpointDescs {
			if err := endpointDesc.Encode(configDescDetailBuf); err != nil {
				return nil, fmt.Errorf("unable to encode HID descriptor: %w", err)
			}
		}

		configDesc := g.getConfigurationDescriptor(uint16(configDescDetailBuf.Len()))
		if err := configDesc.Encode(configDescBuf); err != nil {
			return nil, fmt.Errorf("unable to encode configuration descriptor: %w", err)
		}

		// They are appended because USB/IP client requests 2 times as follow
		// 1. Request for Configuration Descriptor only (9 bytes), and this result will be sliced
		// 	  based on CmdSubmit's TransferBufferLength
		// 2. The client side learn the actual length of data, it request again.
		//    This time, it request the whole data (34 bytes), so no data is sliced out.
		return append(configDescBuf.Bytes(), configDescDetailBuf.Bytes()...), nil
	case descriptor.DESCRIPTOR_TYPE_STRING:
		stringDescBuf := new(bytes.Buffer)
		stringContent := g.getStringDescriptor(index)
		if err := stringContent.Encode(stringDescBuf); err != nil {
			return nil, fmt.Errorf("unable to encode string descriptor: %w", err)
		}
		return stringDescBuf.Bytes(), nil
	case descriptor.DESCRIPTOR_TYPE_HID_REPORT:
		return echoHIDReport, nil
	default:
		return nil, fmt.Errorf("unknown or unimplemented descriptor type for getting descriptor: type: %d, index: %d", descriptorType, index)
	}
}

func (g *genericHIDEchoDevice) getDeviceDescriptor() descriptor.StandardDeviceDescriptor {
	return descriptor.StandardDeviceDescriptor{
		BLength:            descriptor.STANDARD_DEVICE_DESCRIPTOR_LENGTH,
		BDescriptorType:    descriptor.DESCRIPTOR_TYPE_DEVICE,
		BCDUSB:             usbprotocol.HID_SPEC_VERSION,
		BDeviceClass:       usbprotocol.CLASS_BASEDON_INTERFACE,
		BDeviceSubClass:    usbprotocol.SUBCLASS_NONE,
		BDeviceProtocol:    usbprotocol.PROTOCOL_NONE,
		BMaxPacketSize:     64,
		IDVendor:           g.deviceInfo.IDVendor,
		IDProduct:          g.deviceInfo.IDProduct,
		BCDDevice:          g.deviceInfo.BCDDevice,
		IManufacturer:      1, // String descriptor
		IProduct:           2, // String descriptor
		ISerialNumber:      3, // String descriptor
		BNumConfigurations: 1,
	}
}

func (g *genericHIDEchoDevice) getConfigurationDescriptor(totalDetailLength uint16) descriptor.StandardConfigurationDescriptor {
	return descriptor.StandardConfigurationDescriptor{
		BLength:             descriptor.STANDARD_CONFIGURATION_DESCRIPTOR_LENGTH,
		BDescriptorType:     descriptor.DESCRIPTOR_TYPE_CONFIGURATION,
		WTotalLength:        descriptor.STANDARD_CONFIGURATION_DESCRIPTOR_LENGTH + totalDetailLength,
		BNumInterfaces:      1,
		BMAttributes:        0b01000000,
		BMaxPower:           0x32, // 100mA
		BConfigurationValue: 1,
		IConfiguration:      4, // String descriptor
	}
}

func (g *genericHIDEchoDevice) getInterfaceDescriptor() descriptor.StandardInterfaceDescriptor {
	return descriptor.StandardInterfaceDescriptor{
		BLength:            descriptor.STANDARD_INTERFACE_DESCRIPTOR_LENGTH,
		BDescriptorType:    descriptor.DESCRIPTOR_TYPE_INTERFACE,
		BInterfaceNumber:   0,
		BAlternateSetting:  0,
		BNumEndpoints:      2,
		BInterfaceClass:    usbprotocol.CLASS_HID,
		BInterfaceSubClass: usbprotocol.SUBCLASS_NONE,
		BInterfaceProtocol: usbprotocol.PROTOCOL_NONE,
		IInterface:         5, // String descriptor
	}
}

func (g *genericHIDEchoDevice) getEndpointDescriptor() []descriptor.StandardEndpointDescriptor {
	return []descriptor.StandardEndpointDescriptor{
		{
			BLength:          descriptor.STANDARD_ENDPOINT_DESCRIPTOR_LENGTH,
			BDescriptorType:  descriptor.DESCRIPTOR_TYPE_ENDPOINT,
			BEndpointAddress: 0b10000001, // Endpoint IN #1
			BMAttributes:     0b00000011, // Interrupt
			WMaxPacketSize:   64,         // 64 bytes
			BInterval:        128,        // 128ms
		},
		{
			BLength:          descriptor.STANDARD_ENDPOINT_DESCRIPTOR_LENGTH,
			BDescriptorType:  descriptor.DESCRIPTOR_TYPE_ENDPOINT,
			BEndpointAddress: 0b00000010, // Endpoint OUT #2
			BMAttributes:     0b00000011, // Interrupt
			WMaxPacketSize:   64,         // 64 bytes
			BInterval:        128,        // 128ms
		},
	}
}

func (g *genericHIDEchoDevice) getHIDDescriptor(hidReportLength uint16) hid.HIDDescriptor {
	return hid.HIDDescriptor{
		BLength:              hid.HID_DESCRIPTOR_LENGTH,
		BDescriptorType:      descriptor.DESCRIPTOR_TYPE_HID,
		BCDHID:               usbprotocol.HID_CLASS_SPEC_VERSION,
		BCountryCode:         0,
		BNumDescriptors:      1,
		BClassDescriptorType: descriptor.DESCRIPTOR_TYPE_HID_REPORT,
		WDescriptorLength:    hidReportLength,
	}
}

func (g *genericHIDEchoDevice) getStringDescriptor(index uint8) descriptor.StringDescriptor {
	var content []uint16
	switch index {
	case 0: // For zero index, it return list of supported LangIDs
		content = []uint16{uint16(descriptor.LANGID_ENGLISH_UNITED_STATES)}
	case 1: // Manufacturer
		content = utf16.Encode([]rune("ntch.dev"))
	case 2: // Product
		content = utf16.Encode([]rune("String echo device"))
	case 3: // Serial Number
		content = utf16.Encode([]rune("NTCHDEV0002"))
	case 4: // First Configuration
		content = utf16.Encode([]rune("Default Configuration"))
	case 5: // First interface
		content = utf16.Encode([]rune("Default Interface"))
	default:
		g.logger.Error("Unknown string descriptor index", "index", index)
	}

	return descriptor.StringDescriptor{
		BLength:         uint8(2 + len(content)*2),
		BDescriptorType: descriptor.DESCRIPTOR_TYPE_STRING,
		Content:         content,
	}
}

func (g *genericHIDEchoDevice) releaseEchoString() ([]byte, error) {
	select {
	case content, ok := <-g.echoContent:
		{
			if ok {
				return []byte(content), nil
			} else {
				return nil, fmt.Errorf("echo content queue is closed")
			}
		}
	default:
		return nil, nil
	}
}

func (g *genericHIDEchoDevice) queueEchoString(cmd command.CmdSubmit) error {
	str := strings.Trim(string(cmd.TransferBuffer), "\x00")
	g.logger.Debug("Got String to be echoed", "str", str)
	g.echoContent <- str

	return nil
}

func (g *genericHIDEchoDevice) Close() error {
	close(g.echoContent)

	return nil
}
