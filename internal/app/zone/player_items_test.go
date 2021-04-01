package zone

import (
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"testing"
)

func newCharacter(class string) *persistence.Character {
	var (
		bitField byte
		name     string
	)

	switch class {
	case "mage":
		bitField = byte(1 | 16<<2 | 1<<7)
		name = fmt.Sprintf("mage%v", 1)
		break
	case "fighter":
		bitField = byte(1 | 1<<2 | 1<<7)
		name = fmt.Sprintf("fighter%v", 1)
		break
	case "archer":
		bitField = byte(1 | 11<<2 | 1<<7)
		name = fmt.Sprintf("archer%v", 1)
		break
	case "cleric":
		bitField = byte(1 | 6<<2 | 1<<7)
		name = fmt.Sprintf("cleric%v", 1)
		break
	}

	c := structs.NcAvatarCreateReq{
		SlotNum: byte(0),
		Name: structs.Name5{
			Name: name,
		},
		Shape: structs.ProtoAvatarShapeInfo{
			BF:        bitField,
			HairType:  6,
			HairColor: 0,
			FaceShape: 0,
		},
	}

	char, err := persistence.New(1, &c)
	if err != nil {
		log.Fatal(err)
	}

	return char
}

func TestLoadPlayerInventory_BagInventory(t *testing.T) {
	// p.loadInventory(BagInventory)
}

//
//func TestNewItem_Success(t *testing.T) {
//	char := newCharacter("mage")
//
//	player := &player{
//		baseEntity: baseEntity{
//			handle:   1,
//		},
//		char: char,
//	}
//
//	item := item{
//		pItem: &persistence.Item{
//			InventoryType: persistence.EquippedInventory,
//			//Slot:          0,
//			CharacterID:   char.ID,
//			ShnID:         1,
//			Stackable:     false,
//			Amount:        1,
//		},
//	}
//
//	player.newItem(item)
//}
//
//func TestNewItem_WithAttributes(t *testing.T) {
//	char := newCharacter("mage")
//
//	player := &player{
//		baseEntity: baseEntity{
//			handle:   1,
//		},
//		char: char,
//	}
//
//	item := item{
//		pItem: &persistence.Item{
//			InventoryType: persistence.EquippedInventory,
//			//Slot:          0,
//			CharacterID:   char.ID,
//			ShnID:         1,
//			Stackable:     false,
//			Amount:        1,
//		},
//	}
//
//	player.newItem(item)
//}

func TestNewItem_BadItemID(t *testing.T) {

}

func TestNewItem_BadItemIndex(t *testing.T) {

}

func TestNewItemStack_Success(t *testing.T) {

}

func TestNewItemStack_ItemNotStackable(t *testing.T) {

}

func TestSplitItemStack_Success(t *testing.T) {

}

func TestSplitItemStack_NC_Success(t *testing.T) {

}

func TestSplitItemStack_BadDivision(t *testing.T) {

}

func TestSplitItemStack_ItemNotStackable(t *testing.T) {

}

func TestSoftDeleteItem_Success(t *testing.T) {

}

func TestLoadNewPlayer_Mage_EquippedItems(t *testing.T) {
	// should have 1 staff
}

func TestLoadNewPlayer_Warrior_EquippedItems(t *testing.T) {

}

func TestLoadNewPlayer_Archer_EquippedItems(t *testing.T) {

}

func TestLoadNewPlayer_Cleric_EquippedItems(t *testing.T) {

}

func TestPlayer_PicksUpItem(t *testing.T) {

}

func TestPlayer_DropsItem(t *testing.T) {

}

func TestPlayer_DeletesItem(t *testing.T) {

}

func TestItemEquip_Success(t *testing.T) {
	//    SUCCESS = 641, // 0x0281
	//    FAILED = 645, // 0x0285
	player := &player{
		baseEntity: baseEntity{
			handle: 1,
		},
	}

	item := &item{}

	itemSlotChange, err := player.equip(item, data.ItemEquipHat)

	if err != nil {
		t.Fatal(err)
	}

	if itemSlotChange.from != 0 {
		t.Fail()
	}

	if itemSlotChange.to != 1 {
		t.Fail()
	}

	equippedItem, ok := player.inventories.equipped.items[int(data.ItemEquipHat)]

	if !ok {
		t.Fail()
	}

	if equippedItem.pItem.ID != item.pItem.ID {
		t.Fail()
	}

	clauses := make(map[string]interface{})

	clauses["item_id"] = item.pItem.ID
	clauses["character_id"] = player.char.ID
	clauses["inventory_type"] = persistence.EquippedInventory

	_, err = persistence.GetItemWhere(clauses, false)

	if err != nil {
		t.Fail()
	}
}

func TestItemEquip_NC_Success(t *testing.T) {

}

func TestItemEquip_Failed(t *testing.T) {
	// p.equipItem(item{}) (itemSlotChange{}, error)
	// err := error.(ErrorCodeZone)
	// err.Code = ItemEquipFailed
	// err.Details["pHandle"]
	//
}

func TestItemEquip_NC_Failed(t *testing.T) {
	//    FAILED = 645, // 0x0285
	// nc := itemEquipFailNc(err) structs.NcItemEquipFailNc ?
	//nc.Code == 645
}

func TestItemEquip_BadSlot(t *testing.T) {

}

func TestItemUnEquip_NC_Success(t *testing.T) {
}

func TestItemUnEquip_Success(t *testing.T) {

}

func TestChangeItemSlot_Success(t *testing.T) {

}

func TestChangeItemSlot_NC_Success(t *testing.T) {

}

func TestChangeItem_NonExistentSlot(t *testing.T) {

}

func TestChangeItemSlot_BadItemType(t *testing.T) {

}

func TestChangeItemSlot_NoItemInSlot(t *testing.T) {

}

func TestDropItem_NonExistingItem(t *testing.T) {

}

func TestSellItem_Success(t *testing.T) {

}

func TestSellItem_NonExistingItem(t *testing.T) {

}

func TestBuyItem_Success(t *testing.T) {

}

func TestOneUseItem_Success(t *testing.T) {

}

// Like mounts, quest items
func TestMultipleUseItem_Success(t *testing.T) {

}
