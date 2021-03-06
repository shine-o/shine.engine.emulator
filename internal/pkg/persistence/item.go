package persistence

import (
	"github.com/go-pg/pg/v10"
	"time"
)

const (
	EquippedInventory  = 8
	BagInventory       = 9
	MiniHouseInventory = 12

	BagInventoryMin = 9216
	BagInventoryMax = 9377
)

type Item struct {
	tableName     struct{} `pg:"world.character_items"`
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
	Slot        uint16     `pg:",notnull,unique:item"`
	CharacterID uint64     `pg:",notnull,unique:item" `
	Character   *Character `pg:"rel:belongs-to"`
	ShnID       uint16     `pg:",notnull"`
	Stackable   bool       `pg:",notnull,use_zero"`
	Amount      uint32
	Attributes *ItemAttributes `pg:"rel:has-one"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `pg:",soft_delete"`
}

type ItemAttributes struct {
	tableName struct{} `pg:"world.item_attributes"`
	ID        uint64
	ItemID    uint64
	Item      *Item
	Strength  int
}

type ItemParams struct {
	CharacterID uint64
	ShnID       uint16
	Stackable   bool
	Amount      uint32
	Attributes  *ItemAttributes
}

func NewItem(db *pg.DB, params ItemParams) (*Item, error) {
	var (
		item *Item
	)

	if params.Amount == 0 {
		return item, Err{
			Code: ErrItemInvalidAmount,
		}
	}

	if params.ShnID == 0 {
		return item, Err{
			Code: ErrItemInvalidShnId,
		}
	}

	if params.CharacterID == 0 {
		return item, Err{
			Code: ErrItemInvalidCharacterId,
		}
	}

	slot, err := freeSlot(db, params.CharacterID)

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
		var attr = *params.Attributes
		attr.ItemID = item.ID
		attr.Item = item

		_, err = tx.Model(&attr).
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

// UpdateItem attributes, etc...
// should not handle inventory location
func UpdateItem(db *pg.DB, item Item, params ItemParams) (*Item, error) {

	tx, err := db.Begin()

	if err != nil {
		return &item, err
	}

	defer closeTx(tx)

	item.CharacterID = params.CharacterID

	if item.Stackable {
		item.Amount = params.Amount
	}

	item.UpdatedAt = time.Now()

	_, err = tx.Model(&item).
		WherePK().
		Update()

	if err != nil {
		return &item, Err{
			Code:    ErrDB,
			Details: ErrDetails{
				"err": err,
				"txErr": tx.Rollback(),
			},
		}
	}

	attr := &ItemAttributes{
		ItemID:    item.ID,
		Item:      &item,
		Strength:  params.Attributes.Strength,
	}

	// if not exists, update
	if item.Attributes == nil {
		_, err = tx.Model(attr).Insert()
	} else {
		_, err = tx.Model(attr).
			WherePK().
			Update()
	}


	if err != nil {
		return &item, Err{
			Code:    ErrDB,
			Details: ErrDetails{
				"err": err,
				"txErr": tx.Rollback(),
			},
		}
	}

	err = tx.Model(&item).
		WherePK().
		Relation("Attributes").
		Select()

	if err != nil {
		return &item, Err{
			Code:    ErrDB,
			Details: ErrDetails{
				"err": err,
				"txErr": tx.Rollback(),
			},
		}
	}

	return &item, tx.Commit()
}

func GetItemWhere(db *pg.DB, clauses map[string]interface{}, deleted bool) (*Item, error) {
	var item Item

	query := db.Model(&item)

	for k, v := range clauses {
		query.Where(k, v)
	}

	if !deleted {
		query.Where("deleted_at = null")
	}

	err := query.Select()

	return &item, err
}

func freeSlot(db *pg.DB, characterID uint64) (uint16, error) {
	var (
		slots []uint16
		slot  uint16
	)

	err := db.Model((*Item)(nil)).
		Column("slot").
		Where("character_id = ?", characterID).
		Where("inventory_type = ?", BagInventory).
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
			if s+1 <= BagInventoryMax {
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
