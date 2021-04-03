package persistence

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"time"
)

const (
	BufferInventory    = 100
	EquippedInventory  = 8
	BagInventory       = 9
	MiniHouseInventory = 12

	BagInventoryMin = 0
	BagInventoryMax = 117
)

type Item struct {
	tableName     struct{} `pg:"world.items"`
	ID            uint64
	InventoryType int `pg:",notnull,unique:item"`
	// box 2 = reward  inventory
	// box 3 = mini house furniture
	// box 8 = equipped items  / 1-29
	// box 9 = inventory, storage  // 9216 - 9377 (24 slots per page)
	// box 12 = mini houses // 12288 equipped minihouse, 12299-12322 available slots
	// box 13 = mini house accessories
	// box 14 = mini house tile all inventory
	// box 15 = premium actions inventory(dances)
	// box 16 = mini house mini game inventory
	Slot        int        `pg:",use_zero,notnull,unique:item"`
	CharacterID uint64     `pg:",notnull,unique:item" `
	Character   *Character `pg:"rel:belongs-to"`
	ShnID       uint16     `pg:",notnull"`
	ShnInxName  string     `pg:",notnull"`
	Stackable   bool       `pg:",notnull,use_zero"`
	Amount      int
	Attributes  *ItemAttributes `pg:"rel:belongs-to"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	//DeletedAt time.Time `pg:",soft_delete"`
}

type ItemAttributes struct {
	tableName     struct{} `pg:"world.item_attributes"`
	ID            uint64   `pg:",unique:item"`
	ItemID        uint64   `pg:",use_zero,notnull,unique:item"`
	StrengthBase  int      `pg:",use_zero"`
	StrengthExtra int      `pg:",use_zero"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	//DeletedAt time.Time `pg:",soft_delete"`
}

type ItemParams struct {
	CharacterID uint64
	ShnID       uint16
	Stackable   bool
	Amount      int
	Attributes  *ItemAttributes
}

// Insert an Item
// Items can only be persisted if they are linked to a character
func (i *Item) Insert() error {

	err := validateItem(i)

	if err != nil {
		return err
	}

	i.InventoryType = BagInventory

	slot, err := freeSlot(i.CharacterID, BagInventory)

	if err != nil {
		return err
	}

	i.Slot = slot

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer closeTx(tx)

	i.CreatedAt = time.Now()

	_, err = tx.Model(i).
		Returning("*").
		Insert()

	if err != nil {
		return err
	}

	if i.Attributes == nil {
		i.Attributes = &ItemAttributes{}
	}

	i.Attributes.ItemID = i.ID
	i.Attributes.CreatedAt = time.Now()

	_, err = tx.Model(i.Attributes).
		Returning("*").
		Insert()

	if err != nil {
		return errors.Err{
			Code: errors.PersistenceErrDB,
			Details: errors.ErrDetails{
				"err":   err,
				"txErr": tx.Rollback(),
			},
		}
	}

	return tx.Commit()

}
func (i *Item) Update() error {

	err := validateItem(i)

	if err != nil {
		return err
	}

	if i.ShnID != i.ShnID {
		return errors.Err{
			Code: errors.PersistenceErrItemDistinctShnID,
			Details: errors.ErrDetails{
				"originalShnID": i.ShnID,
				"newShnID":      i.ShnID,
			},
		}
	}

	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer closeTx(tx)

	i.UpdatedAt = time.Now()

	_, err = tx.Model(i).
		WherePK().
		Update()

	if err != nil {
		return errors.Err{
			Code: errors.PersistenceErrDB,
			Details: errors.ErrDetails{
				"err":   err,
				"txErr": tx.Rollback(),
			},
		}
	}

	// if not exists, update
	if i.Attributes == nil {
		attr := &ItemAttributes{}
		attr.ID = i.ID
		_, err = tx.Model(attr).Insert()
	} else {
		_, err = tx.Model(i.Attributes).
			WherePK().
			Update()
	}

	if err != nil {
		return errors.Err{
			Code: errors.PersistenceErrDB,
			Details: errors.ErrDetails{
				"err":   err,
				"txErr": tx.Rollback(),
			},
		}
	}

	return tx.Commit()
}
func (i Item) MoveTo(inventoryType int, slot int) (*Item, error) {
	var (
		otherItem Item
	)
	tx, err := db.Begin()

	if err != nil {
		return &i, err
	}

	defer closeTx(tx)

	err = tx.Model(&otherItem).
		Where("character_id = ?", i.CharacterID).
		Where("inventory_type = ?", inventoryType).
		Where("slot = ?", slot).
		Select()

	if err == nil {
		var (
			// values to be used by otherItem
			originalInventory = i.InventoryType
			originalSlot      = i.Slot
		)

		// to avoid unique constraint error, use buffer inventory
		i.InventoryType = BufferInventory
		i.Slot = 0
		_, err := tx.Model(&i).WherePK().Update()
		if err != nil {
			return &i, errors.Err{
				Code: errors.PersistenceErrItemSlotUpdate,
				Details: errors.ErrDetails{
					"err":           err,
					"from":          i.Slot,
					"to":            otherItem.Slot,
					"fromInventory": i.InventoryType,
					"toInventory":   inventoryType,
					"shnID":         i.ShnID,
				},
			}
		}

		// there is an item there,
		// switch slot ids
		otherItem.InventoryType = originalInventory
		otherItem.Slot = originalSlot

		_, err = tx.Model(&otherItem).WherePK().Update()
		if err != nil {
			return &i, errors.Err{
				Code: errors.PersistenceErrItemSlotUpdate,
				Details: errors.ErrDetails{
					"err":           err,
					"from":          i.Slot,
					"to":            otherItem.Slot,
					"fromInventory": i.InventoryType,
					"toInventory":   inventoryType,
					"shnID":         i.ShnID,
				},
			}
		}
	}

	i.InventoryType = inventoryType
	i.Slot = slot

	_, err = tx.Model(&i).WherePK().Update()

	if err != nil {
		return &i, errors.Err{
			Code: errors.PersistenceErrItemSlotUpdate,
			Details: errors.ErrDetails{
				"err":           err,
				"txErr":         tx.Rollback(),
				"from":          i.Slot,
				"to":            otherItem.Slot,
				"fromInventory": i.InventoryType,
				"toInventory":   inventoryType,
				"shnID":         i.ShnID,
			},
		}
	}

	return &i, tx.Commit()
}

