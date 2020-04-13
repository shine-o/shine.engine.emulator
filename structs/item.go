package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_ITEM_EQUIP_REQ
//{
//  char slot;
//};
type NcItemEquipReq struct {
	Slot byte
}

func (nc *NcItemEquipReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemEquipReq) PdbType() string {
	return `
	struct PROTO_NC_ITEM_EQUIP_REQ
	{
	  char slot;
	};
`
}

//typedef unsigned __int16 PROTO_NC_ITEM_DROP_ACK;
type NcItemDropAck uint16

func (nc *NcItemDropAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemDropAck) PdbType() string {
	return `
	typedef unsigned __int16 PROTO_NC_ITEM_DROP_ACK;
`
}

//struct PROTO_NC_ITEM_PICK_REQ
//{
//  unsigned __int16 itemhandle;
//};
type NcItemPickReq struct {
	ItemHandle uint16
}

func (nc *NcItemPickReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemPickReq) PdbType() string {
	return `
	struct PROTO_NC_ITEM_PICK_REQ
	{
	  unsigned __int16 itemhandle;
	};
`
}

//struct PROTO_NC_CHAR_ADMIN_LEVEL_INFORM_CMD
//{
//  char nAdminLevel;
//};
type NcCharAdminLevelInformCmd struct {
	AdminLevel byte
}

func (nc *NcCharAdminLevelInformCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharAdminLevelInformCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_ADMIN_LEVEL_INFORM_CMD
	{
	  char nAdminLevel;
	};
`
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

func (nc *NcItemDropReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemDropReq) PdbType() string {
	return `
	struct PROTO_NC_ITEM_DROP_REQ
	{
	  ITEM_INVEN slot;
	  unsigned int lot;
	  SHINE_XY_TYPE loc;
	};
`
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

func (nc *NcItemChangedInventoryOpenAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemChangedInventoryOpenAck) PdbType() string {
	return `
	struct PROTO_NC_ITEM_CHARGEDINVENOPEN_ACK
	{
	  unsigned __int16 ErrorCode;
	  char nPartMark;
	  unsigned __int16 NumOfChargedItem;
	  PROTO_CHARGED_ITEM_INFO ChargedItemInfoList[];
	};
`
}

//struct PROTO_NC_ITEM_REWARDINVENOPEN_REQ
//{
//  unsigned __int16 page;
//};
type NcItemRewardInventoryOpenReq struct {
	Page uint16
}

func (nc *NcItemRewardInventoryOpenReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemRewardInventoryOpenReq) PdbType() string {
	return `
	struct PROTO_NC_ITEM_REWARDINVENOPEN_REQ
	{
	  unsigned __int16 page;
	};
`
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

func (nc *NcItemCellChangeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemCellChangeCmd) PdbType() string {
	return `
	struct PROTO_NC_ITEM_CELLCHANGE_CMD
	{
	  ITEM_INVEN exchange;
	  ITEM_INVEN location;
	  SHINE_ITEM_VAR_STRUCT item;
	};
`
}

//struct PROTO_NC_ITEM_REWARDINVENOPEN_ACK
//{
//  char itemcounter;
//  PROTO_ITEMPACKET_INFORM itemarray[];
//};
type NcItemRewardInventoryOpenAck struct {
	Count byte
	Items []ItemPacketInfo `struct:"sizefrom=Count"`
}

func (nc *NcItemRewardInventoryOpenAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemRewardInventoryOpenAck) PdbType() string {
	return `
	struct PROTO_NC_ITEM_REWARDINVENOPEN_ACK
	{
	  char itemcounter;
	  PROTO_ITEMPACKET_INFORM itemarray[];
	};
`
}

//struct PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ
//{
//  unsigned __int16 page;
//};
type NcITemChargedInventoryOpenReq struct {
	Page uint16
}

func (nc *NcITemChargedInventoryOpenReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcITemChargedInventoryOpenReq) PdbType() string {
	return `
	struct PROTO_NC_ITEM_CHARGEDINVENOPEN_REQ
	{
	  unsigned __int16 page;
	};
`
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

func (nc *NcItemUseReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemUseReq) PdbType() string {
	return `
	struct PROTO_NC_ITEM_USE_REQ
	{
	  char invenslot;
	  char invenType;
	};
`
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

func (nc *NcItemPickAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemPickAck) PdbType() string {
	return `
	struct PROTO_NC_ITEM_PICK_ACK
	{
	  unsigned __int16 itemid;
	  unsigned int lot;
	  unsigned __int16 error;
	  unsigned __int16 itemhandle;
	};
`
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

func (nc *NcItemUnequipReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcItemUnequipReq) PdbType() string {
	return `
	struct PROTO_NC_ITEM_UNEQUIP_REQ
	{
	  char slotequip;
	  char slotinven;
	};
`
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

func (nc *NcitemRelocateReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcitemRelocateReq) PdbType() string {
	return `
	struct PROTO_NC_ITEM_RELOC_REQ
	{
	  ITEM_INVEN from;
	  ITEM_INVEN to;
	};
`
}
