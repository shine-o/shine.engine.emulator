// Package cmd various CLI commands related to the login service
package cmd

import (
	"github.com/shine-o/shine.engine.emulator/internal/login"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Listen for login connections",
	Long:  `The purpose of the login service is to handle packets related to user account login and server selection.`,
	Run:   login.Start,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	log.Info("serve init()")
}
