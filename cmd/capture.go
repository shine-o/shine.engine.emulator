// Package for cmd config
package cmd

import (
	"github.com/shine-o/shine.engine.packet-sniffer/service"
	"github.com/spf13/cobra"
)

// captureCmd represents the capture command
var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Start capturing and decoding packets",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: service.Capture,
}

func init() {
	rootCmd.AddCommand(captureCmd)
}
