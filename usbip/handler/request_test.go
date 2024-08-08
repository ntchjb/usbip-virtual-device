package handler_test

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"testing"

	"github.com/ntchjb/usbip-virtual-device/usb"
	usbprotocol "github.com/ntchjb/usbip-virtual-device/usb/protocol"
	"github.com/ntchjb/usbip-virtual-device/usbip/handler"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var (
	deviceInfoPath = [256]byte{
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
		0xBB, 0xBB, 0xBB, 0xBB, 0xBB, 0xBB,
	}

	opDevListReq = protocol.OpReqDevList{
		OpHeader: protocol.OpHeader{
			Version:            protocol.VERSION,
			CommandOrReplyCode: protocol.OP_REQ_DEVLIST,
			Status:             protocol.OP_STATUS_OK,
		},
	}
	opDevListRep = protocol.OpRepDevList{
		OpHeader: protocol.OpHeader{
			Version:            protocol.VERSION,
			CommandOrReplyCode: protocol.OP_REP_DEVLIST,
			Status:             0,
		},
		DeviceCount: 2,
		Devices: []protocol.DeviceInfo{
			{
				DeviceInfoTruncated: protocol.DeviceInfoTruncated{
					Path: deviceInfoPath,
					BusID: [32]byte{
						0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
						0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
						0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
						0xBB, 0xBB,
					},
					BusNum:              1,
					DevNum:              5,
					Speed:               usbprotocol.SPEED_USB2_HIGH,
					IDVendor:            0xABCD,
					IDProduct:           0xDCBA,
					BCDDevice:           127,
					BDeviceClass:        usbprotocol.CLASS_HID,
					BDeviceSubclass:     usbprotocol.SUBCLASS_HID_BOOT_INTERFACE,
					BDeviceProtocol:     usbprotocol.PROTOCOL_HID_KEYBOARD,
					BConfigurationValue: 1,
					BNumConfigurations:  3,
					BNumInterfaces:      3,
				},
				Interfaces: []protocol.DeviceInterface{
					{
						BInterfaceClass:    usbprotocol.CLASS_HID,
						BInterfaceSubclass: usbprotocol.SUBCLASS_HID_BOOT_INTERFACE,
						BInterfaceProtocol: usbprotocol.PROTOCOL_HID_MOUSE,
						PaddingAlignment:   0,
					},
					{
						BInterfaceClass:    usbprotocol.CLASS_AUDIO,
						BInterfaceSubclass: 0xAB,
						BInterfaceProtocol: 0xFF,
						PaddingAlignment:   0,
					},
					{
						BInterfaceClass:    usbprotocol.CLASS_AUDIO_AND_VIDEO,
						BInterfaceSubclass: 0xAA,
						BInterfaceProtocol: 0xFE,
						PaddingAlignment:   0,
					},
				},
			},
			{
				DeviceInfoTruncated: protocol.DeviceInfoTruncated{
					Path: deviceInfoPath,
					BusID: [32]byte{
						0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
						0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
						0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xAA,
						0xBB, 0xCC,
					},
					BusNum:              1,
					DevNum:              5,
					Speed:               usbprotocol.SPEED_USB2_HIGH,
					IDVendor:            0xABCD,
					IDProduct:           0xDCBA,
					BCDDevice:           127,
					BDeviceClass:        usbprotocol.CLASS_HID,
					BDeviceSubclass:     usbprotocol.SUBCLASS_HID_BOOT_INTERFACE,
					BDeviceProtocol:     usbprotocol.PROTOCOL_HID_KEYBOARD,
					BConfigurationValue: 1,
					BNumConfigurations:  3,
					BNumInterfaces:      3,
				},
				Interfaces: []protocol.DeviceInterface{
					{
						BInterfaceClass:    usbprotocol.CLASS_HID,
						BInterfaceSubclass: usbprotocol.SUBCLASS_HID_BOOT_INTERFACE,
						BInterfaceProtocol: usbprotocol.PROTOCOL_HID_MOUSE,
						PaddingAlignment:   0,
					},
					{
						BInterfaceClass:    usbprotocol.CLASS_AUDIO,
						BInterfaceSubclass: 0xA1,
						BInterfaceProtocol: 0xF1,
						PaddingAlignment:   0,
					},
					{
						BInterfaceClass:    usbprotocol.CLASS_AUDIO_AND_VIDEO,
						BInterfaceSubclass: 0xA2,
						BInterfaceProtocol: 0xF2,
						PaddingAlignment:   0,
					},
				},
			},
		},
	}
	opDevImport = protocol.OpReqImport{
		OpHeader: protocol.OpHeader{
			Version:            protocol.VERSION,
			CommandOrReplyCode: protocol.OP_REQ_IMPORT,
			Status:             0,
		},
		BusID: opDevListRep.Devices[1].BusID,
	}
	opDevImportRep = protocol.OpRepImport{
		OpHeader: protocol.OpHeader{
			Version:            protocol.VERSION,
			CommandOrReplyCode: protocol.OP_REP_IMPORT,
			Status:             0,
		},
		DeviceInfo: opDevListRep.Devices[1].DeviceInfoTruncated,
	}

	cmdSubmit = protocol.CmdSubmit{
		CmdHeader: protocol.CmdHeader{
			Command:        protocol.CMD_SUBMIT,
			SeqNum:         1,
			DevID:          0x00010001,
			Direction:      protocol.DIR_OUT,
			EndpointNumber: 1,
		},
		TransferFlags:        0,
		TransferBufferLength: 5,
		StartFrame:           0,
		NumberOfPackets:      0xffffffff,
		Interval:             100,
		Setup:                [8]byte{},
		TransferBuffer:       []byte{0x01, 0x02, 0x03, 0x04, 0x05},
		ISOPacketDescriptors: nil,
	}

	retSubmit = protocol.RetSubmit{
		CmdHeader: protocol.CmdHeader{
			Command:        protocol.RET_SUBMIT,
			SeqNum:         cmdSubmit.SeqNum,
			DevID:          cmdSubmit.DevID,
			Direction:      cmdSubmit.Direction,
			EndpointNumber: cmdSubmit.EndpointNumber,
		},
		Status:               0,
		ActualLength:         0,
		StartFrame:           0,
		NumberOfPackets:      0,
		ErrorCount:           0,
		Padding:              0,
		TransferBuffer:       nil,
		ISOPacketDescriptors: nil,
	}

	cmdUnlink = protocol.CmdUnlink{
		CmdHeader: protocol.CmdHeader{
			Command:        protocol.CMD_UNLINK,
			SeqNum:         2,
			DevID:          0x00010001,
			Direction:      protocol.DIR_OUT,
			EndpointNumber: 1,
		},
		UnlinkSeqNum: 1,
	}
)

