// Package cmd various CLI commands related to the login service
package cmd

import (
	wm "github.com/shine-o/shine.engine.emulator/internal/app/world-master"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Listen for world connections to the master",
	Long:  `The purpose of the world master service is to coordinate the registered worlds.`,
	Run:   wm.Start,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	log.Info("serve init()")
}
