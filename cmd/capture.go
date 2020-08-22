// Package cmd used for various command configs
package cmd

import (
	"github.com/shine-o/shine.engine.packet-sniffer/service"
	"github.com/spf13/cobra"
)

// captureCmd represents the capture command
var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Start capturing and decoding packets",
	Run:   service.Capture,
}

func init() {
	rootCmd.AddCommand(captureCmd)
}
