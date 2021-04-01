package persistence

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"testing"
	"time"
)

func TestGetCharacterItems(t *testing.T) {

	//db.AddQueryHook(pgdebug.DebugHook{
	//	// Print all queries.
	//	Verbose: true,
	//})

	cleanDB()

	newCharacter("mage")

	for i := BagInventoryMin; i < BagInventoryMax; i++ {
		_, err := NewItem(ItemParams{
			CharacterID: 1,
			ShnID:       1,
			Stackable:   false,
			Amount:      1,
		})

		if err != nil {
			t.Fatal(err)
		}

	}

	items, err := GetCharacterItems(1, BagInventory)

	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 117 {
		t.Fatalf("expected value %v, got %v", 117, len(items))
	}

	for _, item := range items {
		if item.Attributes == nil {
			t.Fatalf("item attributes should not be nil, item_id=%v", item.ID)
		}
	}

	//t.Fail()
}

func TestCreateItem_Ok(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Amount:      1,
		Stackable:   false,
	})

	if err != nil {
		t.Fatal(err)
	}

	if item.Slot != 0 {
		t.Fatalf("slot = %v, expected slot = %v", item.Slot, 0)
	}

	item2, err := NewItem(ItemParams{
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

	if item2.Slot != 1 {
		t.Fatalf("slot = %v, expected slot = %v", item2.Slot, 1)
	}
}

func TestCreateItem_WithAttributes(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item, err := NewItem(ItemParams{
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
	_, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		//Amount:      1,
		Stackable: false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing shn_id
	_, err = NewItem(ItemParams{
		CharacterID: 1,
		//ShnID:       1,
		Amount:    1,
		Stackable: false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing character_id
	_, err = NewItem(ItemParams{
		//CharacterID: 1,
		ShnID:     1,
		Amount:    1,
		Stackable: false,
	})

	if err == nil {
		t.Error("expected error, got none")
	}
}

func TestCreateItem_CharacterNotExist(t *testing.T) {
	cleanDB()
	//newCharacter("mage")

	_, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
		Attributes: &ItemAttributes{
			Strength: 15,
		},
	})

	if err == nil {
		t.Error("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Error("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceErrCharNotExists {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceErrCharNotExists, cErr.Code)
	}

}

func TestCreateItem_BadAmount(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	_, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      5,
		Attributes: &ItemAttributes{
			Strength: 15,
		},
	})

	if err == nil {
		t.Error("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Error("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceErrItemInvalidAmount {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceErrItemInvalidAmount, cErr.Code)
	}

	// 0 amount
	_, err = NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      0,
		Attributes: &ItemAttributes{
			Strength: 15,
		},
	})

	if err == nil {
		t.Error("expected error, got nil")
	}

	cErr, ok = err.(errors.Err)

	if !ok {
		t.Error("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceErrItemInvalidAmount {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceErrItemInvalidAmount, cErr.Code)
	}

}

func TestCreateItem_BadPKeys(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item0 := &Item{
		CharacterID:   1,
		InventoryType: BagInventory,
		Slot:          0,
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
		Slot:          0,
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

func TestUpdateItem_Ok(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Amount:      1,
		Stackable:   false,
	})

	if err != nil {
		t.Error(err)
	}

	uItem, err := item.Update(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   true,
		Amount:      5,
		Attributes: &ItemAttributes{
			Strength: 15,
		},
	})

	if err != nil {
		t.Error(err)
	}

	if uItem.Attributes.Strength != 15 {
		t.Fatalf("expected value %v, got %v", 15, uItem.Attributes.Strength)
	}

	if uItem.Amount != 5 {
		t.Fatalf("expected value %v, got %v", 5, uItem.Amount)
	}
}

func TestUpdateItem_NoAttributes(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Amount:      1,
		Stackable:   false,
	})

	if err != nil {
		t.Error(err)
	}

	_, err = item.Update(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   true,
		Amount:      5,
	})

	if err != nil {
		t.Error(err)
	}
}

