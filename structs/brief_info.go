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
	Mobs   []NcBriefInfoRegenMobCmd `struct:"sizefrom=MobNum"`
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

//struct PROTO_NC_BRIEFINFO_UNEQUIP_CMD
//{
//  unsigned __int16 handle;
//  char slot;
//};
type NcBriefInfoUnequipCmd struct {
	Handle uint16
	Slot   byte
}

func (nc *NcBriefInfoUnequipCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoUnequipCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_UNEQUIP_CMD
	{
	  unsigned __int16 handle;
	  char slot;
	};
`
}

//struct PROTO_NC_BRIEFINFO_REGENMOB_CMD
//{
//  unsigned __int16 handle;
//  char mode;
//  unsigned __int16 mobid;
//  SHINE_COORD_TYPE coord;
//  char flagstate;
//  PROTO_NC_BRIEFINFO_REGENMOB_CMD::<unnamed-type-flag> flag;
//  char sAnimation[32];
//  char nAnimationLevel;
//  char nKQTeamType;
//  char bRegenAni;
//};
type NcBriefInfoRegenMobCmd struct {
	Handle uint16
	Mode   byte
	MobID  uint16
	Coord  ShineCoordType
	// 0,1
	FlagState      byte
	FlagData       BriefInfoRegenMobCmdFlag
	Animation      [32]byte
	AnimationLevel byte
	KQTeamType     byte
	RegenAni       byte
}

func (nc *NcBriefInfoRegenMobCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoRegenMobCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_REGENMOB_CMD
	{
	  unsigned __int16 handle;
	  char mode;
	  unsigned __int16 mobid;
	  SHINE_COORD_TYPE coord;
	  char flagstate;
	  PROTO_NC_BRIEFINFO_REGENMOB_CMD::<unnamed-type-flag> flag;
	  char sAnimation[32];
	  char nAnimationLevel;
	  char nKQTeamType;
	  char bRegenAni;
	};
`
}

//struct PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD
//{
//  unsigned __int16 handle;
//  char count;
//  ABSTATE_INFORMATION infoList[];
//};
type NcBriefInfoAbstateChangeListCmd struct {
	Handle uint16
	Count  byte
	List   []AbstateInformation `struct:"sizefrom=Count"`
}

func (nc *NcBriefInfoAbstateChangeListCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoAbstateChangeListCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD
	{
	  unsigned __int16 handle;
	  char count;
	  ABSTATE_INFORMATION infoList[];
	};
`
}

//struct PROTO_NC_BRIEFINFO_CHARACTER_CMD
//{
//  char charnum;
//  PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD chars[];
//};
type NcBriefInfoCharacterCmd struct {
	Number byte
	Characters []NcBriefInfoLoginCharacterCmd `struct:"sizefrom=Number"`
}

func (nc *NcBriefInfoCharacterCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoCharacterCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD
	{
	  unsigned __int16 handle;
	  char count;
	  ABSTATE_INFORMATION infoList[];
	};
`
}

//struct PROTO_NC_BRIEFINFO_REGENMOVER_CMD
//{
//  unsigned __int16 nHandle;
//  unsigned int nID;
//  unsigned int nHP;
//  SHINE_COORD_TYPE nCoord;
//  ABNORMAL_STATE_BIT AbstateBit;
//  char nGrade;
//  unsigned __int16 nSlotHandle[10];
//};
type NcBriefInfoRegenMoverCmd struct {
	Handle uint16
	ID uint32
	HP uint32
	Coordinates ShineCoordType
	AbstateBit AbstateBit
	Grade byte
	SlotHandle [9]uint16
}

func (nc *NcBriefInfoRegenMoverCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoRegenMoverCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_REGENMOVER_CMD
	{
	  unsigned __int16 nHandle;
	  unsigned int nID;
	  unsigned int nHP;
	  SHINE_COORD_TYPE nCoord;
	  ABNORMAL_STATE_BIT AbstateBit;
	  char nGrade;
	  unsigned __int16 nSlotHandle[10];
	};
`
}

//struct PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD
//{
//  unsigned __int16 handle;
//  unsigned __int16 item;
//  char upgrade;
//  char nSlotNum;
//};
type NcBriefInfoChangeUpgradeCmd struct {
	Handle uint16
	Item uint16
	Upgrade byte
	SlotNum byte
}

func (nc *NcBriefInfoChangeUpgradeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoChangeUpgradeCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD
	{
	  unsigned __int16 handle;
	  unsigned __int16 item;
	  char upgrade;
	  char nSlotNum;
	};
`
}

