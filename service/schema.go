package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/google/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

var database *gorm.DB

type User struct {
	gorm.Model
	UserName string `gorm:"type:varchar(260);unique_index;not null"`
	Password string `gorm:"type:varchar(36);not null"`
}

func (User) TableName() string {
	return "accounts.user"
}

func Migrate(cmd *cobra.Command, args []string) {
	log = logger.Init("Database Logger", true, false, ioutil.Discard)
	log.Info("Database Logger Migrate()")
	initDatabase()
	defer database.Close()
	if yes, err := cmd.Flags().GetBool("fixtures"); err != nil {
		log.Error(err)
	} else {
		if yes {
			purge()
			migrations()
			fixtures()
		} else {
			migrations()
		}
	}
}

func initDatabase() {
	var (
		dbUser     = viper.GetString("database.postgres.db_user")
		dbPassword = viper.GetString("database.postgres.db_password")
		host       = viper.GetString("database.postgres.host")
		port       = viper.GetString("database.postgres.port")
		dbName     = viper.GetString("database.postgres.db_name")
	)
	dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", dbUser, dbPassword, host, port, dbName)
	if db, err := gorm.Open("postgres", dsn); err != nil {
		log.Fatal(err)
	} else {
		database = db
	}
	log.Infof("connected to the database postgres://v:v@%v:%v/%v?sslmode=disable", host, port, dbName)
}

func migrations() {
	database.Exec("CREATE SCHEMA IF NOT EXISTS accounts;")
	database.AutoMigrate(&User{})
}

func purge() {
	database.DropTableIfExists(&User{})
}

func fixtures() {
	password := md5Hash("admin")
	database.Create(&User{
		UserName: "admin",
		Password: password,
	})
}

func md5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
