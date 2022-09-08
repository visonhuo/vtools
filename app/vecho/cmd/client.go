package cmd

import (
	"net"

	"github.com/spf13/cobra"
	"github.com/vtools/app/vecho/core"
	"github.com/vtools/app/vecho/logger"
)

var (
	IPv4loopback = net.IPv4(127, 0, 0, 1)
)

var clientFlags = struct {
	protocol   string
	listenIP   net.IP
	listenPort int
	remoteIP   net.IP
	remotePort int
	zone       string
}{}

func init() {
	clientCmd.Flags().StringVarP(&clientFlags.protocol, "protocol", "p", "tcp",
		"Protocol of echo server/client, default value: tcp, ranges: tcp/udp/tcp4/udp4/tcp6/tcp6")
	clientCmd.Flags().IPVarP(&clientFlags.listenIP, "listen_ip", "l", net.IPv4zero,
		"Listen IP for echo server, default value: (ipv4)0.0.0.0 / (ipv6)[::]")
	clientCmd.Flags().IntVarP(&clientFlags.listenPort, "listen_port", "L", 0,
		"Listen port for echo server, default value: 0")
	clientCmd.Flags().IPVarP(&clientFlags.remoteIP, "remote_ip", "r", IPv4loopback,
		"Request remote IP for echo client, default value: (ipv1)127.0.0.1 / (ipv6)[::1]")
	clientCmd.Flags().IntVarP(&clientFlags.remotePort, "remote_port", "R", 80,
		"Request remote port for echo client, default value: 80")
	clientCmd.Flags().StringVarP(&clientFlags.zone, "zone", "z", "",
		"Address zone for echo client (ipv6 only), default value: en0")

	rootCmd.AddCommand(clientCmd)
}

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Echo client tools",
	Long: `A simple tool collection of echo client (support tcp4/tcp6/udp4/udp6/ip4:udp).

Command example:
- vecho client --remote_ip=192.168.0.1 --remote_port=54996
- vecho client --protocol=udp --listen_ip=0.0.0.0 --listen_port=64886 --remote_ip=192.168.0.1 --remote_port=54996
- vecho client -p=udp -l=0.0.0.0 -L=64886 -r=192.168.0.1 -R=54996
- vecho client -p=ip4:udp -L=8888 -R=9999 (Need root permission)
Fast send mode example:
- vecho client echo_content -p=udp -R=54996`,
	Run: func(cmd *cobra.Command, args []string) {
		if clientFlags.listenIP == nil {
			logger.Fatalf("Error: invalid listen IP value(%v)", clientFlags.listenIP)
		} else if isIPv6Protocol(clientFlags.protocol) && clientFlags.listenIP.Equal(net.IPv4zero) {
			clientFlags.listenIP = net.IPv6zero
		}

		if !isValidPort(clientFlags.listenPort) {
			logger.Fatalf("Error: invalid listen port value(%v)", clientFlags.listenPort)
		}

		if clientFlags.remoteIP == nil {
			logger.Fatalf("Error: invalid remote IP value(%v)", clientFlags.remoteIP)
		} else if isIPv6Protocol(clientFlags.protocol) && clientFlags.remoteIP.Equal(IPv4loopback) {
			clientFlags.remoteIP = net.IPv6loopback
		}

		if !isValidPort(clientFlags.remotePort) {
			logger.Fatalf("Error: invalid remote port value(%v)", clientFlags.remotePort)
		}

		if isIPv6Protocol(clientFlags.protocol) && clientFlags.zone == "" {
			serverFlags.zone = "en0"
		}

		// srcPort == listenPort
		// dstPort == remotePort
		if err := core.SetupEchoClient(
			clientFlags.protocol,
			clientFlags.listenIP,
			clientFlags.listenPort,
			clientFlags.remoteIP,
			clientFlags.remotePort,
			clientFlags.zone,
			args,
		); err != nil {
			logger.Fatalf("Error: setup echo client failed(%v,%v:%v,%v:%v), err: %v",
				clientFlags.protocol, clientFlags.listenIP, clientFlags.listenPort, clientFlags.remoteIP, clientFlags.remotePort, err)
		}
		logger.Infof("Vecho client closed~")
	},
}
