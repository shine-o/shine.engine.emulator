package persistence

import (
	"github.com/go-pg/pg/v10"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"testing"
	"time"
)

func TestGetCharacterItems(t *testing.T) {
	// create character
	// create dummy items
	// assert items are not nil
	//c := Character{}
	//
	//for _, i := range c.Items {
	//	i.
	//}

}

type NewItemParams struct {
	CharacterID uint64
	ItemData    *ItemData
	ShnID       uint16
	Stackable   bool
	Amount      uint32
}

const (
	EquippedInventory  = 8
	BagInventory       = 9
	MiniHouseInventory = 12

	BagInventoryMin = 9216
	BagInventoryMax = 9377
)

// 	// box 8 = equipped items  / 1-29
//	// box 9 = inventory, storage  // 9216 - 9377 (24 slots per page)
//	// box 12 = mini houses // 12288 equipped minihouse, 12299-12322 available slots
func TestCreateItem(t *testing.T) {
	cleanDB()
	// create character
	// create item for that character
	// check that it was stored correctly

	newCharacter("mage")
	newCharacter("archer")

	item, err := NewItem(db, NewItemParams{
		CharacterID: 1,
		ShnID:       1,
		Amount:      1,
		Stackable:   false,
		ItemData:    &ItemData{},
	})

	if err != nil {
		t.Error(err)
	}

	if item.Slot != 9216 {
		t.Errorf("slot = %v, expected slot = %v", item.Slot, 9216)
	}

	if len(item.Data) == 0 {
		t.Error("item data should not be empty")
	}

	var itemData ItemData
	err = structs.Unpack(item.Data, &itemData)
	if err != nil {
		t.Error(err)
	}
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

func NewItem(db *pg.DB, params NewItemParams) (*Item, error) {

	var (
		item  *Item
		slot  uint16
		slots []uint16
	)

	data, err := structs.Pack(params.ItemData)

	if err != nil {
		return item, err
	}

	err = db.Model((*Item)(nil)).
		Column("slot").
		Where("character_id = ?", params.CharacterID).
		Where("inventory_type = ?", BagInventory).
		Select(&slots)

	if err != nil {
		return item, err
	}

	// for each slot
	// if nextSlot is last one
	//		if nextSlot+1 <= BagInventoryMax
	//		availableSlot = nextSlot+1
	// if nextSlot > currentSlot+1
	//		availableSlot = currentSlot+1
	for i, s := range slots {
		if i+1 == len(slots) {
			if s + 1 <= BagInventoryMax {
				slot = s + 1
				break
			}
			return item, ErrInventoryFull
		}
		if slots[i+1] > s+1 {
			slot = s+ 1
			break
		}
	}

	if len(slots) == 0 {
		slot = BagInventoryMin
	}

	item = &Item{
		//UUID:           uuid.New().String(),
		CharacterID:   params.CharacterID,
		InventoryType: BagInventory,
		Slot:          slot,
		ShnID:         params.ShnID,
		Stackable:     params.Stackable,
		Amount:        params.Amount,
		Data:          data,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = db.Model(item).Insert()

	return item, err

}

func TestCreateItem_BadPKeys(t *testing.T) {
	// du
}

func TestUpdateItem(t *testing.T) {

}

func TestUpdateItem_BadData(t *testing.T) {

}

func TestSoftDeleteItem(t *testing.T) {
	// new entry should be made in another table for deleted items
}

func TestInventoryFull(t *testing.T) {

}

func TestInventory_AutomaticSlot(t *testing.T) {
	// find the first free slot in the inventory
	// select items.Slot with distinct items.Slot
	// var res []struct {
	//    AuthorId  int
	//    BookCount int
	//}
	//err := db.Model((*Book)(nil)).
	//    Column("author_id").
	//    ColumnExpr("count(*) AS book_count").
	//    Group("author_id").
	//    Order("book_count DESC").
	//    Select(&res)
}
