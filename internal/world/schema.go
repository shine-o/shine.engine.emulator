package world

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Migrate schemas and models
func Migrate(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log.Info("Database Logger Migrate()")
	persistence.InitDB(database.ConnectionParams{
		User:     viper.GetString("database.postgres.db_user"),
		Password: viper.GetString("database.postgres.db_password"),
		Host:     viper.GetString("database.postgres.host"),
		Port:     viper.GetString("database.postgres.port"),
		Database: viper.GetString("database.postgres.db_name"),
		Schema:   viper.GetString("database.postgres.schema"),
	})

	if err := database.CreateSchema(persistence.DB(), "world"); err != nil {
		log.Fatal(err)
	}

	defer persistence.CloseDB()

	if yes, err := cmd.Flags().GetBool("purge"); err != nil {
		log.Error(err)
	} else {
		if yes {
			err := persistence.DeleteTables()
			if err != nil {
				log.Fatal(err)
			}
			err = persistence.CreateTables()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = persistence.CreateTables()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
