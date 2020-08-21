// Package cmd various CLI commands related to the service service
package cmd

import (
	"github.com/shine-o/shine.engine.world/service"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Listen for World connections",
	Run:   service.Start,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	log.Info("serve init()")
	err := doc.GenMarkdownTree(serveCmd, "docs")
	if err != nil {
		log.Fatal(err)
	}
}
