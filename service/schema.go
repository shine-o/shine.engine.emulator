package service

import (
	"fmt"
	"github.com/google/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // GORM needs this
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
)

var database *gorm.DB

// User model for schema: accounts
type Character struct {
	gorm.Model
}

type CharacterShape struct {
	gorm.Model
}

type CharacterAttributes struct {
	gorm.Model
}

//
//type CharacterPosition struct {
//	gorm.Model
//}

// TableName schema identifier
func (Character) TableName() string {
	return "player.character"
}

// Migrate schemas and models
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
	database.Exec("CREATE SCHEMA IF NOT EXISTS player;")
	database.AutoMigrate(&Character{})
	database.AutoMigrate(&CharacterShape{})
	database.AutoMigrate(&CharacterAttributes{})
}

func purge() {
	database.DropTableIfExists(&Character{})
	database.DropTableIfExists(&CharacterShape{})
	database.DropTableIfExists(&CharacterAttributes{})
}

func fixtures() {
	//password := md5Hash("admin")
	//database.Create(&User{
	//	UserName: "admin",
	//	Password: password,
	//})
}