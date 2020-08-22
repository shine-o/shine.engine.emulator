package world

import (
	"context"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game/character"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

// Migrate schemas and models
func Migrate(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log = logger.Init("Database Logger", true, false, ioutil.Discard)
	log.Info("Database Logger Migrate()")
	db := database.Connection(ctx, database.ConnectionParams{
		User:     viper.GetString("database.postgres.db_user"),
		Password: viper.GetString("database.postgres.db_password"),
		Host:     viper.GetString("database.postgres.host"),
		Port:     viper.GetString("database.postgres.port"),
		Database: viper.GetString("database.postgres.db_name"),
		Schema:   viper.GetString("database.postgres.schema"),
	})

	if err := database.CreateSchema(db, "world"); err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	if yes, err := cmd.Flags().GetBool("fixtures"); err != nil {
		log.Error(err)
	} else {
		if yes {
			err := character.DeleteTables(db)
			if err != nil {
				log.Fatal(err)
			}
			err = character.CreateTables(db)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = character.CreateTables(db)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
