// Package cmd various CLI commands related to the service service
package cmd

import (
	"github.com/shine-o/shine.engine.emulator/internal/world"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Listen for World connections",
	Run:   world.Start,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	log.Info("serve init()")
}
