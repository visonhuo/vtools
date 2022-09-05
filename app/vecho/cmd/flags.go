package cmd

import (
	"math"
)

var (
	protocol   string
	listenIP   string
	listenPort int
	remoteIP   string
	remotePort int

	// extra options
	reuseAddr bool
	reusePort bool
)

func init() {
	serverCmd.Flags().StringVarP(&protocol, "protocol", "p", "tcp",
		"Protocol of echo server/client, default value: tcp, ranges: tcp/udp/tcp4/udp4/tcp6/tcp6")
	serverCmd.Flags().StringVarP(&listenIP, "listen_ip", "l", "0.0.0.0",
		"Listen IP for echo server, default value: 0.0.0.0")
	serverCmd.Flags().IntVarP(&listenPort, "listen_port", "L", 0,
		"Listen port for echo server, default value: 0")
	serverCmd.Flags().BoolVarP(&reuseAddr, "reuse_addr", "", false,
		"Set SO_REUSERADDR sock option value")
	serverCmd.Flags().BoolVarP(&reusePort, "reuse_port", "", false,
		"Set SO_REUSERPORT sock option value")

	clientCmd.Flags().StringVarP(&protocol, "protocol", "p", "tcp",
		"Protocol of echo server/client, default value: tcp, ranges: tcp/udp/tcp4/udp4/tcp6/tcp6")
	clientCmd.Flags().StringVarP(&listenIP, "listen_ip", "l", "0.0.0.0",
		"Listen IP for echo server, default value: 0.0.0.0")
	clientCmd.Flags().IntVarP(&listenPort, "listen_port", "L", 0,
		"Listen port for echo server, default value: 0")
	clientCmd.Flags().StringVarP(&remoteIP, "remote_ip", "r", "127.0.0.1",
		"Request remote IP for echo client, default value: 127.0.0.1")
	clientCmd.Flags().IntVarP(&remotePort, "remote_port", "R", 80,
		"Request remote port for echo client, default value: 80")
}

func isValidProtocol(protocol string) bool {
	switch protocol {
	case "tcp", "udp", "tcp4", "udp4", "tcp6", "udp6":
		return true
	default:
		return false
	}
}

func isValidPort(port int) bool {
	if port < 0 || port > math.MaxUint16 {
		return false
	}
	return true
}
