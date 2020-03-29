package service

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/google/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"time"
)

// User model for schema: accounts
type Character struct {
	tableName     struct{} `pg:"world.characters"`
	ID            uint64
	//Appearance    *CharacterAppearance    `pg:"unique,on_delete:RESTRICT"`
	Appearance    *CharacterAppearance
	Attributes    *CharacterAttributes
	Location      *CharacterLocation
	Inventory     *CharacterInventory
	EquippedItems *CharacterEquippedItems
	UserID        uint64                  `pg:"notnull"`
	AdminLevel    uint8
	Slot          uint8  `pg:"notnull"`
	OnMap         string `pg:"notnull"`
	IsDeleted     bool
	DeletedAt     time.Time `pg:"soft_delete"`
}

type CharacterAppearance struct {
	tableName   struct{} `pg:"world.character_appearance"`
	ID          uint64
	CharacterID uint64 //
	Character   *Character
	//Race      uint8
	Class     uint8     `pg:"notnull"`
	Gender    uint8     `pg:"notnull"`
	HairType  uint8     `pg:"notnull"`
	HairColor uint8     `pg:"notnull"`
	FaceType  uint8     `pg:"notnull"`
	DeletedAt time.Time `pg:"soft_delete"`
}

type CharacterAttributes struct {
	tableName    struct{} `pg:"world.character_attributes"`
	ID           uint64
	CharacterID  uint64
	Character    *Character  `pg:"unique"`
	Level        uint8 `pg:"notnull"`
	Experience   uint64`pg:"notnull"`
	Fame         uint64`pg:"notnull"`
	Hp           uint32`pg:"notnull"`
	Sp           uint32`pg:"notnull"`
	Intelligence uint32`pg:"notnull"`
	Strength     uint32`pg:"notnull"`
	Dexterity    uint32`pg:"notnull"`
	Endurance    uint32`pg:"notnull"`
	Spirit       uint32`pg:"notnull"`
	Money        uint64`pg:"notnull"`
	KillPoints   uint32`pg:"notnull"`
	HpStones     uint16`pg:"notnull"`
	SpStones     uint16`pg:"notnull"`
	DeletedAt    time.Time `pg:"soft_delete"`
}

type CharacterLocation struct {
	tableName   struct{} `pg:"world.character_location"`
	ID          uint64
	CharacterID uint64 //
	Character   *Character
	MapName     string    `pg:"notnull"`
	X           uint32    `pg:"notnull"`
	Y           uint32    `pg:"notnull"`
	D           uint32    `pg:"notnull"`
	IsKQ        bool      `pg:"notnull"`
	DeletedAt   time.Time `pg:"soft_delete"`
}

type CharacterInventory struct {
	tableName   struct{} `pg:"world.character_inventory"`
	ID          uint64
	CharacterID uint64 //
	Character   *Character
	ShnID       uint16    `pg:"notnull"`
	Slot        uint16    `pg:"unique,notnull"`
	IsStack     bool      `pg:"notnull"`
	IsStored    bool      `pg:"notnull"`
	FromMonarch bool      `pg:"notnull"`
	FromStore   bool      `pg:"notnull"`
	DeletedAt   time.Time `pg:"soft_delete"`
}

type CharacterEquippedItems struct {
	tableName        struct{} `pg:"world.character_equipped_items"`
	ID               uint64
	CharacterID      uint64 //
	Character        *Character
	Head             uint16
	Body             uint16
	Boots            uint16
	LeftHand         uint16
	RightHand        uint16
	MiniPet          uint16
	ApparelHead      uint16
	ApparelEye       uint16
	ApparelBody      uint16
	ApparelBoots     uint16
	ApparelLeftHand  uint16
	ApparelRightHand uint16
	ApparelBack      uint16
	ApparelAura      uint16
	ApparelShield    uint16
	DeletedAt        time.Time `pg:"soft_delete"`
}

// Migrate schemas and models
func Migrate(cmd *cobra.Command, args []string) {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	log = logger.Init("Database Logger", true, false, ioutil.Discard)
	log.Info("Database Logger Migrate()")
	db := dbConn(ctx, "world")
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
			fixtures()
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
		ApplicationName: "world",
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
	for _, model := range []interface{}{
		(*Character)(nil),
		(*CharacterAppearance)(nil),
		(*CharacterAttributes)(nil),
		(*CharacterLocation)(nil),
		(*CharacterInventory)(nil),
		(*CharacterEquippedItems)(nil),
	} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func purge(db *pg.DB) error {
	for _, model := range []interface{}{
		(*Character)(nil),
		(*CharacterAppearance)(nil),
		(*CharacterAttributes)(nil),
		(*CharacterLocation)(nil),
		(*CharacterInventory)(nil),
		(*CharacterEquippedItems)(nil),
	} {
		err := db.DropTable(model, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func fixtures() {

}
