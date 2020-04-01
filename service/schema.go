package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"time"
)

var db *pg.DB

// User model for schema: accounts
type User struct {
	tableName struct{} `pg:"accounts.users"`
	ID        uint64
	UserName  string
	Password  string
	DeletedAt time.Time `pg:"soft_delete"`
}

// Migrate schemas and models
func Migrate(cmd *cobra.Command, args []string) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log = logger.Init("Database Logger", true, false, ioutil.Discard)
	log.Info("Database Logger Migrate()")
	db := dbConn(ctx, "accounts")
	defer db.Close()
	if yes, err := cmd.Flags().GetBool("fixtures"); err != nil {
		log.Error(err)
	} else {
		if yes {
			purge(db)
			createSchema(db)
			fixtures(db)
		} else {
			createSchema(db)
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
		ApplicationName: "login",
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

func createSchema(db *pg.DB) {
	db.Exec("CREATE SCHEMA IF NOT EXISTS accounts;")

	for _, model := range []interface{}{
		(*User)(nil),
	} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func purge(db *pg.DB) {
	for _, model := range []interface{}{
		(*User)(nil),
	} {
		err := db.DropTable(model, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func fixtures(db *pg.DB) {
	password := md5Hash("admin")
	err := db.Insert(&User{
		UserName: "admin",
		Password: password,
	})

	if err != nil {
		log.Fatal(err)
	}
}
