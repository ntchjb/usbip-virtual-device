package main

import (
	"bytes"
	"fmt"
	"log/slog"
	"time"
	"unicode/utf16"

	"github.com/ntchjb/usbip-virtual-device/usb"
	usbprotocol "github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
)

var (
	mouseHIDReport = []byte{
		0x05, 0x01, // Usage Page (Generic Desktop),
		0x09, 0x02, // Usage (Mouse),
		0xA1, 0x01, // 	Collection (Application),
		0x09, 0x01, // 		Usage (Pointer),
		0xA1, 0x00, // 		Collection (Physical),
		0x05, 0x09, // 			Usage Page (Buttons),
		0x19, 0x01, // 			Usage Minimum (01),
		0x29, 0x03, // 			Usage Maximun (03),
		0x15, 0x00, // 			Logical Minimum (0),
		0x25, 0x01, //			Logical Maximum (1),
		0x95, 0x03, //			Report Count (3),
		0x75, 0x01, //			Report Size (1),
		0x81, 0x02, //			Input (Data, Variable, Absolute), ;3 button bits
		0x95, 0x01, //			Report Count (1),
		0x75, 0x05, //			Report Size (5),
		0x81, 0x01, //			Input (Constant), ;5 bit padding
		0x05, 0x01, //			Usage Page (Generic Desktop),
		0x09, 0x30, //			Usage (X),
		0x09, 0x31, //			Usage (Y),
		0x15, 0x81, //			Logical Minimum (-127),
		0x25, 0x7F, //			Logical Maximum (127),
		0x75, 0x08, //			Report Size (8),
		0x95, 0x02, //			Report Count (2),
		0x81, 0x06, //			Input (Data, Variable, Relative), ;2 position bytes (X & Y)
		0xC0, //		End Collection,
		0xC0, //	End Collection
	}
)

type genericHIDMouseDevice struct {
	deviceInfo protocol.DeviceInfo
	logger     *slog.Logger

	state int
}

func NewGenericHIDMouseDevice(logger *slog.Logger) usb.Device {
	return &genericHIDMouseDevice{
		logger: logger,
		deviceInfo: protocol.DeviceInfo{
			DeviceInfoTruncated: protocol.DeviceInfoTruncated{
				Speed:               usbprotocol.SPEED_USB2_HIGH,
				IDVendor:            0x0ff0, // a random vendor ID
				IDProduct:           0x0123, // a random product ID
				BCDDevice:           1,
				BDeviceClass:        usbprotocol.CLASS_BASEDON_INTERFACE,
				BDeviceSubclass:     usbprotocol.SUBCLASS_NONE,
				BDeviceProtocol:     usbprotocol.PROTOCOL_NONE,
				BConfigurationValue: 1,
				BNumConfigurations:  1,
				BNumInterfaces:      1,
			},
			Interfaces: []protocol.DeviceInterface{
				{
					BInterfaceClass:    usbprotocol.CLASS_HID,
					BInterfaceSubclass: usbprotocol.SUBCLASS_HID_BOOT_INTERFACE,
					BInterfaceProtocol: usbprotocol.PROTOCOL_HID_MOUSE,
				},
			},
		},
	}
}

func (g *genericHIDMouseDevice) GetWorkerPoolProfile() usb.WorkerPoolProfile {
	return usb.WorkerPoolProfile{
		MaximumProcWorkers:        1,
		MaximumReplyWorkers:       1,
		MaximumUnlinkReplyWorkers: 1,
	}
}

func (g *genericHIDMouseDevice) GetBusID() usbprotocol.BusID {
	return g.deviceInfo.BusID
}