func DeleteItem(itemID int) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer closeTx(tx)

	_, err = tx.Model((*ItemAttributes)(nil)).Where("item_id = ?", itemID).Delete()

	if err != nil {
		return errors.Err{
			Code: errors.PersistenceErrDB,
			Details: errors.ErrDetails{
				"err":    err,
				"txErr":  tx.Rollback(),
				"itemID": itemID,
			},
		}
	}

	_, err = tx.Model((*Item)(nil)).Where("id = ?", itemID).Delete()

	if err != nil {
		return errors.Err{
			Code: errors.PersistenceErrDB,
			Details: errors.ErrDetails{
				"err":    err,
				"txErr":  tx.Rollback(),
				"itemID": itemID,
			},
		}
	}

	return tx.Commit()
}

func GetItemWhere(clauses map[string]interface{}, deleted bool) (*Item, error) {
	var item Item

	query := db.Model(&item)

	for k, v := range clauses {
		query.Where(k, v)
	}

	//if !deleted {
	//	query.Where("deleted_at = null")
	//}

	err := query.Select()

	return &item, err
}

func GetCharacterItems(characterID int, inventoryType int) ([]*Item, error) {
	var items []*Item

	err := db.Model(&items).
		Where("character_id = ?", characterID).
		Where("inventory_type = ?", inventoryType).
		Relation("Attributes").
		Select()

	return items, err
}

func freeSlot(characterID uint64, inventoryType uint64) (int, error) {
	var (
		slots []int
		slot  int
	)

	err := db.Model((*Item)(nil)).
		Column("slot").
		Where("character_id = ?", characterID).
		Where("inventory_type = ?", inventoryType).
		Select(&slots)

	if err != nil {
		return slot, err
	}

	// for each slot
	// if nextSlot is last one
	//		if nextSlot+1 <= BagInventoryMax
	//		availableSlot = nextSlot+1
	// if nextSlot > currentSlot+1
	//		availableSlot = currentSlot+1
	for i, s := range slots {
		if i+1 == len(slots) {
			if s+1 < BagInventoryMax {
				slot = s + 1
				break
			}
			return slot, errors.Err{
				Code: errors.PersistenceErrInventoryFull,
			}
		}
		if slots[i+1] > s+1 {
			slot = s + 1
			break
		}
	}

	if len(slots) == 0 {
		slot = BagInventoryMin
	}

	return slot, nil
}

func validateItem(item *Item) error {
	if item.Amount == 0 {
		return errors.Err{
			Code: errors.PersistenceErrItemInvalidAmount,
			Details: errors.ErrDetails{
				"stackable": item.Stackable,
				"amount":    item.Amount,
			},
		}
	}

	if item.ShnID == 0 {
		return errors.Err{
			Code: errors.PersistenceErrItemInvalidShnId,
			Details: errors.ErrDetails{
				"shineID": item.ShnID,
			},
		}
	}

	if item.CharacterID == 0 {
		return errors.Err{
			Code: errors.PersistenceErrItemInvalidCharacterId,
			Details: errors.ErrDetails{
				"characterID": item.CharacterID,
			},
		}
	}

	if !item.Stackable {
		if item.Amount > 1 {
			return errors.Err{
				Code: errors.PersistenceErrItemInvalidAmount,
				Details: errors.ErrDetails{
					"stackable": item.Stackable,
					"amount":    item.Amount,
				},
			}
		}
	}

	var cId uint64 = 0
	err := db.Model((*Character)(nil)).Column("id").Where("id = ?", item.CharacterID).Select(&cId)

	if err != nil {
		return errors.Err{
			Code:    errors.PersistenceErrCharNotExists,
			Message: "could not fetch character with id",
			Details: errors.ErrDetails{
				"err":          err,
				"character_id": item.CharacterID,
			},
		}
	}

	return nil
}
