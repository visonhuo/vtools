package core

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/vtools/app/vecho/logger"
	"golang.org/x/net/ipv4"
)

type clientConn interface {
	io.Reader
	io.Writer
	io.Closer
	LocalAddr() net.Addr
	RemoteAddr() net.Addr
}

func SetupEchoClient(protocol string, srcIP net.IP, srcPort int, dstIP net.IP, dstPort int, zone string, args []string) error {
	var (
		conn clientConn
		err  error
	)

	switch protocol {
	case "tcp", "tcp4", "tcp6":
		conn, err = net.DialTCP(protocol, &net.TCPAddr{IP: srcIP, Port: srcPort, Zone: zone}, &net.TCPAddr{IP: dstIP, Port: dstPort, Zone: zone})
	case "udp", "udp4", "udp6":
		conn, err = net.DialUDP(protocol, &net.UDPAddr{IP: srcIP, Port: srcPort, Zone: zone}, &net.UDPAddr{IP: dstIP, Port: dstPort, Zone: zone})
	case "ip4:udp":
		conn, err = dialIPv4UDP(protocol, &net.UDPAddr{IP: srcIP, Port: srcPort}, &net.UDPAddr{IP: dstIP, Port: dstPort})
	default:
		return fmt.Errorf("invalid protocol(%v)", protocol)
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

type ipv4UDPClientConn struct {
	srcAddr *net.UDPAddr
	dstAddr *net.UDPAddr
	rawConn *ipv4.RawConn
	readBuf []byte
}

func dialIPv4UDP(protocol string, lAddr *net.UDPAddr, rAddr *net.UDPAddr) (*ipv4UDPClientConn, error) {
	if lAddr.Port == 0 {
		lAddr.Port = randomPort()
	}

	packetConn, err := net.ListenPacket(protocol, lAddr.IP.String())
	if err != nil {
		return nil, fmt.Errorf("create packet conn failed, err=%w", err)
	}
	defer func() {
		if err != nil && packetConn != nil {
			_ = packetConn.Close()
		}
	}()

	rawConn, err := ipv4.NewRawConn(packetConn)
	if err != nil {
		return nil, fmt.Errorf("new raw conn failed, err=%w", err)
	}

	return &ipv4UDPClientConn{
		srcAddr: lAddr,
		dstAddr: rAddr,
		rawConn: rawConn,
		readBuf: make([]byte, 1500),
	}, nil
}

func (c *ipv4UDPClientConn) Read(bytes []byte) (n int, err error) {
	// TODO read operation still got some issue
	return c.rawConn.Read(bytes)
	//ipv4Header, payload, _, err := c.rawConn.ReadFrom(c.readBuf)
	//if err != nil {
	//	logger.Infof("read from IP conn failed, err=%v", err)
	//	return
	//}
	//// check packet source
	//if ipv4Header.Src.Equal(c.dstAddr.IP) {
	//	logger.Infof("IP mismatch, expect=%v, actual=%v", c.dstAddr.IP, ipv4Header.Src)
	//	return
	//}
	//
	//newPacket := gopacket.NewPacket(payload, layers.LayerTypeUDP, gopacket.Default)
	//if newPacket == nil {
	//	logger.Infof("Packet decode failed, payload=%v", payload)
	//	return
	//}
	//
	//udpPacket := newPacket.Layer(layers.LayerTypeUDP).(*layers.UDP)
	//if int(udpPacket.DstPort) != c.srcAddr.Port || int(udpPacket.SrcPort) != c.dstAddr.Port {
	//	logger.Infof("Port mismatch, srcPort=%v, dstPort=%v", udpPacket.SrcPort, udpPacket.DstPort)
	//	return
	//}
	//
	//logger.Infof("Recv payload data: %v", udpPacket.Payload)
	//n = copy(bytes, udpPacket.Payload)
	//return
}

func (c *ipv4UDPClientConn) Write(bytes []byte) (int, error) {
	ipv4Header, payload, err := c.generatePacket(bytes)
	if err != nil {
		return -1, fmt.Errorf("generate packet failed, err=%w", err)
	}

	_ = c.rawConn.SetWriteDeadline(time.Now().Add(100 * time.Millisecond))
	if err := c.rawConn.WriteTo(ipv4Header, payload, nil); err != nil {
		return -1, fmt.Errorf("write message failed, err=%w", err)
	}
	return len(bytes), nil
}

func (c *ipv4UDPClientConn) generatePacket(data []byte) (*ipv4.Header, []byte, error) {
	const defaultTTL = 64
	udpHeader := &layers.UDP{
		SrcPort: layers.UDPPort(c.srcAddr.Port),
		DstPort: layers.UDPPort(c.dstAddr.Port),
	}
	_ = udpHeader.SetNetworkLayerForChecksum(&layers.IPv4{
		SrcIP:    c.srcAddr.IP,
		DstIP:    c.dstAddr.IP,
		Protocol: layers.IPProtocolUDP,
		TTL:      defaultTTL,
	})

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths:       true,
	}
	if err := gopacket.SerializeLayers(buf, opts, udpHeader, gopacket.Payload(data)); err != nil {
		return nil, nil, fmt.Errorf("serialize layers failed, err=%w", err)
	}

	udpPacketBytes := buf.Bytes()
	ipv4Header := &ipv4.Header{
		Version:  ipv4.Version,
		Len:      ipv4.HeaderLen,
		TOS:      0x00,
		TotalLen: ipv4.HeaderLen + len(udpPacketBytes),
		Flags:    ipv4.DontFragment,
		FragOff:  0,
		TTL:      defaultTTL,
		Protocol: int(layers.IPProtocolUDP),
		Src:      c.srcAddr.IP,
		Dst:      c.dstAddr.IP,
	}
	return ipv4Header, udpPacketBytes, nil
}

func (c *ipv4UDPClientConn) Close() error {
	return c.rawConn.Close()
}

func (c *ipv4UDPClientConn) LocalAddr() net.Addr {
	return c.srcAddr
}

func (c *ipv4UDPClientConn) RemoteAddr() net.Addr {
	return c.dstAddr
}
