package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vtools/app/vecho/core"
	"github.com/vtools/app/vecho/logger"
	"net"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Echo server tools",
	Long: `A tool collection of echo server.

Command example:
- vecho server --protocol=tcp --listen_ip=127.0.0.1 
- vecho server --protocol=udp --listen_ip=0.0.0.0 --listen_port=64886
- vecho server -p=udp -l=0.0.0.0 -L=64886`,
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

		if err := core.SetupEchoServer(
			protocol, srcIP, listenPort,
			core.WithReuseAddr(reuseAddr),
			core.WithReusePort(reusePort),
		); err != nil {
			logger.Fatalf("Error: setup echo server failed(%v,%v,%v), err: %v", protocol, srcIP, listenPort, err)
		}

		logger.Infof("Vecho server closed~")
	},
}
