// Package cmd used for various command configs
package cmd

import (
	"github.com/shine-o/shine.engine.packet-sniffer/service"
	"github.com/spf13/cobra"
)

// captureCmd represents the capture command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode file with packet data",
	Run:   service.Capture,
}

func init() {
	rootCmd.AddCommand(decodeCmd)
}
