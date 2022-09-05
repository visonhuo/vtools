package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vtools/app/vecho/core"
	"github.com/vtools/app/vecho/logger"
	"net"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Echo client tools",
	Long: `A tool collection of echo server.

Command example:
- vecho client --remote_ip=192.168.0.1 --remote_port=54996
- vecho client --protocol=udp4 --listen_ip=0.0.0.0 --listen_port=64886 --remote_ip=192.168.0.1 --remote_port=54996
- vecho client -p=udp -l=0.0.0.0 -L=64886 -r=192.168.0.1 -R=54996

Write content only example:
- vecho client echo_content -p=udp -R=54996`,
	Run: func(cmd *cobra.Command, args []string) {
		if !isValidProtocol(protocol) {
			logger.Fatalf("Error: invalid protocol flag value(%v)", protocol)
		}

		srcIP := net.ParseIP(listenIP)
		if srcIP == nil {
			logger.Fatalf("Error: invalid listen IP value(%v)", listenIP)
		}

		if !isValidPort(listenPort) {
			logger.Fatalf("Error: invalid listen port value(%v)", listenPort)
		}

		dstIP := net.ParseIP(remoteIP)
		if dstIP == nil {
			logger.Fatalf("Error: invalid remote IP value(%v)", remoteIP)
		}

		if !isValidPort(remotePort) {
			logger.Fatalf("Error: invalid remote port value(%v)", remotePort)
		}

		// srcPort == listenPort
		// dstPort == remotePort
		if err := core.SetupEchoClient(protocol, srcIP, listenPort, dstIP, remotePort, args); err != nil {
			logger.Fatalf("Error: setup echo client failed(%v,%v:%v,%v:%v), err: %v",
				protocol, srcIP, listenPort, dstIP, remotePort, err)
		}

		logger.Infof("Vecho client closed~")
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
}
