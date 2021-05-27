package persistence

import (
	"testing"
	"time"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
)

func TestGetCharacterItems(t *testing.T) {
	//db.AddQueryHook(pgdebug.DebugHook{
	//	// Print all queries.
	//	Verbose: true,
	//})

	cleanDB()

	newCharacter("mage")
	for i := BagInventoryMin; i < BagInventoryMax; i++ {
		item := &Item{
			CharacterID:   1,
			Stackable:     false,
			Amount:        1,
			ShnID:         2,
			ShnInxName:    "ShortStaff",
			InventoryType: int(BagInventory),
		}

		err := item.Insert()
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

	// t.Fail()
}

func TestCreateItemOk(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err := item.Insert()
	if err != nil {
		t.Fatal(err)
	}

	if item.Slot != 0 {
		t.Fatalf("slot = %v, expected slot = %v", item.Slot, 0)
	}

	item2 := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err = item2.Insert()

	if err != nil {
		t.Error(err)
	}

	if item2.Slot != 1 {
		t.Fatalf("slot = %v, expected slot = %v", item2.Slot, 1)
	}
}

func TestCreateItemWithAttributes(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
		Attributes: &ItemAttributes{
			StrengthBase: 15,
		},
	}

	err := item.Insert()
	if err != nil {
		t.Error(err)
	}

	if item.Attributes.ID == 0 {
		t.Error("id should not be 0")
	}
}

func TestCreateItemMissingValues(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	// missing amount
	item := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		InventoryType: int(BagInventory),
		Attributes: &ItemAttributes{
			StrengthBase: 15,
		},
	}

	err := item.Insert()

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing shn_id
	item = &Item{
		CharacterID:   1,
		Stackable:     false,
		InventoryType: int(BagInventory),
		Attributes: &ItemAttributes{
			StrengthBase: 15,
		},
	}

	err = item.Insert()

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing shn_inx_name
	item = &Item{
		CharacterID:   1,
		ShnID:         1,
		Stackable:     false,
		InventoryType: int(BagInventory),
		Attributes: &ItemAttributes{
			StrengthBase: 15,
		},
	}

	err = item.Insert()

	if err == nil {
		t.Error("expected error, got none")
	}

	// missing character_id
	item = &Item{
		CharacterID:   0,
		Stackable:     false,
		InventoryType: int(BagInventory),
		Attributes: &ItemAttributes{
			StrengthBase: 15,
		},
	}

	err = item.Insert()

	if err == nil {
		t.Error("expected error, got none")
	}
}

func TestCreateItemCharacterNotExist(t *testing.T) {
	cleanDB()
	// newCharacter("mage")

	item := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
		Attributes: &ItemAttributes{
			StrengthBase: 15,
		},
	}

	err := item.Insert()

	if err == nil {
		t.Error("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Error("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceCharNotExists {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceCharNotExists, cErr.Code)
	}
}

func TestCreateItemBadAmount(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        5,
		InventoryType: int(BagInventory),
		Attributes: &ItemAttributes{
			StrengthBase: 15,
		},
	}

	err := item.Insert()

	if err == nil {
		t.Error("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Error("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceItemInvalidAmount {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceItemInvalidAmount, cErr.Code)
	}

	// 0 amount
	item = &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        0,
		InventoryType: int(BagInventory),
		Attributes: &ItemAttributes{
			StrengthBase: 15,
		},
	}

	err = item.Insert()

	if err == nil {
		t.Error("expected error, got nil")
	}

	cErr, ok = err.(errors.Err)

	if !ok {
		t.Error("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceItemInvalidAmount {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceItemInvalidAmount, cErr.Code)
	}
}

