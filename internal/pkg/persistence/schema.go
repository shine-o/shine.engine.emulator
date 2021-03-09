package persistence

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
)

// CreateTables if not yet created
func CreateTables(db *pg.DB) error {

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
func DeleteTables(db *pg.DB) error {
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

func closeTx(tx *pg.Tx) {
	err := tx.Close()
	if err != nil {
		log.Error(err)
	}
}
