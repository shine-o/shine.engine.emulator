package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/logger"
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
	db := dbConn(ctx, "service")
	defer db.Close()
	if yes, err := cmd.Flags().GetBool("fixtures"); err != nil {
		log.Error(err)
	} else {
		if yes {
			err := purge(db)
			if err != nil {
				log.Fatal(err)
			}
			err = createSchema(db)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = createSchema(db)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func dbConn(ctx context.Context, schema string) *pg.DB {
	var (
		dbUser     = viper.GetString("database.postgres.db_user")
		dbPassword = viper.GetString("database.postgres.db_password")
		host       = viper.GetString("database.postgres.host")
		port       = viper.GetString("database.postgres.port")
		dbName     = viper.GetString("database.postgres.db_name")
	)

	db := pg.Connect(&pg.Options{
		Addr:            fmt.Sprintf("%v:%v", host, port),
		User:            dbUser,
		Password:        dbPassword,
		Database:        dbName,
		ApplicationName: "service",
		TLSConfig:       nil,
		//DialTimeout:     15,
		//ReadTimeout:     5,
		//WriteTimeout:    5,
		PoolSize:    5,
		PoolTimeout: 5,
	})
	log.Info(db)
	db = db.WithParam(schema, pg.Safe(schema))
	return db.WithContext(ctx)
}

func createSchema(db *pg.DB) error {
	schemaTx,err := db.Begin()
	if err != nil {
		return err
	}
	defer schemaTx.Close()

	_, err = schemaTx.Exec("CREATE SCHEMA IF NOT EXISTS service;")

	if err != nil {
		log.Fatal(err)
	}
	for _, model := range []interface{}{
		(*Character)(nil),
		(*CharacterAppearance)(nil),
		(*CharacterAttributes)(nil),
		(*CharacterLocation)(nil),
		(*CharacterInventory)(nil),
		(*CharacterEquippedItems)(nil),
	} {
		err := schemaTx.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return errors.New(fmt.Sprintf("%v, %v", err, schemaTx.Rollback()))
		}
	}
	return schemaTx.Commit()
}

func purge(db *pg.DB) error {
	purgeTx, err := db.Begin()
	if err != nil {
		return err
	}
	defer purgeTx.Close()

	for _, model := range []interface{}{
		(*Character)(nil),
		(*CharacterAppearance)(nil),
		(*CharacterAttributes)(nil),
		(*CharacterLocation)(nil),
		(*CharacterInventory)(nil),
		(*CharacterEquippedItems)(nil),
	} {
		err := purgeTx.DropTable(model, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			return errors.New(fmt.Sprintf("%v, %v", err, purgeTx.Rollback()))
		}
	}
	return purgeTx.Commit()
}