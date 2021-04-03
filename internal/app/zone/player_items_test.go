package zone

import (
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"sync"
	"testing"
)



func TestLoadPlayerInventory_BagInventory(t *testing.T) {
	// p.loadInventory(BagInventory)
}

//
func TestNewItem_Success(t *testing.T) {
	_ = persistence.NewCharacter("mage")

	player := &player{
		baseEntity: baseEntity{
			handle: 1,
		},
		//char: char,
	}


	err := player.load("mage")

	if err != nil {
		t.Fatal(err)
	}

	// item is not persisted here, only in memory
	item, err := makeItem("ShortStaff")

	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(item, player)
	//// item is persisted here
	err = player.newItem(item)

	if err != nil {
		t.Fatal(err)
	}

	// check if item is in player inventory
	item1, ok := player.inventories.inventory.items[1]
	if !ok {
		t.Fail()
	}

	if item1.itemData.itemInfo.InxName != "ShortStaff" {
		t.Fail()
	}

}

func makeItem(itemIndex string) (*item, error) {
	var i = &item{}

	itemData := getItemData(itemIndex)

	if itemData.itemInfo == nil {
		return i, errors.Err{
			Code:    errors.ZoneItemMissingData,
			Details: errors.ErrDetails{
				"itemIndex": itemIndex,
				"type": "ItemInfo",
			},
		}
	}

	if itemData.itemInfoServer == nil {
		return i, errors.Err{
			Code:    errors.ZoneItemMissingData,
			Details: errors.ErrDetails{
				"itemIndex": itemIndex,
				"type": "ItemInfoServer",
			},
		}
	}
	i.itemData = itemData
	i.pItem = &persistence.Item{}

	// first check if there are any random stats using (RandomOption / RandomOptionCount)
	// apply those first, after that check GradeItemOption for fixed stats
	i.stats = itemStats{}

	if itemData.itemInfo.MaxLot > 1 {
		i.stackable = true
	}

	// will vary when created through the ItemDropTables
	// will vary when created through admin command with quantity parameter
	// will not vary if stackable is false
	i.amount = int(itemData.itemInfo.MaxLot)

	return i, nil
}

func getItemData(itemIndex string) *itemData {
	var (
		id = &itemData{}
		wg = &sync.WaitGroup{}
	)

	wg.Add(3)

	go addItemInfoRow(itemIndex, id, wg)

	go addItemInfoServerRow(itemIndex, id, wg)

	go addGradeItemOptionRow(itemIndex, id, wg)

	wg.Wait()

	if id.itemInfoServer.RandomOptionDropGroup != "" {
		wg.Add(2)
		go addRandomOptionRow(id.itemInfoServer.RandomOptionDropGroup, id, wg)
		go addRandomOptionCountRow(id.itemInfoServer.RandomOptionDropGroup, id, wg)
	}

	wg.Wait()

	return id
}

func addItemInfoRow(itemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.ItemInfo.ShineRow {
		if row.InxName == itemIndex {
			id.Lock()
			id.itemInfo = &itemsData.ItemInfo.ShineRow[i]
			id.Unlock()
			return
		}
	}
}

func addItemInfoServerRow(itemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.ItemInfoServer.ShineRow {
		if row.InxName == itemIndex {
			id.Lock()
			id.itemInfoServer = &itemsData.ItemInfoServer.ShineRow[i]
			id.Unlock()
			return
		}
	}
}

func addGradeItemOptionRow(itemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.GradeItemOptions.ShineRow {
		if row.ItemIndex == itemIndex {
			id.Lock()
			id.gradeItemOption = &itemsData.GradeItemOptions.ShineRow[i]
			id.Unlock()
			return
		}
	}
}

func addRandomOptionRow(dropItemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.RandomOption.ShineRow {
		if row.DropItemIndex == dropItemIndex {
			id.Lock()
			id.randomOption = &itemsData.RandomOption.ShineRow[i]
			id.Unlock()
			return
		}
	}
}

func addRandomOptionCountRow(dropItemIndex string, id *itemData, wg *sync.WaitGroup) {
	defer wg.Done()
	for i, row := range itemsData.RandomOptionCount.ShineRow {
		if row.DropItemIndex == dropItemIndex {
			id.Lock()
			id.randomOptionCount = &itemsData.RandomOptionCount.ShineRow[i]
			id.Unlock()
			return
		}
	}
}

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
