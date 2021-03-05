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

// todo: move to inventory package
type Item struct {
	tableName struct{} `pg:"world.character_items"`
	ID         				 uint64
	//UUID          string `pg:",pk,notnull,type:uuid,default:gen_random_uuid"`
	CharacterID   uint64 `pg:",pk,notnull,unique:item" `
	InventoryType int    `pg:",pk,notnull,unique:item"`
	// box 2 = reward  inventory
	// box 3 = mini house furniture
	// box 8 = equipped items  / 1-29
	// box 9 = inventory, storage  // 9216 - 9377 (24 slots per page)
	// box 12 = mini houses // 12288 equipped minihouse, 12299-12322 available slots
	// box 13 = mini house accessories
	// box 14 = mini house tile all inventory
	// box 15 = premium actions inventory(dances)
	// box 16 = mini house mini game inventory
	Slot      uint16     `pg:",pk,notnull,unique:item"`
	Character *Character `pg:"rel:belongs-to"`

	ShnID     uint16 `pg:",notnull"`
	Stackable bool   `pg:",notnull,use_zero"`
	Amount    uint32
	Attributes *ItemAttributes
	Enchantments *ItemEnchantments
	Licences *ItemLicences
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `pg:",soft_delete"`
}

type ItemAttributes struct {
	tableName struct{} `pg:"world.item_attributes"`
	ID         				 uint64
	ItemID uint64
	Item *Item `pg:"rel:belongs-to"`
}

type ItemEnchantments struct {
	tableName struct{} `pg:"world.item_enchantments"`
	ID         				 uint64
	ItemID uint64
	Item *Item `pg:"rel:belongs-to"`
}

type ItemLicences struct {
	tableName struct{} `pg:"world.item_licences"`
	ID         				 uint64
	ItemID uint64
	Item *Item `pg:"rel:belongs-to"`
	ShnID uint16
}

type ItemParams struct {
	CharacterID uint64
	ShnID       uint16
	Stackable   bool
	Amount      uint32
	Attributes *ItemAttributes
	Enchantments *ItemEnchantments
	Licences * ItemLicences
}

type ErrItem struct {
	Code    int
	Message string
}

func (ei *ErrItem) Error() string {
	return ei.Message
}

// ErrNameTaken name is reserved or in use
var ErrInventoryFull = &ErrItem{
	Code:    1,
	Message: "inventory full",
}

// ErrNameTaken name is reserved or in use
var ErrInvalidAmount = &ErrItem{
	Code:    2,
	Message: "invalid amount specified",
}

func NewItem(db *pg.DB, params ItemParams) (*Item, error) {
	var (
		item *Item
	)

	if params.Amount == 0 {
		return item, ErrInvalidAmount
	}

	slot, err := freeSlot(db, params.CharacterID)

	if err != nil {
		return item, err
	}

	item = &Item{
		//UUID:           uuid.New().String(),
		CharacterID:   params.CharacterID,
		InventoryType: BagInventory,
		Slot:          slot,
		ShnID:         params.ShnID,
		Stackable:     params.Stackable,
		Amount:        params.Amount,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = db.Model(item).Insert()

	return item, err

}

func getItemWhere(db *pg.DB, clauses map[string]interface{}, deleted bool) (*Item, error) {
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

func freeSlot(db * pg.DB, characterID uint64) (uint16, error) {
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
			return slot, ErrInventoryFull
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
