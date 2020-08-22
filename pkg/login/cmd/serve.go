// Package cmd various CLI commands related to the login service
package cmd

import (
	"github.com/shine-o/shine.engine.emulator/internal/app/login"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
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
	err := doc.GenMarkdownTree(serveCmd, "docs")
	if err != nil {
		log.Fatal(err)
	}
}