func (g *genericHIDMouseDevice) SetBusID(busNum, devNum uint) {
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

func (g *genericHIDMouseDevice) GetDeviceInfo() protocol.DeviceInfo {
	return g.deviceInfo
}

func (g *genericHIDMouseDevice) Process(data protocol.CmdSubmit) protocol.RetSubmit {
	retTransferBuffer := make([]byte, data.TransferBufferLength)
	switch data.EndpointNumber {
	case usbprotocol.ENDPOINT_CONTROL:
		{
			var setupPacket usbprotocol.SetupPacket
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
			retData, err := g.proceeHIDData(data)
			if err != nil {
				g.logger.Error("unable to process data message", "err", err)
				return g.createErrorRetSubmit(data.CmdHeader)
			}
			copy(retTransferBuffer, retData)
			return g.createSuccessRetSubmit(data.CmdHeader, retTransferBuffer)
		}
	default:
		g.logger.Error("unknown endpoint number", "endpoint", data.EndpointNumber)

		return g.createErrorRetSubmit(data.CmdHeader)
	}
}

func (g *genericHIDMouseDevice) createErrorRetSubmit(header protocol.CmdHeader) protocol.RetSubmit {
	return protocol.RetSubmit{
		CmdHeader: protocol.CmdHeader{
			Command: protocol.RET_SUBMIT,
			SeqNum:  header.SeqNum,
		},
		Status: 99,
	}
}

func (g *genericHIDMouseDevice) createSuccessRetSubmit(header protocol.CmdHeader, returnData []byte) protocol.RetSubmit {
	return protocol.RetSubmit{
		CmdHeader: protocol.CmdHeader{
			Command: protocol.RET_SUBMIT,
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

func (g *genericHIDMouseDevice) getDeviceDescriptor() usbprotocol.StandardDeviceDescriptor {
	return usbprotocol.StandardDeviceDescriptor{
		BLength:            usbprotocol.STANDARD_DEVICE_DESCRIPTOR_LENGTH,
		BDescriptorType:    usbprotocol.DESCRIPTOR_TYPE_DEVICE,
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

func (g *genericHIDMouseDevice) getConfigurationDescriptor(totalDetailLength uint16) usbprotocol.StandardConfigurationDescriptor {
	return usbprotocol.StandardConfigurationDescriptor{
		BLength:             usbprotocol.STANDARD_CONFIGURATION_DESCRIPTOR_LENGTH,
		BDescriptorType:     usbprotocol.DESCRIPTOR_TYPE_CONFIGURATION,
		WTotalLength:        usbprotocol.STANDARD_CONFIGURATION_DESCRIPTOR_LENGTH + totalDetailLength,
		BNumInterfaces:      1,
		BMAttributes:        0b01000000,
		BMaxPower:           0x32,
		BConfigurationValue: 1,
		IConfiguration:      4, // String descriptor
	}
}

func (g *genericHIDMouseDevice) getInterfaceDescriptor() usbprotocol.StandardInterfaceDescriptor {
	return usbprotocol.StandardInterfaceDescriptor{
		BLength:            usbprotocol.STANDARD_INTERFACE_DESCRIPTOR_LENGTH,
		BDescriptorType:    usbprotocol.DESCRIPTOR_TYPE_INTERFACE,
		BInterfaceNumber:   0,
		BAlternateSetting:  0,
		BNumEndpoints:      1,
		BInterfaceClass:    usbprotocol.CLASS_HID,
		BInterfaceSubClass: usbprotocol.SUBCLASS_HID_BOOT_INTERFACE,
		BInterfaceProtocol: usbprotocol.PROTOCOL_HID_MOUSE,
		IInterface:         5, // String descriptor
	}
}

func (g *genericHIDMouseDevice) getEndpointDescriptor() []usbprotocol.StandardEndpointDescriptor {
	return []usbprotocol.StandardEndpointDescriptor{
		{
			BLength:          usbprotocol.STANDARD_ENDPOINT_DESCRIPTOR_LENGTH,
			BDescriptorType:  usbprotocol.DESCRIPTOR_TYPE_ENDPOINT,
			BEndpointAddress: 0b10000001,
			BMAttributes:     0b00000011,
			WMaxPacketSize:   8,
			BInterval:        255,
		},
	}
}

func (g *genericHIDMouseDevice) getHIDDescriptor(hidReportLength uint16) usbprotocol.HIDDescriptor {
	return usbprotocol.HIDDescriptor{
		BLength:              usbprotocol.HID_DESCRIPTOR_LENGTH,
		BDescriptorType:      usbprotocol.DESCRIPTOR_TYPE_HID,
		BCDHID:               usbprotocol.HID_CLASS_SPEC_VERSION,
		BCountryCode:         0,
		BNumDescriptors:      1,
		BClassDescriptorType: usbprotocol.DESCRIPTOR_TYPE_HID_REPORT,
		WDescriptorLength:    hidReportLength,
	}
}

func (g *genericHIDMouseDevice) getStringDescriptor(index uint8) usbprotocol.StringDescriptor {
	var content []uint16
	switch index {
	case 0: // For zero index, it return list of supported LangIDs
		content = []uint16{usbprotocol.LangIDEnglishUSA}
	case 1: // Manufacturer
		content = utf16.Encode([]rune("ntch.dev"))
	case 2: // Product
		content = utf16.Encode([]rune("Virtual Mouse"))
	case 3: // Serial Number
		content = utf16.Encode([]rune("1ABBA1BABA1"))
	case 4: // First Configuration
		content = utf16.Encode([]rune("Default Configuration"))
	case 5: // First interface
		content = utf16.Encode([]rune("Default Interface"))
	default:
		g.logger.Error("Unknown string descriptor index", "index", index)
	}

	return usbprotocol.StringDescriptor{
		BLength:         uint8(2 + len(content)*2),
		BDescriptorType: usbprotocol.DESCRIPTOR_TYPE_STRING,
		Content:         content,
	}
}

func (g *genericHIDMouseDevice) getDescriptor(descriptorType usbprotocol.DescriptorType, index uint8) ([]byte, error) {
	switch descriptorType {
	case usbprotocol.DESCRIPTOR_TYPE_DEVICE:
		desc := g.getDeviceDescriptor()
		buf := new(bytes.Buffer)
		if err := desc.Encode(buf); err != nil {
			return nil, fmt.Errorf("unable to encode standard device descriptor: %w", err)
		}
		return buf.Bytes(), nil
	case usbprotocol.DESCRIPTOR_TYPE_CONFIGURATION:
		hidDesc := g.getHIDDescriptor(uint16(len(mouseHIDReport)))
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
	case usbprotocol.DESCRIPTOR_TYPE_STRING:
		stringDescBuf := new(bytes.Buffer)
		stringContent := g.getStringDescriptor(index)
		if err := stringContent.Encode(stringDescBuf); err != nil {
			return nil, fmt.Errorf("unable to encode string descriptor: %w", err)
		}
		return stringDescBuf.Bytes(), nil
	case usbprotocol.DESCRIPTOR_TYPE_HID_REPORT:
		return mouseHIDReport, nil
	default:
		return nil, fmt.Errorf("unknown or unimplemented descriptor type for getting descriptor: type: %d, index: %d", descriptorType, index)
	}
}

func (g *genericHIDMouseDevice) processControlMsg(setup usbprotocol.SetupPacket) ([]byte, error) {
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
func (g *genericHIDMouseDevice) processControlDeviceMsg(setup usbprotocol.SetupPacket) ([]byte, error) {
	g.logger.Debug("Processing control device msg")
	switch setup.BRequest {
	case usbprotocol.REQUEST_GET_DESCRIPTOR:
		descriptorType, index := usbprotocol.GetDescriptorTypeAndIndex(setup.WValue)
		return g.getDescriptor(descriptorType, index)
	case usbprotocol.REQUEST_GET_STATUS:
		// It's a self-powered device, in little-endian form
		return []byte{0x01, 0x00}, nil
	case usbprotocol.REQUEST_SET_CONFIGURATION:
		// Configuration only have one and cannot be changed, so no-op
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown or unimplemented SetupPacket's BRequest: %d", setup.BRequest)
	}
}

// processControlInterfaceMsg processes requests that are sent to Interface, which should be HID requests
func (g *genericHIDMouseDevice) processControlInterfaceMsg(setup usbprotocol.SetupPacket) ([]byte, error) {
	g.logger.Debug("Processing control interface msg")
	switch setup.BRequest {
	case usbprotocol.REQUEST_GET_DESCRIPTOR:
		descriptorType, index := usbprotocol.GetDescriptorTypeAndIndex(setup.WValue)
		return g.getDescriptor(descriptorType, index)
	case usbprotocol.REQUEST_HID_SET_IDLE:
		// this is software device, so no-op
		return nil, nil
	case usbprotocol.REQUEST_HID_SET_PROTOCOL:
		// we always use boot protocol
		return nil, nil
	default:
		return nil, fmt.Errorf("unknown BRequest type: %d", setup.BRequest)
	}
}

func (g *genericHIDMouseDevice) proceeHIDData(_ protocol.CmdSubmit) ([]byte, error) {
	// Use Sleep to replies every 100ms, which is polling rate of 10Hz
	time.Sleep(100 * time.Millisecond)

	buf := make([]byte, 3)
	var minusFive int8 = -5

	if g.state == 5 {
		buf[0] = 0 // Button 1,2,3 and device specific (1 byte)
		buf[1] = 5 // X
		buf[2] = 5 // Y
	} else if g.state == 0 {
		buf[0] = 0                // Button 1,2,3 and device specific (1 byte)
		buf[1] = uint8(minusFive) // X
		buf[2] = uint8(minusFive) // Y
	} else {
		buf[0] = 0 // Button 1,2,3 and device specific (1 byte)
		buf[1] = 0 // X
		buf[2] = 0 // Y
	}

	g.state = (g.state + 1) % 10

	return buf, nil
}
