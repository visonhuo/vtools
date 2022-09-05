package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vecho",
	Short: "Echo server/client tools",
	Long: `vecho is a CLI personal experiment tools, it can support 
to setup a echo server or a echo client. Support both TCP and UDP protocol.

Command example:
- vecho server --protocol=tcp --listen_ip=127.0.0.1
- vecho server --protocol=udp --listen_ip=0.0.0.0 --listen_port=64886
- vecho client tcp --remote_ip=192.168.0.1 --remote_port=54666`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(cmd.Long)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vecho.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
