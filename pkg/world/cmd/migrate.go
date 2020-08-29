// Package cmd various CLI commands related to the login service
package cmd

import (
	"github.com/shine-o/shine.engine.emulator/internal/app/world"
	"github.com/spf13/cobra"
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Update database schema",
	Run:   world.Migrate,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	migrateCmd.Flags().Bool("purge", false, "load dummy data. WARNING: it wil purge the database")
}
