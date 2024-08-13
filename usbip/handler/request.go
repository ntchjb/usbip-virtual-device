package handler

import (
	"encoding/hex"
	"fmt"
	"io"
	"log/slog"
	"net"

	"github.com/ntchjb/usbip-virtual-device/usb"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol/command"
	operation "github.com/ntchjb/usbip-virtual-device/usbip/protocol/op"
)

type HandlerLevel uint8

const (
	HANDLER_LEVEL_OP  HandlerLevel = 0
	HANDLER_LEVEL_CMD HandlerLevel = 1
)

type RequestHandler interface {
	HandleOpHeader() (operation.OpHeader, error)
	HandleCmdHeader() (command.CmdHeader, error)
	HandleOpDevList(opHeader operation.OpHeader) error
	HandleOpImport(opHeader operation.OpHeader) error
	HandleCmdSubmit(cmdHeader command.CmdHeader) error
	HandleCmdUnlink(cmdHeader command.CmdHeader) error
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

func (op *requestHandlerImpl) HandleOpHeader() (operation.OpHeader, error) {
	var opHeader operation.OpHeader
	// When idle, logic will stuck here, waiting for new data
	if err := opHeader.Decode(op.conn); err != nil {
		return opHeader, fmt.Errorf("unable to decode OpHeader: %w", err)
	}
	if opHeader.Version != operation.VERSION {
		return opHeader, fmt.Errorf("unsupported USBIP protocol version, expected: %x, actual: %x", operation.VERSION, opHeader)
	}

	return opHeader, nil
}

func (op *requestHandlerImpl) HandleCmdHeader() (command.CmdHeader, error) {
	var cmdHeader command.CmdHeader
	// When idle, logic will stuck here, waiting for new data
	if err := cmdHeader.Decode(op.conn); err != nil {
		return cmdHeader, fmt.Errorf("unable to decode CmdHeader: %w", err)
	}

	return cmdHeader, nil
}

func (op *requestHandlerImpl) HandleOpDevList(opHeader operation.OpHeader) error {
	// TODO:
	// 1. Get device list from registrar
	// 2. Reply a list of USB devices
	devices := op.registrar.GetAvailableDevices()
	var reply operation.OpRepDevList
	var replyHeader operation.OpHeader

	replyHeader.Version = opHeader.Version
	replyHeader.CommandOrReplyCode = operation.OP_REP_DEVLIST
	replyHeader.Status = 0

	reply.DeviceCount = uint32(len(devices))
	if len(devices) > 0 {
		reply.Devices = make([]operation.DeviceInfo, len(devices))
	}
	for i, device := range devices {
		reply.Devices[i] = device.GetDeviceInfo()
	}

	op.logger.Debug("OP_DEVLIST_REPLY", "reply", reply, "replyHeader", replyHeader)
	if err := replyHeader.Encode(op.conn); err != nil {
		return fmt.Errorf("unable to encode OpHeader for OpRepDevList: %w", err)
	}
	if err := reply.Encode(op.conn); err != nil {
		return fmt.Errorf("unable to encode OpRepDevList: %w", err)
	}

	return nil
}

func (op *requestHandlerImpl) HandleOpImport(opHeader operation.OpHeader) error {
	opReqImport := operation.OpReqImport{
		OpHeader: opHeader,
	}
	var reply operation.OpRepImport
	var replyHeader operation.OpHeader

	if err := opReqImport.Decode(op.conn); err != nil {
		return fmt.Errorf("unable to decode OpReqImport: %w", err)
	}

	replyHeader.Version = opHeader.Version
	replyHeader.CommandOrReplyCode = operation.OP_REP_IMPORT
	device, err := op.registrar.GetDevice(opReqImport.BusID)
	if err != nil {
		op.logger.Error("unable to get USB device from registrar", "err", err)
		replyHeader.Status = operation.OP_STATUS_ERROR
	} else {
		replyHeader.Status = operation.OP_STATUS_OK
		reply.DeviceInfo = device.GetDeviceInfo().DeviceInfoTruncated
	}

	op.logger.Debug("OP_IMPORT_REPLY", "reply", reply, "replyHeader", replyHeader)
	if err := replyHeader.Encode(op.conn); err != nil {
		return fmt.Errorf("unable to encode OpHeader: for OpRepImport: %w", err)
	}
	if replyHeader.Status != operation.OP_STATUS_OK {
		return io.EOF
	}

	if err := reply.Encode(op.conn); err != nil {
		return fmt.Errorf("unable to encode OpRepImport: %w", err)
	}
	op.worker.SetDevice(device)
	op.worker.Start()

	op.level = HANDLER_LEVEL_CMD

	op.logger.Info("Device attached", "busID", hex.EncodeToString(opReqImport.BusID[:]), "id", fmt.Sprintf("%04x:%04x", reply.DeviceInfo.IDVendor, reply.DeviceInfo.IDProduct))
	return nil
}

func (op *requestHandlerImpl) HandleCmdSubmit(cmdHeader command.CmdHeader) error {
	// TODO: Handle CmdSubmit in asynchronous way
	cmdSubmit := command.CmdSubmit{
		CmdHeader: cmdHeader,
	}

	if err := cmdSubmit.Decode(op.conn); err != nil {
		return fmt.Errorf("unable to decode CmdSubmit: %w", err)
	}
	op.worker.PublishCmdSubmit(cmdSubmit)

	return nil
}

func (op *requestHandlerImpl) HandleCmdUnlink(cmdHeader command.CmdHeader) error {
	cmdUnlink := command.CmdUnlink{
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
