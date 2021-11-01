// Package cmd used for various command configs
package cmd

import (
	ps "github.com/shine-o/shine.engine.emulator/internal/packet-sniffer"
	"github.com/spf13/cobra"
)

// captureCmd represents the capture command
var captureCmd = &cobra.Command{
	Use:   "capture",
	Short: "Start capturing and decoding packets",
	Run:   ps.Capture,
}

func init() {
	rootCmd.AddCommand(captureCmd)
}
