package core

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"syscall"

	"github.com/vtools/app/vecho/logger"
)

func SetupEchoServer(protocol string, srcIP net.IP, srcPort int, opts ...SockOption) error {
	// setup options and build listen config
	newSockOptions := defaultSockOptions()
	for _, opt := range opts {
		opt.apply(newSockOptions)
	}
	listenCfg := &net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			return c.Control(newSockOptions.Control())
		},
	}

	switch protocol {
	case "tcp", "tcp4", "tcp6":
		return setupEchoServerTCP(listenCfg, protocol, srcIP, srcPort)
	case "udp", "udp4", "udp6":
		return setupEchoServerUDP(listenCfg, protocol, srcIP, srcPort)
	default:
		return errors.New("unknown protocol")
	}
}

func setupEchoServerTCP(listenCfg *net.ListenConfig, protocol string, srcIP net.IP, srcPort int) error {
	address := &net.TCPAddr{IP: srcIP, Port: srcPort}
	listener, err := listenCfg.Listen(context.TODO(), protocol, address.String())
	if err != nil {
		return err
	}
	defer listener.Close()

	logger.Infof("Setup TCP echo server on %v", listener.Addr())
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		go func(conn net.Conn) {
			defer conn.Close()

			reader := bufio.NewReader(conn)
			for {
				data, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						logger.Infof("[%v] Disconnected", conn.RemoteAddr())
					} else {
						logger.Infof("[%v] Read data failed, err=%v", conn.RemoteAddr(), err)
					}
					return
				}
				logger.Infof("[%v] -> %s", conn.RemoteAddr(), bytes.Trim(data, "\n"))

				// echo reply
				_, err = conn.Write(data)
				if err != nil {
					if err == io.EOF {
						logger.Infof("[%v] Disconnected", conn.RemoteAddr())
					} else {
						logger.Infof("[%v] Write data failed, err=%v", conn.RemoteAddr(), err)
					}
					return
				}
			}
		}(conn)
	}
}

func setupEchoServerUDP(listenCfg *net.ListenConfig, protocol string, srcIP net.IP, srcPort int) error {
	address := &net.UDPAddr{IP: srcIP, Port: srcPort}
	conn, err := listenCfg.ListenPacket(context.TODO(), protocol, address.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	logger.Infof("Setup UDP echo server on %v", conn.LocalAddr())
	buf := make([]byte, 1500)
	for {
		n, rAddr, err := conn.ReadFrom(buf)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		logger.Infof("[%v] -> %s", rAddr, bytes.Trim(buf[:n], "\n"))

		// echo reply
		_, err = conn.WriteTo(buf[:n], rAddr)
		if err != nil {
			if err == io.EOF {
				logger.Infof("[%v] Disconnected", rAddr)
			} else {
				logger.Infof("[%v] Write data failed, err=%v", rAddr, err)
			}
			return err
		}
	}
}
