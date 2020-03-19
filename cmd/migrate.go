// Package cmd various CLI commands related to the login service
package cmd

import (
	"github.com/shine-o/shine.engine.login/service"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Update database schema",
	Run:   service.Migrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	err := doc.GenMarkdownTree(migrateCmd, "docs")
	if err != nil {
		log.Fatal(err)
	}
	migrateCmd.Flags().Bool("fixtures", false, "load dummy data. WARNING: it wil purge the database")
}
