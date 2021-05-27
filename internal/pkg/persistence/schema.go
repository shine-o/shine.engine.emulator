package persistence

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/database"
)

var db *pg.DB

func InitDB(cp database.ConnectionParams) {
	ctx := context.Background()
	db = database.Connection(ctx, cp)
}

func DB() *pg.DB {
	return db
}

func CloseDB() {
	log.Info(db.Close())
}

// CreateTables if not yet created
func CreateTables() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer closeTx(tx)

	for _, model := range []interface{}{
		(*Character)(nil),
		(*Appearance)(nil),
		(*Attributes)(nil),
		(*Location)(nil),
		(*ClientOptions)(nil),
		(*EquippedItems)(nil),
		(*Item)(nil),
		(*ItemAttributes)(nil),
	} {
		err := tx.Model(model).CreateTable(&orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return fmt.Errorf("%v, %v", err, tx.Rollback())
		}
	}
	return tx.Commit()
}

// DeleteTables if they exist
// TODO: https://github.com/go-pg/migrations#example
func DeleteTables() error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer tx.Close()

	for _, model := range []interface{}{
		(*Character)(nil),
		(*Appearance)(nil),
		(*Attributes)(nil),
		(*Location)(nil),
		(*ClientOptions)(nil),
		(*EquippedItems)(nil),
		(*Item)(nil),
		(*ItemAttributes)(nil),
	} {
		err := tx.Model(model).DropTable(&orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			return fmt.Errorf("%v, %v", err, tx.Rollback())
		}
	}
	return tx.Commit()
}

// // TODO: https://github.com/go-pg/migrations#example
func CleanDB() {
	err := DeleteTables()
	if err != nil {
		log.Fatal(err)
	}
	err = CreateTables()
	if err != nil {
		log.Fatal(err)
	}
}

func closeTx(tx *pg.Tx) {
	err := tx.Close()
	if err != nil {
		log.Error(err)
	}
}
