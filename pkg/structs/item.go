package structs

//struct PROTO_NC_ITEM_EQUIP_REQ
//{
//  char slot;
//};
type NcItemEquipReq struct {
	Slot byte
}

//typedef unsigned __int16 PROTO_NC_ITEM_DROP_ACK;
type NcItemDropAck uint16

//struct PROTO_NC_ITEM_PICK_REQ
//{
//  unsigned __int16 itemhandle;
//};
type NcItemPickReq struct {
	ItemHandle uint16
}

//struct PROTO_NC_CHAR_ADMIN_LEVEL_INFORM_CMD
//{
//  char nAdminLevel;
//};
type NcCharAdminLevelInformCmd struct {
	AdminLevel byte
}

//struct PROTO_NC_ITEM_DROP_REQ
//{
//  ITEM_INVEN slot;
//  unsigned int lot;
//  SHINE_XY_TYPE loc;
//};
type NcItemDropReq struct {
	Slot     ItemInventory
	Lot      uint32
	Location ShineXYType
}

//struct PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK
//{
//  unsigned __int16 ErrorCode;
//  char nPartMark;
//  unsigned __int16 NumOfChargedItem;
//  PROTO_CHARGED_ITEM_INFO ChargedItemInfoList[];
//};
type NcItemChangedInventoryOpenAck struct {
	ErrorCode         uint16
	PartMark          byte
	NumOfChargedItems uint16
	ChargedItems      []ChargedItemInfo `struct:"sizefrom=NumOfChargedItems"`
}

//struct PROTO_NC_ITEM_REWARDINVENOPEN_REQ
//{
//  unsigned __int16 page;
//};
type NcItemRewardInventoryOpenReq struct {
	Page uint16
}

//struct PROTO_NC_ITEM_CELLCHANGE_CMD
//{
//  ITEM_INVEN exchange;
//  ITEM_INVEN location;
//  SHINE_ITEM_VAR_STRUCT item;
//};
type NcItemCellChangeCmd struct {
	Exchange ItemInventory
	Location ItemInventory
	Item     ShineItemVar
}

//struct PROTO_NC_ITEM_REWARDINVENOPEN_ACK
//{
//  char itemcounter;
//  PROTO_ITEMPACKET_INFORM itemarray[];
//};
type NcItemRewardInventoryOpenAck struct {
	Count byte
	Items []ItemPacketInfo `struct:"sizefrom=Count"`
	Unk   byte             // grrr
}

//struct PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ
//{
//  unsigned __int16 page;
//};
type NcITemChargedInventoryOpenReq struct {
	Page uint16
}

//struct PROTO_NC_ITEM_USE_REQ
//{
//  char invenslot;
//  char invenType;
//};
type NcItemUseReq struct {
	Slot byte
	Type byte
}

//struct PROTO_NC_ITEM_PICK_ACK
//{
//  unsigned __int16 itemid;
//  unsigned int lot;
//  unsigned __int16 error;
//  unsigned __int16 itemhandle;
//};
type NcItemPickAck struct {
	ItemID     uint16
	Lot        uint32
	Error      uint16
	ItemHandle uint16
}

//struct PROTO_NC_ITEM_UNEQUIP_REQ
//{
//  char slotequip;
//  char slotinven;
//};
type NcItemUnequipReq struct {
	SlotEquip byte
	SlotInven byte
}

//struct PROTO_NC_ITEM_RELOC_REQ
//{
//  ITEM_INVEN from;
//  ITEM_INVEN to;
//};
type NcitemRelocateReq struct {
	From ItemInventory
	To   ItemInventory
}
