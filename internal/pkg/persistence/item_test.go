package persistence

import (
	"errors"
	"github.com/go-pg/pg/v10"
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

// 	// box 8 = equipped items  / 1-29
//	// box 9 = inventory, storage  // 9216 - 9377 (24 slots per page)
//	// box 12 = mini houses // 12288 equipped minihouse, 12299-12322 available slots
func TestCreateItem(t *testing.T) {
	cleanDB()
	// create character
	// create item for that character
	// check that it was stored correctly

	newCharacter("mage")

	item, err := NewItem(db, ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Amount:      1,
		Stackable:   false,
	})

	if err != nil {
		t.Error(err)
	}

	if item.Slot != 9216 {
		t.Errorf("slot = %v, expected slot = %v", item.Slot, 9216)
	}

}

func TestCreateItem_MissingValues(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	// missing amount
	_, err := NewItem(db, ItemParams{
		CharacterID: 1,
		ShnID:       1,
		//Amount:      1,
		Stackable:    false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing shn_id
	_, err = NewItem(db, ItemParams{
		CharacterID: 1,
		//ShnID:       1,
		Amount:      1,
		Stackable:   false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing character_id
	_, err = NewItem(db, ItemParams{
		//CharacterID: 1,
		ShnID:       1,
		Amount:      1,
		Stackable:   false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing stackable
	_, err = NewItem(db, ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Amount:      1,
		//Stackable:   false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing stackable
	_, err = NewItem(db, ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Amount:      1,
		Stackable:   false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}
}

func TestCreateItem_BadPKeys(t *testing.T) {
	cleanDB()
	// du
	newCharacter("mage")

	item0 := &Item{
		//UUID:           uuid.New().String(),
		CharacterID:   1,
		InventoryType: BagInventory,
		Slot:          9216,
		ShnID:         1,
		Stackable:     false,
		Amount:        1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err := db.Model(item0).Insert()

	if err != nil {
		t.Error(err)
	}

	// second try should fail
	item1 := &Item{
		//UUID:           uuid.New().String(),
		CharacterID:   1,
		InventoryType: BagInventory,
		Slot:          9216,
		ShnID:         1,
		Stackable:     false,
		Amount:        1,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	_, err = db.Model(item1).Insert()

	if err == nil {
		t.Error("expected error, got nil")
	}

}

func TestUpdateItem(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item, err := NewItem(db, ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Amount:      1,
		Stackable:   false,
		//ItemData:    &ItemData{},
	})

	if err != nil {
		t.Error(err)
	}

	_, err = UpdateItem(db, item, ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Amount:      1,
		Stackable:   false,
	})

	if err != nil {
		t.Error(err)
	}
}

func UpdateItem(db *pg.DB, item *Item, params ItemParams) (*Item, error) {
	var uItem Item
	return &uItem, errors.New("unimplemented")
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
