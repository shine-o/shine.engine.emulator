package service

import (
	"fmt"
	"github.com/google/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // GORM needs this
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"time"
)

var database *gorm.DB

type Model struct {
	ID        uint64 `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// User model for schema: accounts
type Character struct {
	Model
	Appearance CharacterAppearance
	Attributes CharacterAttributes
	Location   CharacterLocation
	Inventory CharacterInventory
	EquippedItems CharacterEquippedItems
	UserID     uint64
	AdminLevel uint8
	slot       uint8
	OnMap      string
	IsDeleted  bool
}

type CharacterAppearance struct {
	Model
	CharacterID uint16
	Race      uint8
	Class     uint8
	Gender    uint8
	HairType  uint8
	HairColor uint8
	FaceType  uint8
}

type CharacterAttributes struct {
	Model
	CharacterID uint16
	Level        uint8
	Experience   uint64
	Fame         uint64
	Hp           uint32
	Sp           uint32
	Intelligence uint32
	Strength     uint32
	Dexterity    uint32
	Endurance    uint32
	Spirit       uint32
	Money        uint64
	KillPoints   uint32
	HpStones     uint16
	SpStones     uint16
}

type CharacterLocation struct {
	Model
	CharacterID uint16
	MapName string
	X       uint32
	Y       uint32
	D       uint32
	isKQ    bool
}

type CharacterInventory struct {
	Model
	CharacterID uint16
	ShnID uint16
	Slot  uint16
	IsStack bool
	IsStored bool
	FromMonarch bool
	FromStore bool
}

type CharacterEquippedItems struct {
	Model
	CharacterID uint16
	Head uint16
	Body uint16
	Boots uint16
	LeftHand uint16
	RightHand uint16
	MiniPet uint16
	ApparelHead uint16
	ApparelEye uint16
	ApparelBody uint16
	ApparelBoots uint16
	ApparelLeftHand uint16
	ApparelRightHand uint16
	ApparelBack uint16
	ApparelAura uint16
	ApparelShield uint16
}

// TableName schema identifier
func (Character) TableName() string {
	return "world.character"
}

// TableName schema identifier
func (CharacterAppearance) TableName() string {
	return "world.character_appearance"
}

// TableName schema identifier
func (CharacterAttributes) TableName() string {
	return "world.character_attributes"
}

// TableName schema identifier
func (CharacterLocation) TableName() string {
	return "world.character_location"
}

// TableName schema identifier
func (CharacterInventory) TableName() string {
	return "world.character_inventory"
}

// TableName schema identifier
func (CharacterEquippedItems) TableName() string {
	return "world.character_equipped_items"
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
			purge()

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
	database.Exec("CREATE SCHEMA IF NOT EXISTS world;")
	database.AutoMigrate(&Character{})
	database.AutoMigrate(&CharacterAppearance{})
	database.AutoMigrate(&CharacterAttributes{})
	database.AutoMigrate(&CharacterLocation{})
	database.AutoMigrate(&CharacterInventory{})
	database.AutoMigrate(&CharacterEquippedItems{})

	//database.Model(&Character{}).AddForeignKey("")
}

func purge() {
	database.DropTableIfExists(&Character{})
	database.DropTableIfExists(&CharacterAppearance{})
	database.DropTableIfExists(&CharacterAttributes{})
	database.DropTableIfExists(&CharacterLocation{})
	database.DropTableIfExists(&CharacterInventory{})
	database.DropTableIfExists(&CharacterEquippedItems{})

}

func fixtures() {
	//password := md5Hash("admin")
	//database.Create(&User{
	//	UserName: "admin",
	//	Password: password,
	//})
}
