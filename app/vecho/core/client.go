package core

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/vtools/app/vecho/logger"
)

type clientConn interface {
	io.Reader
	io.Writer
	io.Closer
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

func SetupEchoClient(protocol string, srcIP net.IP, srcPort int, dstIP net.IP, dstPort int, args []string) error {
	var (
		conn clientConn
		err  error
	)

	switch protocol {
	case "tcp", "tcp4", "tcp6":
		conn, err = net.DialTCP(protocol, &net.TCPAddr{IP: srcIP, Port: srcPort}, &net.TCPAddr{IP: dstIP, Port: dstPort})
	case "udp", "udp4", "udp6":
		conn, err = net.DialUDP(protocol, &net.UDPAddr{IP: srcIP, Port: srcPort}, &net.UDPAddr{IP: dstIP, Port: dstPort})
	default:
		return errors.New("unknown protocol")
	}

	if err != nil {
		return err
	}
	defer conn.Close()

	// fast echo mode
	if len(args) != 0 {
		for _, arg := range args {
			_, _ = conn.Write(append([]byte(arg), '\n'))
		}
		return nil
	}

	logger.Infof("Connect to server[%v] on [%v]", conn.RemoteAddr(), conn.LocalAddr())
	go func() {
		reader := bufio.NewReader(conn)
		for {
			data, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					logger.Infof("[%v -> %v] Disconnected!", conn.LocalAddr(), conn.RemoteAddr())
				} else {
					logger.Infof("[%v -> %v] Read data failed, err=%v", conn.LocalAddr(), conn.RemoteAddr(), err)
				}
				return
			}
			logger.Infof("[%v] <- %v", conn.LocalAddr(), strings.TrimRight(string(data), "\n"))
		}
	}()

	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		data := scan.Bytes()
		_, err := conn.Write(append(data, '\n'))
		if err != nil {
			if err == io.EOF {
				logger.Infof("[%v -> %v] Disconnected!", conn.LocalAddr(), conn.RemoteAddr())
				return nil
			}
			return fmt.Errorf("[%v -> %v] Write data failed, err=%v", conn.LocalAddr(), conn.RemoteAddr(), err)
		}
	}
	return nil
}