func TestRequestHandler(t *testing.T) {
	ctrl := gomock.NewController(t)

	registrar := usb.NewMockDeviceRegistrar(ctrl)
	device1 := usb.NewMockDevice(ctrl)
	device2 := usb.NewMockDevice(ctrl)
	worker := handler.NewMockWorkerPool(ctrl)
	logger := slog.Default()
	server, client := net.Pipe()
	reqHandler := handler.NewRequestHandler(server, registrar, worker, logger)

	// DevList
	registrar.EXPECT().GetAvailableDevices().Return([]usb.Device{
		device1, device2,
	})
	device1.EXPECT().GetDeviceInfo().Return(opDevListRep.Devices[0]).AnyTimes()
	device2.EXPECT().GetDeviceInfo().Return(opDevListRep.Devices[1]).AnyTimes()

	// DevImport
	registrar.EXPECT().GetDevice(opDevListRep.Devices[0].BusID).Return(device1, nil).AnyTimes()
	registrar.EXPECT().GetDevice(opDevListRep.Devices[1].BusID).Return(device2, nil)
	registrar.EXPECT().GetDevice(gomock.Any).Return(nil, fmt.Errorf("unknown device")).AnyTimes()
	worker.EXPECT().SetDevice(device1).Return()
	worker.EXPECT().Start().Return(nil)

	// CmdSubmit & CmdUnlink
	worker.EXPECT().PublishCmdSubmit(cmdSubmit).Return()
	worker.EXPECT().Unlink(cmdUnlink).Return(nil)

	var wg sync.WaitGroup

	// Client code
	wg.Add(1)
	go func() {
		defer wg.Done()

		err := opDevListReq.Encode(client)
		assert.NoError(t, err)

		actualOpDevListRep := protocol.OpRepDevList{}
		err = actualOpDevListRep.OpHeader.Decode(client)
		assert.NoError(t, err)
		err = actualOpDevListRep.Decode(client)
		assert.NoError(t, err)
		assert.Equal(t, opDevListRep, actualOpDevListRep)

		err = opDevImport.OpHeader.Encode(client)
		assert.NoError(t, err)
		err = opDevImport.Encode(client)
		assert.NoError(t, err)

		actualOpImportRep := protocol.OpRepImport{}
		err = actualOpImportRep.OpHeader.Decode(client)
		assert.NoError(t, err)
		err = actualOpImportRep.Decode(client)
		assert.NoError(t, err)
		assert.Equal(t, actualOpImportRep, opDevImportRep)

		err = cmdSubmit.CmdHeader.Encode(client)
		assert.NoError(t, err)
		err = cmdSubmit.Encode(client)
		assert.NoError(t, err)

		err = cmdUnlink.CmdHeader.Encode(client)
		assert.NoError(t, err)
		err = cmdUnlink.Encode(client)
		assert.NoError(t, err)
	}()

	// Server code
	assert.Equal(t, handler.HANDLER_LEVEL_OP, reqHandler.GetHandlerLevel())

	header, err := reqHandler.HandleOpHeader()
	assert.NoError(t, err)
	assert.Equal(t, opDevListReq.OpHeader, header)

	err = reqHandler.HandleOpDevList(header)
	assert.NoError(t, err)

	header, err = reqHandler.HandleOpHeader()
	assert.NoError(t, err)
	assert.Equal(t, opDevImport.OpHeader, header)
	err = reqHandler.HandleOpImport(header)
	assert.NoError(t, err)
	assert.Equal(t, handler.HANDLER_LEVEL_CMD, reqHandler.GetHandlerLevel())

	cmdHeader, err := reqHandler.HandleCmdHeader()
	assert.NoError(t, err)
	err = reqHandler.HandleCmdSubmit(cmdHeader)
	assert.NoError(t, err)

	cmdHeader, err = reqHandler.HandleCmdHeader()
	assert.NoError(t, err)
	err = reqHandler.HandleCmdUnlink(cmdHeader)
	assert.NoError(t, err)

	wg.Wait()
}
