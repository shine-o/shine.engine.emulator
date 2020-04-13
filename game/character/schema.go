package character

import (
	"errors"
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

func CreateTables(db *pg.DB) error {
	schemaTx,err := db.Begin()
	if err != nil {
		return err
	}
	defer schemaTx.Close()
	for _, model := range []interface{}{
		(*Character)(nil),
		(*Appearance)(nil),
		(*Attributes)(nil),
		(*Location)(nil),
		(*Inventory)(nil),
		(*EquippedItems)(nil),
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

func DeleteTables(db *pg.DB) error {
	purgeTx, err := db.Begin()
	if err != nil {
		return err
	}
	defer purgeTx.Close()

	for _, model := range []interface{}{
		(*Character)(nil),
		(*Appearance)(nil),
		(*Attributes)(nil),
		(*Location)(nil),
		(*Inventory)(nil),
		(*EquippedItems)(nil),
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