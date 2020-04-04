package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_CMD
//{
//  unsigned __int16 handle;
//  ABSTATE_INFORMATION info;
//};
type NcBriefInfoAbstateChangeCmd struct {
	Handle uint16
	Info   AbstateInformation
}

func (nc *NcBriefInfoAbstateChangeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoAbstateChangeCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_CMD
	{
	  unsigned __int16 handle;
	  ABSTATE_INFORMATION info;
	};
`
}

//struct PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD
//{
//  unsigned __int16 hnd;
//};
type NcBriefInfoDeleteCmd struct {
	Handle uint16
}

func (nc *NcBriefInfoDeleteCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoDeleteCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_BRIEFINFODELETE_CMD
	{
	  unsigned __int16 hnd;
	};
`
}

//struct PROTO_NC_BRIEFINFO_DROPEDITEM_CMD
//{
//  unsigned __int16 handle;
//  unsigned __int16 itemid;
//  SHINE_XY_TYPE location;
//  unsigned __int16 dropmobhandle;
//  PROTO_NC_BRIEFINFO_DROPEDITEM_CMD::<unnamed-type-attr> attr;
//};
type NcBriefInfoDroppedItemCmd struct {
	Handle        uint16
	ItemID        uint16
	Location      ShineXYType
	DropMobHandle uint16
	Attr          NcBriefInfoDroppedItemCmdAttr
}

func (nc *NcBriefInfoDroppedItemCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoDroppedItemCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_DROPEDITEM_CMD
	{
	  unsigned __int16 handle;
	  unsigned __int16 itemid;
	  SHINE_XY_TYPE location;
	  unsigned __int16 dropmobhandle;
	  PROTO_NC_BRIEFINFO_DROPEDITEM_CMD::<unnamed-type-attr> attr;
	};
`
}

//struct PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD
//{
//  unsigned __int16 handle;
//  unsigned __int16 item;
//  char nSlotNum;
//};
type NcBriefInfoChangeDecorateCmd struct {
	Handle  uint16
	Item    uint16
	SlotNum byte
}

func (nc *NcBriefInfoChangeDecorateCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoChangeDecorateCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_CHANGEDECORATE_CMD
	{
	  unsigned __int16 handle;
	  unsigned __int16 item;
	  char nSlotNum;
	};
`
}

//struct PROTO_NC_BRIEFINFO_MOB_CMD
//{
//  char mobnum;
//  PROTO_NC_BRIEFINFO_REGENMOB_CMD mobs[];
//};
type NcBriefInfoMobCmd struct {
	MobNum byte
	Mobs []NcBriefInfoRegenMobCmd `struct:"sizefrom=MobNum"`
}

func (nc *NcBriefInfoMobCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoMobCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_MOB_CMD
	{
	  char mobnum;
	  PROTO_NC_BRIEFINFO_REGENMOB_CMD mobs[];
	};
`
}