func TestUpdateItem_BadAmount(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = item.Update(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      5,
		Attributes: &ItemAttributes{
			Strength: 15,
		},
	})

	if err == nil {
		t.Fatal("expected error, got none")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatal("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceErrItemInvalidAmount {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceErrItemInvalidAmount, cErr.Code)
	}

	// zero amount
	_, err = item.Update(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      0,
		Attributes: &ItemAttributes{
			Strength: 15,
		},
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok = err.(errors.Err)

	if !ok {
		t.Fatal("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceErrItemInvalidAmount {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceErrItemInvalidAmount, cErr.Code)
	}
}

func TestUpdateItem_DistinctShnID(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = item.Update(ItemParams{
		CharacterID: 1,
		ShnID:       2,
		Stackable:   false,
		Amount:      1,
		Attributes: &ItemAttributes{
			Strength: 15,
		},
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatal("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceErrItemDistinctShnID {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceErrItemDistinctShnID, cErr.Code)
	}
}

func TestItemSlot_MoveToUnusedSlot(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item0, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err != nil {
		t.Fatal(err)
	}

	uItem0, err := item0.MoveTo(BagInventory, 1)
	if err != nil {
		t.Fatal(err)
	}

	if uItem0.ID != item0.ID {
		t.Fatalf("expected id %v, got %v", item0.ID, uItem0.ID)
	}
}

func TestItemSlot_MoveToUsedSlot(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	// item 1, slot 0
	item0, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err != nil {
		t.Fatal(err)
	}

	// item 2, slot 1
	item1, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = item0.MoveTo(BagInventory, 1)
	if err != nil {
		t.Fatal(err)
	}

	var uItem0 Item
	err = db.Model(&uItem0).
		Where("character_id = ?", 1).
		Where("inventory_type = ?", BagInventory).
		Where("slot = ?", 1).
		Select()
	if err != nil {
		t.Fatal(err)
	}

	// items should have switched places
	if uItem0.ID != item0.ID {
		t.Fatalf("expected id %v, got %v", item0.ID, uItem0.ID)
	}

	var uItem1 Item
	err = db.Model(&uItem1).
		Where("character_id = ?", 1).
		Where("inventory_type = ?", BagInventory).
		Where("slot = ?", 0).
		Select()

	if err != nil {
		t.Fatal(err)
	}

	if uItem1.ID != item1.ID {
		t.Fatalf("expected id %v, got %v", item0.ID, uItem0.ID)
	}
}

func TestSoftDeleteItem(t *testing.T) {
	// new entry should be made in another table for deleted items
	cleanDB()

	newCharacter("mage")

	// item 1, slot 0
	_, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err != nil {
		t.Fatal(err)
	}

	err = DeleteItem(1)

	if err != nil {
		t.Fatal(err)
	}

	var uItem1 Item

	err = db.Model(&uItem1).
		Where("id = ?", 1).
		Select()

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestInventoryFull(t *testing.T) {
	// new entry should be made in another table for deleted items
	cleanDB()

	newCharacter("mage")

	for i := BagInventoryMin; i < BagInventoryMax; i++ {
		_, err := NewItem(ItemParams{
			CharacterID: 1,
			ShnID:       1,
			Stackable:   false,
			Amount:      1,
		})

		if err != nil {
			t.Fatal("expected error, got nil")
		}

	}

	_, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Error("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceErrInventoryFull {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceErrInventoryFull, cErr.Code)
	}

}

func TestInventory_AutomaticSlot(t *testing.T) {
	cleanDB()

	newCharacter("mage")

	// create 3 items
	// delete 2nd item
	// try to create another item
	// the slot should be the one freed up by deleting the 2nd item
	_, err := NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err != nil {
		t.Fatal(err)
	}

	_, err = NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err != nil {
		t.Fatal(err)
	}

	//

	err = DeleteItem(2)

	if err != nil {
		t.Fatal(err)
	}

	_, err = NewItem(ItemParams{
		CharacterID: 1,
		ShnID:       1,
		Stackable:   false,
		Amount:      1,
	})

	if err != nil {
		t.Fatal(err)
	}

	var uItem1 Item

	err = db.Model(&uItem1).
		Where("id = ?", 4).
		Select()

	if uItem1.Slot != 1 {
		t.Fatalf("expected slot %v, got %v", 1, uItem1.Slot)
	}
}
