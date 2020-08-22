// Package cmd used for various command configs
package cmd

import (
	ps "github.com/shine-o/shine.engine.emulator/internal/app/packet-sniffer"
	"github.com/spf13/cobra"
)

// captureCmd represents the capture command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode file with packet data",
	Run:   ps.Capture,
}

func init() {
	rootCmd.AddCommand(decodeCmd)
}
