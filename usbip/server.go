package usbip

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"
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
	conf   USBIPServerConfig
	logger *slog.Logger

	listener net.Listener
	connWg   sync.WaitGroup
	quit     chan any
}

// NewUSBIPServer returns an instance of USB/IP server,
// which is a TCP server using net package.
//
// listenAddress uses format of [ip:port] such as "127.0.0.1:3456", ":3456"
func NewUSBIPServer(config USBIPServerConfig, logger *slog.Logger) USBIPServer {
	if logger == nil {
		panic(fmt.Errorf("a logger instance required for USB/IP server"))
	}
	return &usbIPServerImpl{
		conf:   config,
		logger: logger,
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
	s.connWg.Wait()

	return err
}

func (s *usbIPServerImpl) handleConnection(conn net.Conn) {
	defer conn.Close()

	// TODO: handle data in this TCP connection

	s.logger.Debug("received TCP connection", "addr", conn.RemoteAddr())
}
