package usbip

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ntchjb/usbip-virtual-device/usb"
	"github.com/ntchjb/usbip-virtual-device/usbip/handler"
	"github.com/ntchjb/usbip-virtual-device/usbip/protocol"
)

type USBIPServer interface {
	Open() error
	Close() error
}

type USBIPServerConfig struct {
	ListenAddress        string
	TCPConnectionTimeout time.Duration
	MaxTCPConnection     uint
}

type usbIPServerImpl struct {
	conf      USBIPServerConfig
	logger    *slog.Logger
	registrar usb.DeviceRegistrar

	listener net.Listener
	connWg   sync.WaitGroup
	quit     chan any
}

// NewUSBIPServer returns an instance of USB/IP server,
// which is a TCP server using net package.
//
// listenAddress uses format of [ip:port] such as "127.0.0.1:3456", ":3456"
func NewUSBIPServer(config USBIPServerConfig, registrar usb.DeviceRegistrar, logger *slog.Logger) USBIPServer {
	if logger == nil {
		panic(fmt.Errorf("a logger instance required for USB/IP server"))
	}
	return &usbIPServerImpl{
		conf:      config,
		logger:    logger,
		registrar: registrar,
	}
}

func (s *usbIPServerImpl) Open() error {
	var err error
	s.quit = make(chan any)

	s.listener, err = net.Listen("tcp", s.conf.ListenAddress)
	if err != nil {
		return fmt.Errorf("unable to listen to address %s: %w", s.conf.ListenAddress, err)
	}

	// TCP connection acceptor
	s.connWg.Add(1)
	go func() {
		var connCount atomic.Int64
		defer s.connWg.Done()
		for {
			conn, err := s.listener.Accept()
			// Error occurred when
			// 1. Connection error
			// 2. The listener is closed (quit channel is closed)
			if err != nil {
				select {
				case <-s.quit:
					return
				default:
					s.logger.Error("unable to accept request", "address", s.conf.ListenAddress, "err", err)
				}
			} else {
				// Check if TCP connection reached limit specified in given config
				count := connCount.Add(1)
				if count > int64(s.conf.MaxTCPConnection) {
					s.logger.Error("maximum TCP connection reached, drop the connection", "count", count)
					conn.Close()
					connCount.Add(-1)
					continue
				}

				// TCP connection handler
				s.connWg.Add(1)
				go func() {
					defer connCount.Add(-1)
					defer s.connWg.Done()
					s.logger.Info("new connection established", "addr", conn.RemoteAddr())
					s.handleConnection(conn)
				}()
			}
		}
	}()

	return nil
}

func (s *usbIPServerImpl) Close() error {
	var err error

	close(s.quit)
	listenerErr := s.listener.Close()
	if listenerErr != nil {
		err = fmt.Errorf("cannot close TCP listener: %w", listenerErr)
	}
	s.logger.Info("Closing server, waiting for all devices to be disconnected. Please make sure that USB/IP client-side devices are all unbinded and disconnected from USB/IP server")
	s.connWg.Wait()
	s.logger.Info("Server closed, bye.")

	return err
}

func (s *usbIPServerImpl) handleConnection(conn net.Conn) {
	worker := handler.NewWorkerPool(conn, s.logger)
	reqHandler := handler.NewRequestHandler(conn, s.registrar, worker, s.logger)

	defer conn.Close()
	defer worker.Stop()

	for {
		switch reqHandler.GetHandlerLevel() {
		case handler.HANDLER_LEVEL_OP:
			if err := s.handleOp(reqHandler); err != nil {
				if !errors.Is(err, io.EOF) {
					s.logger.Error("unable to handle Op request", "err", err)
				}
				return
			}
		case handler.HANDLER_LEVEL_CMD:
			if err := s.handleCmd(reqHandler); err != nil {
				if !errors.Is(err, io.EOF) {
					s.logger.Error("unable to handle Cmd request", "err", err)
				}
				return
			}
		}
	}
}

func (s *usbIPServerImpl) handleOp(reqHandler handler.RequestHandler) error {
	opHeader, err := reqHandler.HandleOpHeader()
	if err != nil {
		return fmt.Errorf("unable to handle OpHeader: %w", err)
	}

	switch opHeader.CommandOrReplyCode {
	case protocol.OP_REQ_DEVLIST:
		if err := reqHandler.HandleOpDevList(opHeader); err != nil {
			return fmt.Errorf("error occurred when handling OpDevList: %w", err)
		}
		// Close TCP connection after replied
		return io.EOF
	case protocol.OP_REQ_IMPORT:
		if err := reqHandler.HandleOpImport(opHeader); err != nil {
			return fmt.Errorf("error occurred when handling OpImport: %w", err)
		}
		// OP_REQ_IMPORT does not close TCP connection
	default:
		return fmt.Errorf("unknown operation: %x", opHeader.CommandOrReplyCode)
	}

	return nil
}

func (s *usbIPServerImpl) handleCmd(reqHandler handler.RequestHandler) error {
	cmdHeader, err := reqHandler.HandleCmdHeader()
	if err != nil {
		return fmt.Errorf("unable to handle CmdHeader: %w", err)
	}

	switch cmdHeader.Command {
	case protocol.CMD_SUBMIT:
		if err := reqHandler.HandleCmdSubmit(cmdHeader); err != nil {
			return fmt.Errorf("unable to handle CmdSubmit: %w", err)
		}
	case protocol.CMD_UNLINK:
		if err := reqHandler.HandleCmdUnlink(cmdHeader); err != nil {
			return fmt.Errorf("unable to handle CmdUnlink: %w", err)
		}
	default:
		return fmt.Errorf("unknown command: %x", cmdHeader.Command)
	}

	return nil
}