//struct PROTO_NC_BRIEFINFO_MOVER_CMD
//{
//  char nMoverNum;
//  PROTO_NC_BRIEFINFO_REGENMOVER_CMD Movers[];
//};
type NcBriefInfoMoverCmd struct {
	Count byte
	Movers []NcBriefInfoRegenMoverCmd `struct:"sizefrom=Count"`
}

func (nc *NcBriefInfoMoverCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoMoverCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_MOVER_CMD
	{
	  char nMoverNum;
	  PROTO_NC_BRIEFINFO_REGENMOVER_CMD Movers[];
	};
`
}

//struct PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD
//{
//  unsigned __int16 handle;
//  Name5 charid;
//  SHINE_COORD_TYPE coord;
//  char mode;
//  char chrclass;
//  PROTO_AVATAR_SHAPE_INFO shape;
//  PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD::<unnamed-type-shapedata> shapedata;
//  unsigned __int16 polymorph;
//  STOPEMOTICON_DESCRIPT emoticon;
//  CHARTITLE_BRIEFINFO chartitle;
//  ABNORMAL_STATE_BIT abstatebit;
//  unsigned int myguild;
//  char type;
//  char isGuildAcademyMember;
//  char IsAutoPick;
//  char Level;
//  char sAnimation[32];
//  unsigned __int16 nMoverHnd;
//  char nMoverSlot;
//  char nKQTeamType;
//  char IsUseItemMinimon;
//};
type NcBriefInfoLoginCharacterCmd struct {
	Handle      uint16
	CharID      Name5
	Coordinates ShineCoordType
	Mode        byte
	Class       byte
	Shape       ProtoAvatarShapeInfo
	ShapeData   NcBriefInfoLoginCharacterCmdShapeData
	Polymorph uint16
	Emoticon StopEmoticonDescript
	CharTitle CharTitleBriefInfo
	AbstateBit AbstateBit
	MyGuild uint32
	Type byte
	IsAcademyMember byte
	IsAutoPick byte
	Level byte
	Animation [32]byte
	MoverHandle uint16
	MoverSlot byte
	KQTeamType byte
	UsingMinipet byte
}

func (nc *NcBriefInfoLoginCharacterCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoLoginCharacterCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_MOVER_CMD
	{
	  char nMoverNum;
	  PROTO_NC_BRIEFINFO_REGENMOVER_CMD Movers[];
	};
`
}

//struct PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD
//{
//  PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD upgradeinfo;
//  unsigned __int16 currentmobid;
//  char currentkilllevel;
//};
type NcBriefInfoChangeWeaponCmd struct {
	UpgradeInfo NcBriefInfoChangeUpgradeCmd
	CurrentMobID uint16
	CurrentKillLevel byte
}

func (nc *NcBriefInfoChangeWeaponCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoChangeWeaponCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_CHANGEWEAPON_CMD
	{
	  PROTO_NC_BRIEFINFO_CHANGEUPGRADE_CMD upgradeinfo;
	  unsigned __int16 currentmobid;
	  char currentkilllevel;
	};
`
}

//struct PROTO_NC_BRIEFINFO_INFORM_CMD
//{
//  unsigned __int16 nMyHnd;
//  NETCOMMAND ReceiveNetCommand;
//  unsigned __int16 hnd;
//};
type NcBriefInfoInformCmd struct {
	MyHandle uint16
	ReceiveNetCommand NetCommand
	Handle uint16
}

func (nc *NcBriefInfoInformCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoInformCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_INFORM_CMD
	{
	  unsigned __int16 nMyHnd;
	  NETCOMMAND ReceiveNetCommand;
	  unsigned __int16 hnd;
	};
`
}

