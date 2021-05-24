package structs

// struct PROTO_NC_ITEM_EQUIP_REQ
type NcItemEquipReq struct {
	Slot byte
}

// typedef unsigned __int16 PROTO_NC_ITEM_DROP_ACK;
type NcItemDropAck uint16

// struct PROTO_NC_ITEM_PICK_REQ
type NcItemPickReq struct {
	ItemHandle uint16
}

// struct PROTO_NC_CHAR_ADMIN_LEVEL_INFORM_CMD
type NcCharAdminLevelInformCmd struct {
	AdminLevel byte
}

// struct PROTO_NC_ITEM_DROP_REQ
type NcItemDropReq struct {
	Slot     ItemInventory
	Lot      uint32
	Location ShineXYType
}

// struct PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK
type NcItemChangedInventoryOpenAck struct {
	ErrorCode         uint16
	PartMark          byte
	NumOfChargedItems uint16
	ChargedItems      []ChargedItemInfo `struct:"sizefrom=NumOfChargedItems"`
}

// struct PROTO_NC_ITEM_REWARDINVENOPEN_REQ
type NcItemRewardInventoryOpenReq struct {
	Page uint16
}

// struct PROTO_NC_ITEM_CELLCHANGE_CMD
type NcItemCellChangeCmd struct {
	Exchange ItemInventory
	Location ItemInventory
	Item     ShineItemVar
}

// struct PROTO_NC_ITEM_REWARDINVENOPEN_ACK
type NcItemRewardInventoryOpenAck struct {
	Count byte
	Items []ItemPacketInfo `struct:"sizefrom=Count"`
	Unk   byte             // grrr
}

// struct PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ
type NcITemChargedInventoryOpenReq struct {
	Page uint16
}

// struct PROTO_NC_ITEM_USE_REQ
type NcItemUseReq struct {
	Slot byte
	Type byte
}

// struct PROTO_NC_ITEM_PICK_ACK
type NcItemPickAck struct {
	ItemID     uint16
	Lot        uint32
	Error      uint16
	ItemHandle uint16
}

// struct PROTO_NC_ITEM_UNEQUIP_REQ
type NcItemUnequipReq struct {
	SlotEquip byte
	SlotInven byte
}

// struct PROTO_NC_ITEM_RELOC_REQ
type NcitemRelocateReq struct {
	From ItemInventory
	To   ItemInventory
}

// struct NC_ITEM_RELOC_ACK
type NcItemRelocateAck struct {
	Code uint16
}

// struct PROTO_NC_ITEM_EQUIPCHANGE_CMD
type NcItemEquipChangeCmd struct {
	From      ItemInventory
	EquipSlot byte
	ItemData  ShineItemVar
}

// NcItemEquipAck
type NcItemEquipAck struct {
	Code uint16
}

// PROTO_NC_ITEM_SPLIT_REQ
type NcItemSplitReq struct {
	From ItemInventory
	To   ItemInventory
	Lot  uint32
}

// struct PROTO_NC_MENU_OPENSTORAGE_CMD
type NcMenuOpenStorageCmd struct {
	Cen         uint64
	MaxPage     byte
	CurrentPage byte
	OpenType    byte
	CountItems  byte
	Items       []ProtoItemPacketInformation `struct:"sizefrom=CountItems"`
}

// struct PROTO_CHARGED_ITEM_INFO
type ChargedItemInfo struct {
	ItemOrderNo      uint32
	ItemCode         uint32
	ItemAmount       uint32
	ItemRegisterDate ShineDateTime
}

//struct CChargedItem
//{
//  PROTO_CHARGED_ITEM_INFO m_ChargedItemBF[24];
//  int m_NumOfChargedItem;
//};
//type ChargedItems [24]ChargedItemInfo
type ChargedItems struct {
	// Count int
	Items [24]ChargedItemInfo
}

// NC_ITEM_REWARDINVENOPEN_ACK
type NcItemRewardInvenOpenAck struct {
	Count byte
	Items []ItemPacketInfo `struct:"sizefrom=Count"`
}
