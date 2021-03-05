package persistence

import (
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

	item2, err := NewItem(db, ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
		Attributes: &ItemAttributes{
			Strength: 15,
		},
	})

	if err != nil {
		t.Error(err)
	}

	if item2.Slot != 9217 {
		t.Errorf("slot = %v, expected slot = %v", item2.Slot, 9217)
	}
}

func TestCreateItem_Relations(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item, err := NewItem(db, ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
		Attributes: &ItemAttributes{
			Strength: 15,
		},
	})

	if err != nil {
		t.Error(err)
	}

	if item.Attributes.ID == 0 {
		t.Error("id should not be 0")
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
		Stackable: false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing shn_id
	_, err = NewItem(db, ItemParams{
		CharacterID: 1,
		//ShnID:       1,
		Amount:    1,
		Stackable: false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing character_id
	_, err = NewItem(db, ItemParams{
		//CharacterID: 1,
		ShnID:     1,
		Amount:    1,
		Stackable: false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}
}

func TestCreateItem_BadPKeys(t *testing.T) {
	cleanDB()
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
	})

	if err != nil {
		t.Error(err)
	}

	uItem, err := UpdateItem(db, *item, ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
		Attributes:  &ItemAttributes{
			Strength:  15,
		},
	})

	if err != nil {
		t.Error(err)
	}

	if uItem.Attributes.Strength != 15 {
		t.Errorf("expected value %v, got %v", 15, uItem.Attributes.Strength)
	}
}

// UpdateItem attributes, etc...
// should not handle inventory location
func UpdateItem(db *pg.DB, item Item, params ItemParams) (*Item, error) {

	var uItem = item

	tx, err := db.Begin()

	if err != nil {
		return &item, err
	}

	defer txCloseLog(tx)

	uItem.CharacterID = params.CharacterID

	if uItem.Stackable {
		uItem.Amount = params.Amount
	}

	uItem.Attributes = &ItemAttributes{
		ItemID:    item.ID,
		Item:      &item,
		Strength:  params.Attributes.Strength,
	}

	uItem.UpdatedAt = time.Now()

	_, err = tx.Model(&uItem).
		WherePK().
		Update()

	if err != nil {
		return &uItem, Err{
			Code:    ErrDB,
			Details: ErrDetails{
				"err": err,
				"txErr": tx.Rollback(),
			},
		}
	}

	// if not exists, update
	if item.Attributes == nil {
		_, err = tx.Model(uItem.Attributes).
			Insert()
	} else {
		_, err = tx.Model(uItem.Attributes).
			WherePK().
			Update()
	}


	if err != nil {
		return &uItem, Err{
			Code:    ErrDB,
			Details: ErrDetails{
				"err": err,
				"txErr": tx.Rollback(),
			},
		}
	}

	err = tx.Model(uItem).
		WherePK().
		Relation("Attributes").
		Select()

	if err != nil {
		return &uItem, Err{
			Code:    ErrDB,
			Details: ErrDetails{
				"err": err,
				"txErr": tx.Rollback(),
			},
		}
	}

	return &uItem, tx.Commit()
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