func TestCreateItemBadPKeys(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item0 := &Item{
		CharacterID:   1,
		InventoryType: int(BagInventory),
		Slot:          1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
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
		InventoryType: int(BagInventory),
		Slot:          1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
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

func TestUpdateItemOk(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err := item.Insert()
	if err != nil {
		t.Error(err)
	}

	item.Stackable = true
	item.Amount = 2
	err = item.Update()
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateItemBadAmount(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err := item.Insert()
	if err != nil {
		t.Fatal(err)
	}

	item.Amount = 5

	err = item.Update()

	if err == nil {
		t.Fatal("expected error, got none")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Fatal("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceItemInvalidAmount {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceItemInvalidAmount, cErr.Code)
	}

	// zero amount
	item.Amount = 0

	err = item.Update()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok = err.(errors.Err)

	if !ok {
		t.Fatal("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceItemInvalidAmount {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceItemInvalidAmount, cErr.Code)
	}
}

// not sure if needed
//func TestUpdateItem_DistinctShnID(t *testing.T) {
//	cleanDB()
//	newCharacter("mage")
//
//	item := &Item{
//		CharacterID:   1,
//		ShnID:         1,
//		ShnInxName: "ShortStaff",
//		Stackable:     false,
//		Amount:        1,
//		InventoryType: int(BagInventory),
//	}
//
//	err := item.Insert()
//
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	item.ShnID = 2
//	err = item.UpdateCharacter()
//
//	if err == nil {
//		t.Fatal("expected error, got nil")
//	}
//
//	cErr, ok := err.(errors.Err)
//
//	if !ok {
//		t.Fatal("expected custom error type Err")
//	}
//
//	if cErr.Code != errors.PersistenceItemDistinctShnID {
//		t.Fatalf("expected error code %v, got %v", errors.PersistenceItemDistinctShnID, cErr.Code)
//	}
//}

func TestItemSlotMoveToUnusedSlot(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	item0 := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err := item0.Insert()
	if err != nil {
		t.Fatal(err)
	}

	uItem0, err := item0.MoveTo(BagInventory, 1)
	if err != nil {
		t.Fatal(err)
	}

	if uItem0 != nil {
		t.Fatal("updated item should be nil")
	}

	if item0.Slot != 1 {
		t.Fatalf("expected slot %v, got %v", 1, item0.Slot)
	}
}

func TestItemSlotMoveToUsedSlot(t *testing.T) {
	cleanDB()
	newCharacter("mage")

	// item 1, slot 0
	item0 := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err := item0.Insert()
	if err != nil {
		t.Fatal(err)
	}

	// item 2, slot 1
	item1 := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err = item1.Insert()
	if err != nil {
		t.Fatal(err)
	}

	uItem1, err := item0.MoveTo(BagInventory, 1)
	if err != nil {
		t.Fatal(err)
	}

	if uItem1 == nil {
		t.Fatal("item should not be nil")
	}

	if uItem1.ID != item1.ID {
		t.Fatalf("expected id %v, got %v", item1.ID, uItem1.ID)
	}
}

func TestSoftDeleteItem(t *testing.T) {
	// new entry should be made in another table for deleted items
	cleanDB()

	newCharacter("mage")

	// item 1, slot 0
	item := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err := item.Insert()
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
	// -1 making up for default items
	for i := BagInventoryMin; i < BagInventoryMax; i++ {
		item := &Item{
			CharacterID:   1,
			ShnID:         1,
			ShnInxName:    "ShortStaff",
			Stackable:     false,
			Amount:        1,
			InventoryType: int(BagInventory),
		}

		err := item.Insert()
		if err != nil {
			t.Fatal(err)
		}

	}

	item := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err := item.Insert()

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	cErr, ok := err.(errors.Err)

	if !ok {
		t.Error("expected custom error type Err")
	}

	if cErr.Code != errors.PersistenceInventoryFull {
		t.Fatalf("expected error code %v, got %v", errors.PersistenceInventoryFull, cErr.Code)
	}
}

func TestInventoryAutomaticSlot(t *testing.T) {
	cleanDB()

	newCharacter("mage")

	// create 3 items
	// delete 2nd item
	// try to create another item
	// the slot should be the one freed up by deleting the 2nd item
	item := &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err := item.Insert()
	if err != nil {
		t.Fatal(err)
	}

	item = &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err = item.Insert()

	if err != nil {
		t.Fatal(err)
	}

	item = &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err = item.Insert()

	if err != nil {
		t.Fatal(err)
	}

	err = DeleteItem(2)

	if err != nil {
		t.Fatal(err)
	}

	item = &Item{
		CharacterID:   1,
		ShnID:         1,
		ShnInxName:    "ShortStaff",
		Stackable:     false,
		Amount:        1,
		InventoryType: int(BagInventory),
	}

	err = item.Insert()

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
