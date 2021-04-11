package structs

//struct PROTO_NC_ITEM_EQUIP_REQ
type NcItemEquipReq struct {
	Slot byte
}

//typedef unsigned __int16 PROTO_NC_ITEM_DROP_ACK;
type NcItemDropAck uint16

//struct PROTO_NC_ITEM_PICK_REQ
type NcItemPickReq struct {
	ItemHandle uint16
}

//struct PROTO_NC_CHAR_ADMIN_LEVEL_INFORM_CMD
type NcCharAdminLevelInformCmd struct {
	AdminLevel byte
}

//struct PROTO_NC_ITEM_DROP_REQ
type NcItemDropReq struct {
	Slot     ItemInventory
	Lot      uint32
	Location ShineXYType
}

//struct PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK
type NcItemChangedInventoryOpenAck struct {
	ErrorCode         uint16
	PartMark          byte
	NumOfChargedItems uint16
	ChargedItems      []ChargedItemInfo `struct:"sizefrom=NumOfChargedItems"`
}

//struct PROTO_NC_ITEM_REWARDINVENOPEN_REQ
type NcItemRewardInventoryOpenReq struct {
	Page uint16
}

//struct PROTO_NC_ITEM_CELLCHANGE_CMD
type NcItemCellChangeCmd struct {
	Exchange ItemInventory
	Location ItemInventory
	Item     ShineItemVar
}

//struct PROTO_NC_ITEM_REWARDINVENOPEN_ACK
type NcItemRewardInventoryOpenAck struct {
	Count byte
	Items []ItemPacketInfo `struct:"sizefrom=Count"`
	Unk   byte             // grrr
}

//struct PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ
type NcITemChargedInventoryOpenReq struct {
	Page uint16
}

//struct PROTO_NC_ITEM_USE_REQ
type NcItemUseReq struct {
	Slot byte
	Type byte
}

//struct PROTO_NC_ITEM_PICK_ACK
type NcItemPickAck struct {
	ItemID     uint16
	Lot        uint32
	Error      uint16
	ItemHandle uint16
}

//struct PROTO_NC_ITEM_UNEQUIP_REQ
type NcItemUnequipReq struct {
	SlotEquip byte
	SlotInven byte
}

//struct PROTO_NC_ITEM_RELOC_REQ
type NcitemRelocateReq struct {
	From ItemInventory
	To   ItemInventory
}
