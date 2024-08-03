package handler

import (
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net"

	"github.com/ntchjb/usbip-virtual-device/usb"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
)

type HandlerLevel uint8

const (
	HANDLER_LEVEL_OP  HandlerLevel = 0
	HANDLER_LEVEL_CMD HandlerLevel = 1
)

type RequestHandler interface {
	HandleOpHeader() (protocol.OpHeader, error)
	HandleCmdHeader() (protocol.CmdHeader, error)
	HandleOpDevList(opHeader protocol.OpHeader) error
	HandleOpImport(opHeader protocol.OpHeader) error
	HandleCmdSubmit(cmdHeader protocol.CmdHeader) error
	HandleCmdUnlink(cmdHeader protocol.CmdHeader) error
	GetHandlerLevel() HandlerLevel
}

type requestHandlerImpl struct {
	conn      net.Conn
	registrar usb.DeviceRegistrar
	logger    *slog.Logger
	worker    WorkerPool

	level HandlerLevel
}

func NewRequestHandler(conn net.Conn, registrar usb.DeviceRegistrar, worker WorkerPool, logger *slog.Logger) RequestHandler {
	return &requestHandlerImpl{
		conn:      conn,
		registrar: registrar,
		level:     HANDLER_LEVEL_OP,
		logger:    logger,
		worker:    worker,
	}
}

func (op *requestHandlerImpl) HandleOpHeader() (protocol.OpHeader, error) {
	var opHeader protocol.OpHeader
	// When idle, logic will stuck here, waiting for new data
	if err := opHeader.Decode(op.conn); err != nil {
		return opHeader, fmt.Errorf("unable to decode OpHeader: %w", err)
	}
	if opHeader.Version != protocol.VERSION {
		return opHeader, fmt.Errorf("unsupported USBIP protocol version, expected: %x, actual: %x", protocol.VERSION, opHeader)
	}

	return opHeader, nil
}

func (op *requestHandlerImpl) HandleCmdHeader() (protocol.CmdHeader, error) {
	var cmdHeader protocol.CmdHeader
	// When idle, logic will stuck here, waiting for new data
	if err := cmdHeader.Decode(op.conn); err != nil {
		return cmdHeader, fmt.Errorf("unable to decode CmdHeader: %w", err)
	}

	return cmdHeader, nil
}

func (op *requestHandlerImpl) HandleOpDevList(opHeader protocol.OpHeader) error {
	// TODO:
	// 1. Get device list from registrar
	// 2. Reply a list of USB devices
	devices := op.registrar.GetAvailableDevices()
	var reply protocol.OpRepDevList

	reply.OpHeader.Version = opHeader.Version
	reply.CommandOrReplyCode = protocol.OP_REP_DEVLIST
	reply.Status = 0
	reply.DeviceCount = uint32(len(devices))

	if len(devices) > 0 {
		reply.Devices = make([]protocol.DeviceInfo, len(devices))
	}
	for i, device := range devices {
		reply.Devices[i] = device.GetDeviceInfo()
	}

	op.logger.Debug("OP_DEVLIST_REPLY", "reply", reply)
	if err := reply.Encode(op.conn); err != nil {
		return fmt.Errorf("unable to encode OpRepDevList: %w", err)
	}

	return nil
}

func (op *requestHandlerImpl) HandleOpImport(opHeader protocol.OpHeader) error {
	opReqImport := protocol.OpReqImport{
		OpHeader: opHeader,
	}
	var reply protocol.OpRepImport

	if err := opReqImport.Decode(op.conn); err != nil {
		return fmt.Errorf("unable to decode OpReqImport: %w", err)
	}

	reply.Version = opHeader.Version
	reply.CommandOrReplyCode = protocol.OP_REP_IMPORT
	device, err := op.registrar.GetDevice(opReqImport.BusID)
	if err != nil {
		if !errors.Is(err, usb.ErrDeviceNotFound) {
			op.logger.Error("unable to get USB device from registrar", "err", err)
		}
		reply.Status = protocol.OP_STATUS_ERROR
	} else {
		reply.Status = protocol.OP_STATUS_OK
		reply.DeviceInfo = device.GetDeviceInfo().DeviceInfoTruncated
	}

	op.logger.Debug("OP_IMPORT_REPLY", "reply", reply)
	if err := reply.Encode(op.conn); err != nil {
		return fmt.Errorf("unable to encode OpRepImport: %w", err)
	}

	op.worker.SetProcessor(device.GetURBProcessor())
	op.worker.Start()

	op.level = HANDLER_LEVEL_CMD

	op.logger.Info("Device attached", "busID", hex.EncodeToString(opReqImport.BusID[:]), "id", fmt.Sprintf("%04x:%04x", reply.DeviceInfo.IDVendor, reply.DeviceInfo.IDProduct))
	return nil
}

func (op *requestHandlerImpl) HandleCmdSubmit(cmdHeader protocol.CmdHeader) error {
	// TODO: Handle CmdSubmit in asynchronous way
	cmdSubmit := protocol.CmdSubmit{
		CmdHeader: cmdHeader,
	}

	if err := cmdSubmit.Decode(op.conn); err != nil {
		return fmt.Errorf("unable to decode CmdSubmit: %w", err)
	}
	op.worker.PublishCmdSubmit(cmdSubmit)

	return nil
}

func (op *requestHandlerImpl) HandleCmdUnlink(cmdHeader protocol.CmdHeader) error {
	cmdUnlink := protocol.CmdUnlink{
		CmdHeader: cmdHeader,
	}
	if err := cmdUnlink.Decode(op.conn); err != nil {
		return fmt.Errorf("unable to decode CmdUnlink: %w", err)
	}
	if err := op.worker.Unlink(cmdUnlink); err != nil {
		return fmt.Errorf("unable to unlink with seqNum %d: %w", cmdHeader.SeqNum, err)
	}

	return nil
}

func (op *requestHandlerImpl) GetHandlerLevel() HandlerLevel {
	return op.level
}
