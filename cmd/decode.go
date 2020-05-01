// Package cmd used for various command configs
package cmd

import (
	"github.com/shine-o/shine.engine.packet-sniffer/service"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"log"
)

// captureCmd represents the capture command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode file with packet data",
	Run:   service.Capture,
}

func init() {
	rootCmd.AddCommand(decodeCmd)
	err := doc.GenMarkdownTree(decodeCmd, "docs")
	if err != nil {
		log.Fatal(err)
	}
}
