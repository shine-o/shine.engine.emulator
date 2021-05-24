package persistence

import (
	"time"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
)

type InventoryType int

const (
	UnknownInventory   InventoryType = 0
	BufferInventory    InventoryType = 100
	EquippedInventory  InventoryType = 8
	BagInventory       InventoryType = 9
	DepositInventory   InventoryType = 6
	RewardInventory    InventoryType = 2
	MiniHouseInventory InventoryType = 12

	BufferInventoryMin = 0
	BufferInventoryMax = 1024

	BagInventoryMin = 0
	BagInventoryMax = 191

	DepositInventoryMin       = 0
	DepositInventoryMax       = 144
	DepositInventoryPageLimit = 36

	RewardInventoryMin = 0
	RewardInventoryMax = 4096

	EquippedInventoryMin = 1
	EquippedInventoryMax = 29

	MiniHouseInventoryMin = 0
	MiniHouseInventoryMax = 23
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
	// DeletedAt time.Time `pg:",soft_delete"`
}

type ItemAttributes struct {
	tableName         struct{} `pg:"world.item_attributes"`
	ID                uint64   `pg:",unique:item"`
	ItemID            uint64   `pg:",use_zero,notnull,unique:item"`
	StrengthBase      int      `pg:",use_zero"`
	StrengthExtra     int      `pg:",use_zero"`
	DexterityBase     int      `pg:",use_zero"`
	DexterityExtra    int      `pg:",use_zero"`
	IntelligenceBase  int      `pg:",use_zero"`
	IntelligenceExtra int      `pg:",use_zero"`
	EnduranceBase     int      `pg:",use_zero"`
	EnduranceExtra    int      `pg:",use_zero"`
	SpiritBase        int      `pg:",use_zero"`
	SpiritExtra       int      `pg:",use_zero"`
	PAttackBase       int      `pg:",use_zero"`
	PAttackExtra      int      `pg:",use_zero"`
	MAttackBase       int      `pg:",use_zero"`
	MAttackExtra      int      `pg:",use_zero"`
	MDefenseBase      int      `pg:",use_zero"`
	MDefenseExtra     int      `pg:",use_zero"`
	PDefenseBase      int      `pg:",use_zero"`
	PDefenseExtra     int      `pg:",use_zero"`
	AimBase           int      `pg:",use_zero"`
	AimExtra          int      `pg:",use_zero"`
	EvasionBase       int      `pg:",use_zero"`
	EvasionExtra      int      `pg:",use_zero"`
	MaxHPBase         int      `pg:",use_zero"`
	MaxHPExtra        int      `pg:",use_zero"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	// DeletedAt time.Time `pg:",soft_delete"`
}

type ItemParams struct {
	CharacterID uint64
	ShnID       uint16
	Stackable   bool
	Amount      int
	Attributes  *ItemAttributes
}

func getInventoryType(val int) InventoryType {
	switch val {
	case 8:
		return EquippedInventory
	case 9:
		return BagInventory
	case 6:
		return DepositInventory
	case 12:
		return MiniHouseInventory
	case 2:
		return RewardInventory
	default:
		return UnknownInventory
	}
}

// Insert an Item
// Items can only be persisted if they are linked to a character
func (i *Item) Insert() error {
	err := validateItem(i)
	if err != nil {
		return err
	}

	if InventoryType(i.InventoryType) != EquippedInventory {
		slot, err := freeSlot(i.CharacterID, InventoryType(i.InventoryType))
		if err != nil {
			return err
		}

		i.Slot = slot
	}

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
			Code: errors.PersistenceItemDistinctShnID,
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

func limitExceeded(inventoryType InventoryType, slot int) bool {
	var max int
	switch inventoryType {
	case BagInventory:
		max = BagInventoryMax
		break
	case EquippedInventory:
		max = EquippedInventoryMax
		break
	case DepositInventory:
		max = DepositInventoryMax
		break
	case MiniHouseInventory:
		max = MiniHouseInventoryMax
		break
	case RewardInventory:
		max = RewardInventoryMax
		break
	}

	if slot > max {
		return true
	}

	return false
}

func (i *Item) MoveTo(inventoryType InventoryType, slot int) (*Item, error) {
	otherItem := &Item{}

	tx, err := db.Begin()
	if err != nil {
		return otherItem, err
	}

	defer closeTx(tx)

	_, err = freeSlot(i.CharacterID, inventoryType)

	if err != nil {
		return otherItem, err
	}

	if limitExceeded(inventoryType, slot) {
		return otherItem, errors.Err{
			Code: errors.PersistenceOutOfRangeSlot,
			Details: errors.ErrDetails{
				"slot":          slot,
				"inventoryType": inventoryType,
			},
		}
	}

	err = tx.Model(otherItem).
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
		i.InventoryType = int(BufferInventory)
		i.Slot = 0
		_, err := tx.Model(i).WherePK().Update()
		if err != nil {
			return otherItem, errors.Err{
				Code: errors.PersistenceItemSlotUpdate,
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

		_, err = tx.Model(otherItem).WherePK().Update()
		if err != nil {
			return otherItem, errors.Err{
				Code: errors.PersistenceItemSlotUpdate,
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
	} else {
		otherItem = nil
	}

	i.InventoryType = int(inventoryType)
	i.Slot = slot

	_, err = tx.Model(i).WherePK().Update()

	if err != nil {
		return otherItem, errors.Err{
			Code: errors.PersistenceItemSlotUpdate,
			Details: errors.ErrDetails{
				"err":           err,
				"txErr":         tx.Rollback(),
				"from":          i.Slot,
				"to":            slot,
				"fromInventory": i.InventoryType,
				"toInventory":   inventoryType,
				"shnID":         i.ShnID,
			},
		}
	}

	return otherItem, tx.Commit()
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

func GetCharacterItems(characterID int, inventoryType InventoryType) ([]*Item, error) {
	var items []*Item

	err := db.Model(&items).
		Where("character_id = ?", characterID).
		Where("inventory_type = ?", inventoryType).
		Relation("Attributes").
		Select()

	return items, err
}

// for each slot
// if nextSlot is last one
//		if nextSlot+1 <= BagInventoryMax
//		availableSlot = nextSlot+1
// if nextSlot > currentSlot+1
//		availableSlot = currentSlot+1
func freeSlot(characterID uint64, inventoryType InventoryType) (int, error) {
	var (
		slots    []int
		slot     int
		min, max int
	)

	err := db.Model((*Item)(nil)).
		Column("slot").
		Where("character_id = ?", characterID).
		Where("inventory_type = ?", inventoryType).
		Select(&slots)
	if err != nil {
		return slot, err
	}

	switch inventoryType {
	case BagInventory:
		min = BagInventoryMin
		max = BagInventoryMax
		break
	case DepositInventory:
		min = DepositInventoryMin
		max = DepositInventoryMax
		break
	case BufferInventory:
		min = BufferInventoryMin
		max = BufferInventoryMax
	case EquippedInventory:
		min = EquippedInventoryMin
		max = EquippedInventoryMax
	case RewardInventory:
		min = RewardInventoryMin
		max = RewardInventoryMax
	}

	for i, s := range slots {
		if i+1 == len(slots) {
			if s+1 < max {
				slot = s + 1
				break
			}
			return slot, errors.Err{
				Code: errors.PersistenceInventoryFull,
			}
		}
		if slots[i+1] > s+1 {
			slot = s + 1
			break
		}
	}

	if len(slots) == 0 {
		slot = min
	}

	return slot, nil
}

func validateItem(item *Item) error {
	inventoryType := getInventoryType(item.InventoryType)

	if inventoryType == UnknownInventory {
		return errors.Err{
			Code: errors.PersistenceUnknownInventory,
			Details: errors.ErrDetails{
				"inventoryType": inventoryType,
			},
		}
	}

	if item.Amount == 0 {
		return errors.Err{
			Code: errors.PersistenceItemInvalidAmount,
			Details: errors.ErrDetails{
				"stackable": item.Stackable,
				"amount":    item.Amount,
			},
		}
	}

	if item.ShnID == 0 {
		return errors.Err{
			Code: errors.PersistenceItemInvalidShnId,
			Details: errors.ErrDetails{
				"shineID": item.ShnID,
			},
		}
	}

	if item.CharacterID == 0 {
		return errors.Err{
			Code: errors.PersistenceItemInvalidCharacterId,
			Details: errors.ErrDetails{
				"characterID": item.CharacterID,
			},
		}
	}

	if !item.Stackable {
		if item.Amount > 1 {
			return errors.Err{
				Code: errors.PersistenceItemInvalidAmount,
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
			Code:    errors.PersistenceCharNotExists,
			Message: "could not fetch character with id",
			Details: errors.ErrDetails{
				"err":          err,
				"character_id": item.CharacterID,
			},
		}
	}

	return nil
}
