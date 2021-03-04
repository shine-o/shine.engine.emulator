package zone

import "testing"

func TestLoadPlayerInventory(t *testing.T) {

}

func TestNewItem_Ok(t *testing.T) {

}

func TestNewItem_BadItemID(t *testing.T) {

}

func TestNewItem_BadItemIndex(t *testing.T) {

}

func TestNewItemStack_Ok(t *testing.T) {

}

func TestNewItemStack_ItemNotStackable(t *testing.T) {

}

func TestSplitItemStack_Ok(t *testing.T) {

}

func TestSplitItemStack_NC_Ok(t *testing.T) {

}

func TestSplitItemStack_BadDivision(t *testing.T) {

}

func TestSplitItemStack_ItemNotStackable(t *testing.T) {

}

func TestItemUnEquip(t *testing.T) {

}

func TestSoftDeleteItem_Ok(t *testing.T) {

}

func TestLoadNewPlayer_BagInventory(t *testing.T) {

}

func TestLoadNewPlayer_BankInventory(t *testing.T) {

}

func TestLoadNewPlayer_RewardInventory(t *testing.T) {

}

func TestLoadNewPlayer_PremiumInventory(t *testing.T) {

}

func TestLoadNewPlayer_MiniHouseInventory(t *testing.T) {
	// should always have the default house

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

func TestItemEquip_Ok(t *testing.T) {

}

func TestItemEquip_NC_Ok(t *testing.T) {

}

func TestItemEquip_BadSlot(t *testing.T) {

}

func TestItemUnEquip_NC_Ok(t *testing.T) {

}

func TestItemUnEquip_Ok(t *testing.T) {

}

func TestChangeItemSlot_Ok(t *testing.T) {

}

func TestChangeItemSlot_NC_Ok(t *testing.T) {

}

func TestChangeItem_NonExistentSlot(t *testing.T) {

}

func TestChangeItemSlot_BadItemType(t *testing.T) {

}

func TestChangeItemSlot_NoItemInSlot(t *testing.T) {

}

func TestDropItem_NonExistingItem(t *testing.T) {

}

func TestSellItem_OK(t *testing.T) {

}

func TestSellItem_NonExistingItem(t *testing.T) {

}

func TestBuyItem_OK(t *testing.T) {

}

func TestOneUseItem_OK(t *testing.T) {

}

// Like mounts, quest items
func TestMultipleUseItem_OK(t *testing.T) {

}
