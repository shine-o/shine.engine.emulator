// Package cmd various CLI commands related to the login service
package cmd

import (
	"github.com/shine-o/shine.engine.emulator/internal/app/login"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Update database schema",
	Run:   login.Migrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().Bool("fixtures", false, "load dummy data. WARNING: it wil purge the database")
}
