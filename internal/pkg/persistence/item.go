package persistence

import (
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
	Stackable   bool       `pg:",notnull,use_zero"`
	Amount      uint32
	Attributes  *ItemAttributes `pg:"rel:belongs-to"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	//DeletedAt time.Time `pg:",soft_delete"`
}

type ItemAttributes struct {
	tableName struct{} `pg:"world.item_attributes"`
	ID        uint64   `pg:",unique:item"`
	ItemID    uint64   `pg:",use_zero,notnull,unique:item"`
	Strength  int
	CreatedAt time.Time
	UpdatedAt time.Time
	//DeletedAt time.Time `pg:",soft_delete"`
}

type ItemParams struct {
	CharacterID uint64
	ShnID       uint16
	Stackable   bool
	Amount      uint32
	Attributes  *ItemAttributes
}

func NewItem(params ItemParams) (*Item, error) {
	var (
		item *Item
		attr *ItemAttributes
	)

	err := validateItemParams(params)

	if err != nil {
		return item, err
	}

	slot, err := freeSlot(params.CharacterID, BagInventory)

	if err != nil {
		return item, err
	}

	tx, err := db.Begin()

	if err != nil {
		return item, err
	}

	defer closeTx(tx)

	item = &Item{
		CharacterID:   params.CharacterID,
		InventoryType: BagInventory,
		Slot:          slot,
		ShnID:         params.ShnID,
		Stackable:     params.Stackable,
		Amount:        params.Amount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = tx.Model(item).
		Returning("*").
		Insert()

	if err != nil {
		return item, err
	}

	if params.Attributes != nil {
		attr = params.Attributes
	} else {
		attr = &ItemAttributes{}
	}

	attr.ItemID = item.ID
	attr.CreatedAt = time.Now()
	attr.UpdatedAt = time.Now()

	_, err = tx.Model(attr).
		Returning("*").
		Insert()

	if err != nil {
		return item, Err{
			Code: ErrDB,
			Details: ErrDetails{
				"err":   err,
				"txErr": tx.Rollback(),
			},
		}
	}

	err = tx.Model(item).
		WherePK().
		Returning("*").
		Relation("Attributes").
		Select()

	if err != nil {
		return item, err
	}

	return item, tx.Commit()

}

func DeleteItem(itemID int) error {
	tx, err := db.Begin()

	if err != nil {
		return err
	}

	defer closeTx(tx)

	_, err = tx.Model((*ItemAttributes)(nil)).Where("item_id = ?", itemID).Delete()

	if err != nil {
		return Err{
			Code: ErrDB,
			Details: ErrDetails{
				"err":    err,
				"txErr":  tx.Rollback(),
				"itemID": itemID,
			},
		}
	}

	_, err = tx.Model((*Item)(nil)).Where("id = ?", itemID).Delete()

	if err != nil {
		return Err{
			Code: ErrDB,
			Details: ErrDetails{
				"err":    err,
				"txErr":  tx.Rollback(),
				"itemID": itemID,
			},
		}
	}

	return tx.Commit()
}

func (i Item) Update(params ItemParams) (*Item, error) {

	err := validateItemParams(params)

	if params.ShnID != i.ShnID {
		return &i, Err{
			Code: ErrItemDistinctShnID,
			Details: ErrDetails{
				"originalShnID": i.ShnID,
				"newShnID":      params.ShnID,
			},
		}
	}

	if err != nil {
		return &i, err
	}

	tx, err := db.Begin()

	if err != nil {
		return &i, err
	}

	defer closeTx(tx)

	i.CharacterID = params.CharacterID

	i.Stackable = params.Stackable

	i.Amount = params.Amount

	i.UpdatedAt = time.Now()

	_, err = tx.Model(&i).
		WherePK().
		Update()

	if err != nil {
		return &i, Err{
			Code: ErrDB,
			Details: ErrDetails{
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

		if params.Attributes != nil {
			i.Attributes.Strength = params.Attributes.Strength
		}

		_, err = tx.Model(i.Attributes).
			WherePK().
			Update()
	}

	if err != nil {
		return &i, Err{
			Code: ErrDB,
			Details: ErrDetails{
				"err":   err,
				"txErr": tx.Rollback(),
			},
		}
	}

	err = tx.Model(&i).
		Returning("*").
		WherePK().
		Relation("Attributes").
		Select()

	if err != nil {
		return &i, Err{
			Code: ErrDB,
			Details: ErrDetails{
				"err":   err,
				"txErr": tx.Rollback(),
			},
		}
	}

	return &i, tx.Commit()
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
			return &i, Err{
				Code: ErrItemSlotUpdate,
				Details: ErrDetails{
					"err":           err,
					"from":      i.Slot,
					"to": otherItem.Slot,
					"fromInventory": i.InventoryType,
					"toInventory": inventoryType,
					"shnID": i.ShnID,
				},
			}
		}

		// there is an item there,
		// switch slot ids
		otherItem.InventoryType = originalInventory
		otherItem.Slot = originalSlot

		_, err = tx.Model(&otherItem).WherePK().Update()
		if err != nil {
			return &i, Err{
				Code: ErrItemSlotUpdate,
				Details: ErrDetails{
					"err":           err,
					"from":      i.Slot,
					"to": otherItem.Slot,
					"fromInventory": i.InventoryType,
					"toInventory": inventoryType,
					"shnID": i.ShnID,
				},
			}
		}
	}

	i.InventoryType = inventoryType
	i.Slot = slot

	_, err = tx.Model(&i).WherePK().Update()

	if err != nil {
		return &i, Err{
			Code: ErrItemSlotUpdate,
			Details: ErrDetails{
				"err":           err,
				"txErr":         tx.Rollback(),
				"from":      i.Slot,
				"to": otherItem.Slot,
				"fromInventory": i.InventoryType,
				"toInventory": inventoryType,
				"shnID": i.ShnID,
			},
		}
	}

	return &i, tx.Commit()
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
			return slot, Err{
				Code: ErrInventoryFull,
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

func validateItemParams(params ItemParams) error {
	if params.Amount == 0 {
		return Err{
			Code: ErrItemInvalidAmount,
			Details: ErrDetails{
				"stackable": params.Stackable,
				"amount":    params.Amount,
			},
		}
	}

	if params.ShnID == 0 {
		return Err{
			Code: ErrItemInvalidShnId,
			Details: ErrDetails{
				"shineID": params.ShnID,
			},
		}
	}

	if params.CharacterID == 0 {
		return Err{
			Code: ErrItemInvalidCharacterId,
			Details: ErrDetails{
				"characterID": params.CharacterID,
			},
		}
	}

	if !params.Stackable {
		if params.Amount > 1 {
			return Err{
				Code: ErrItemInvalidAmount,
				Details: ErrDetails{
					"stackable": params.Stackable,
					"amount":    params.Amount,
				},
			}
		}
	}

	var cId uint64 = 0
	err := db.Model((*Character)(nil)).Column("id").Where("id = ?", params.CharacterID).Select(&cId)

	if err != nil {
		return Err{
			Code:    ErrCharNotExists,
			Message: "could not fetch character with id",
			Details: ErrDetails{
				"err":          err,
				"character_id": params.CharacterID,
			},
		}
	}

	return nil
}
