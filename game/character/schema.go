package character

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
)

// CreateTables if not yet created
func CreateTables(db *pg.DB) error {
	createTx, err := db.Begin()
	if err != nil {
		return err
	}
	defer createTx.Close()
	for _, model := range []interface{}{
		(*Character)(nil),
		(*Appearance)(nil),
		(*Attributes)(nil),
		(*Location)(nil),
		(*Inventory)(nil),
		(*EquippedItems)(nil),
	} {
		err := createTx.CreateTable(model, &orm.CreateTableOptions{
			IfNotExists:   true,
			FKConstraints: true,
		})
		if err != nil {
			return fmt.Errorf("%v, %v", err, createTx.Rollback())
		}
	}
	return createTx.Commit()
}

// DeleteTables if they exist
func DeleteTables(db *pg.DB) error {
	deleteTx, err := db.Begin()
	if err != nil {
		return err
	}
	defer deleteTx.Close()

	for _, model := range []interface{}{
		(*Character)(nil),
		(*Appearance)(nil),
		(*Attributes)(nil),
		(*Location)(nil),
		(*Inventory)(nil),
		(*EquippedItems)(nil),
	} {
		err := deleteTx.DropTable(model, &orm.DropTableOptions{
			IfExists: true,
			Cascade:  true,
		})
		if err != nil {
			return fmt.Errorf("%v, %v", err, deleteTx.Rollback())
		}
	}
	return deleteTx.Commit()
}
