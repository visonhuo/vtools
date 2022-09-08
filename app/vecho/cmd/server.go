package cmd

import (
	"net"

	"github.com/spf13/cobra"
	"github.com/vtools/app/vecho/core"
	"github.com/vtools/app/vecho/logger"
)

var serverFlags = struct {
	protocol   string
	listenIP   net.IP
	listenPort int
	zone       string
	// extra options
	reuseAddr bool
	reusePort bool
}{}

func init() {
	serverCmd.Flags().StringVarP(&serverFlags.protocol, "protocol", "p", "tcp",
		"Protocol of echo server/client, default value: tcp, ranges: tcp/udp/tcp4/udp4/tcp6/tcp6")
	serverCmd.Flags().IPVarP(&serverFlags.listenIP, "listen_ip", "l", net.IPv4zero,
		"Listen IP for echo server, default value: (ipv4)0.0.0.0 / (ipv6)[::]")
	serverCmd.Flags().IntVarP(&serverFlags.listenPort, "listen_port", "L", 0,
		"Listen port for echo server, default value: 0")
	serverCmd.Flags().StringVarP(&serverFlags.zone, "zone", "z", "",
		"Address zone for echo server (ipv6 only), default value: en0")
	serverCmd.Flags().BoolVarP(&serverFlags.reuseAddr, "reuse_addr", "", false,
		"Set SO_REUSERADDR sock option value")
	serverCmd.Flags().BoolVarP(&serverFlags.reusePort, "reuse_port", "", false,
		"Set SO_REUSERPORT sock option value")

	rootCmd.AddCommand(serverCmd)
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Echo server tools",
	Long: `A simple tool collection of echo server. (support tcp4/tcp6/udp4/udp6)

Command example:
- vecho server --protocol=tcp --listen_ip=127.0.0.1
- vecho server --protocol=udp --listen_ip=0.0.0.0 --listen_port=64886
- vecho server --protocol=udp6 --listen_ip=fe80::8c:574d:b224:3f2f --listen_port=64886 --zone=en0
- vecho server -p=udp6 -l=fe80::8c:574d:b224:3f2f -L=64886 -z=en0
`,
	Run: func(cmd *cobra.Command, args []string) {
		if serverFlags.listenIP == nil {
			logger.Fatalf("Error: invalid listen IP value(%v)", serverFlags.listenIP)
		} else if isIPv6Protocol(serverFlags.protocol) && serverFlags.listenIP.Equal(net.IPv4zero) {
			// upgrade to IPv6zero if try to use ipv6 protocol
			serverFlags.listenIP = net.IPv6zero
		}

		if !isValidPort(serverFlags.listenPort) {
			logger.Fatalf("Error: invalid listen port value(%v)", serverFlags.listenPort)
		}

		if isIPv6Protocol(serverFlags.protocol) && serverFlags.zone == "" {
			serverFlags.zone = "en0"
		}

		if err := core.SetupEchoServer(
			serverFlags.protocol,
			serverFlags.listenIP,
			serverFlags.listenPort,
			serverFlags.zone,
			core.WithReuseAddr(serverFlags.reuseAddr),
			core.WithReusePort(serverFlags.reusePort),
		); err != nil {
			logger.Fatalf("Error: setup echo server failed(%v,%v,%v), err: %v",
				serverFlags.protocol, serverFlags.listenIP, serverFlags.listenPort, err)
		}

		logger.Infof("Vecho server closed~")
	},
}
