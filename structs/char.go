package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_CHAR_CLIENT_SKILL_CMD
//{
//	char restempow;
//	PARTMARK PartMark;
//	unsigned __int16 nMaxNum;
//	PROTO_NC_CHAR_SKILLCLIENT_CMD skill;
//};
type NcCharClientSkillCmd struct {
	RestEmpower byte
	PartMark    PartMark
	MaxNum      uint16
	Skills      NcCharSkillClientCmd
}

func (nc *NcCharClientSkillCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharClientSkillCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_CLIENT_SKILL_CMD
	{
	  char restempow;
	  PARTMARK PartMark;
	  unsigned __int16 nMaxNum;
	  PROTO_NC_CHAR_SKILLCLIENT_CMD skill;
	};
`
}

//struct PROTO_NC_CHAR_SKILLCLIENT_CMD
//{
//	unsigned int chrregnum;
//	unsigned __int16 number;
//	PROTO_SKILLREADBLOCKCLIENT skill[];
//};
type NcCharSkillClientCmd struct {
	ChrRegNum uint32
	Number    uint16
	Skills    []SkillReadBlockClient `struct:"sizefrom=Number"`
}

func (nc *NcCharSkillClientCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharSkillClientCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_SKILLCLIENT_CMD
	{
	  unsigned int chrregnum;
	  unsigned __int16 number;
	  PROTO_SKILLREADBLOCKCLIENT skill[];
	};
`
}

//struct PROTO_NC_CHAR_CLIENT_ITEM_CMD
//{
//	char numofitem;
//	char box;
//	PROTO_NC_CHAR_CLIENT_ITEM_CMD::<unnamed-type-flag> flag;
//	PROTO_ITEMPACKET_INFORM ItemArray[];
//};
type NcCharClientItemCmd struct {
	NumOfItem byte `struct:"byte"`
	Box       byte `struct:"byte"`
	Flag      ProtoNcCharClientItemCmdFlag
	Items     []ProtoItemPacketInformation `struct:"sizefrom=NumOfItem"`
}

func (nc *NcCharClientItemCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharClientItemCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_CLIENT_ITEM_CMD
	{
	  char numofitem;
	  char box;
	  PROTO_NC_CHAR_CLIENT_ITEM_CMD::<unnamed-type-flag> flag;
	  PROTO_ITEMPACKET_INFORM ItemArray[];
	};
`
}

//struct PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD
//{
//  char CurrentTitle;
//  char CurrentTitleElement;
//  unsigned __int16 CurrentTitleMobID;
//  unsigned __int16 NumOfTitle;
//  CT_INFO TitleArray[];
//};
type NcClientCharTitleCmd struct {
	CurrentTitle        byte
	CurrentTitleElement byte
	CurrentTitleMobID   uint16
	NumOfTitle          uint16
	Titles              []CharTitleInfo `struct:"sizefrom=NumOfTitle"`
}

func (nc *NcClientCharTitleCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcClientCharTitleCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD
	{
	  char CurrentTitle;
	  char CurrentTitleElement;
	  unsigned __int16 CurrentTitleMobID;
	  unsigned __int16 NumOfTitle;
	  CT_INFO TitleArray[];
	};
`
}

// struct PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ
//{
//  char dummy[1];
//};
type NcCharOptionGetShortcutSizeReq struct {
}

func (nc *NcCharOptionGetShortcutSizeReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharOptionGetShortcutSizeReq) PdbType() string {
	return `
	struct PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ
	{
	  char dummy[1];
	};
`
}

//struct PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK
//{
//  char bSuccess;
//  PROTO_NC_CHAR_OPTION_SHORTCUTSIZE Data;
//};
type NcCharOptionGetShortcutSizeAck struct {
	Success byte
	Data    NcCharOptionShortcutSize
}

func (nc *NcCharOptionGetShortcutSizeAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharOptionGetShortcutSizeAck) PdbType() string {
	return `
	struct PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK
	{
	  char bSuccess;
	  PROTO_NC_CHAR_OPTION_SHORTCUTSIZE Data;
	};
`
}

//struct PROTO_NC_CHAR_GUILD_CMD
//{
//  unsigned int nGuildNo;
//  GUILD_CLIENT Guild[];
//};
type NcCharGuildCmd struct {
	GuildNumber uint32
	Guilds      GuildClient
}

func (nc *NcCharGuildCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharGuildCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_GUILD_CMD
	{
	  unsigned int nGuildNo;
	  GUILD_CLIENT Guild[];
	};
`
}

//struct PROTO_NC_CHAR_LOGIN_REQ
//{
//  char slot;
//};
type NcCharLoginReq struct {
	Slot byte
}

func (nc *NcCharLoginReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharLoginReq) PdbType() string {
	return `
	struct PROTO_NC_CHAR_LOGIN_REQ
	{
	  char slot;
	};
`
}

//struct PROTO_NC_CHAR_MYSTERYVAULT_UI_STATE_CMD
//{
//  char mystery_vault;
//};
type CharMysteryVaultUiStateCmd struct {
	MysteryVault byte
}

func (nc *CharMysteryVaultUiStateCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *CharMysteryVaultUiStateCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_MYSTERYVAULT_UI_STATE_CMD
	{
	  char mystery_vault;
	};
`
}

//struct PROTO_NC_CHAR_USEITEM_MINIMON_INFO_CLIENT_CMD
//{
//  USEITEM_MINIMON_INFO UseItemMinimonInfo;
//};
type CharUseItemMiniMonsterInfoClientCmd struct {
	MiniMonsterInfo UseItemMiniMonsterInfo
}

func (nc *CharUseItemMiniMonsterInfoClientCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *CharUseItemMiniMonsterInfoClientCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_USEITEM_MINIMON_INFO_CLIENT_CMD
	{
	  USEITEM_MINIMON_INFO UseItemMinimonInfo;
	};
`
}

//struct PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ
//{
//  char btMysteryType;
//};
type NcCharUiStateSaveReq struct {
	MysteryType byte
}

func (nc *NcCharUiStateSaveReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharUiStateSaveReq) PdbType() string {
	return `
	struct PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ
	{
	  char btMysteryType;
	};
`
}

//struct PROTO_NC_CHAR_LOGIN_ACK
//{
//  Name4 zoneip;
//  unsigned __int16 zoneport;
//};
type NcCharLoginAck struct {
	ZoneIP   Name4
	ZonePort uint16
}

func (nc *NcCharLoginAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharLoginAck) PdbType() string {
	return `
	struct PROTO_NC_CHAR_LOGIN_ACK
	{
	  Name4 zoneip;
	  unsigned __int16 zoneport;
	};
`
}

//struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD
//{
//  unsigned __int16 nGameOptionDataCnt;
//  GAME_OPTION_DATA GameOptionData[];
//};
type NcCharOptionImproveGetGameOptionCmd struct {
	OptionCount uint16
	GameOptions []GameOptionData `struct:"sizefrom=OptionCount"`
}

func (nc *NcCharOptionImproveGetGameOptionCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharOptionImproveGetGameOptionCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD
	{
	  unsigned __int16 nGameOptionDataCnt;
	  GAME_OPTION_DATA GameOptionData[];
	};
`
}

//struct PROTO_NC_CHAR_GUILD_ACADEMY_CMD
//{
//  unsigned int nGuildAcademyNo;
//  char isGuildAcademyMember;
//  GUILD_ACADEMY_CLIENT GuildAcademy[];
//};
type NcCharGuildAcademyCmd struct {
	GuildAcademyNo uint32
	//GuildAcademyNo uint16
	IsGuildAcademyMember byte
	//GuildAcademyClients []GuildAcademyClient `struct:"sizefrom=GuildAcademyNo"`
	GuildAcademyClient GuildAcademyClient
}

func (nc *NcCharGuildAcademyCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharGuildAcademyCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_GUILD_ACADEMY_CMD
	{
	  unsigned int nGuildAcademyNo;
	  char isGuildAcademyMember;
	  GUILD_ACADEMY_CLIENT GuildAcademy[];
	};
`
}

//struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD
//{
//unsigned __int16 nShortCutDataCnt;
//SHORT_CUT_DATA ShortCutData[];
//};
type NcCharGetShortcutDataCmd struct {
	Count     uint16
	Shortcuts []ShortCutData `struct:"sizefrom=Count"`
}

func (nc *NcCharGetShortcutDataCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharGetShortcutDataCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD
	{
	  unsigned __int16 nShortCutDataCnt;
	  SHORT_CUT_DATA ShortCutData[];
	};
`
}

//struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD
//{
//  unsigned __int16 nKeyMapDataCnt;
//  KEY_MAP_DATA KeyMapData[];
//};
type NcCharGetKeyMapCmd struct {
	Count uint16
	Keys  []KeyMapData `struct:"sizefrom=Count"`
}

func (nc *NcCharGetKeyMapCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharGetKeyMapCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD
	{
	  unsigned __int16 nKeyMapDataCnt;
	  KEY_MAP_DATA KeyMapData[];
	};
`
}

//struct PROTO_NC_CHAR_STAT_REMAINPOINT_CMD
//{
//  char remain;
//};
type NcCharStatRemainPointCmd struct {
	Remain byte
}

func (nc *NcCharStatRemainPointCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharStatRemainPointCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_STAT_REMAINPOINT_CMD
	{
	  char remain;
	};
`
}

//struct PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK
//{
//  char bSuccess;
//  PROTO_NC_CHAR_OPTION_WINDOWPOS Data;
//};
type NcCharOptionGetWindowPosAck struct {
	Success byte
	Data    NcCharOptionWindowPos
}

func (nc *NcCharOptionGetWindowPosAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharOptionGetWindowPosAck) PdbType() string {
	return `
	struct PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK
	{
	  char bSuccess;
	  PROTO_NC_CHAR_OPTION_WINDOWPOS Data;
	};
`
}

//struct PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD
//{
//  char nGuideView;
//};
type NcCharNewbieGuideViewSetCmd struct {
	GuideView byte
}

func (nc *NcCharNewbieGuideViewSetCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharNewbieGuideViewSetCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK
	{
	  char bSuccess;
	  PROTO_NC_CHAR_OPTION_WINDOWPOS Data;
	};
`
}

//struct PROTO_NC_CHAR_CLIENT_AUTO_PICK_CMD
//{
//  unsigned __int16 player;
//  char bEnable;
//};
type NcCharClientAutoPickCmd struct {
	Player uint16
	Enable byte
}

func (nc *NcCharClientAutoPickCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharClientAutoPickCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_CLIENT_AUTO_PICK_CMD
	{
	  unsigned __int16 player;
	  char bEnable;
	};
`
}

//struct PROTO_NC_CHAR_USEITEM_MINIMON_USE_BROAD_CMD
//{
//  unsigned __int16 nCharHandle;
//  char bUse;
//};
type NcCharUseItemMinimonUseBroadCmd struct {
	CharHandle uint16
	Use        byte
}

func (nc *NcCharUseItemMinimonUseBroadCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharUseItemMinimonUseBroadCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_USEITEM_MINIMON_USE_BROAD_CMD
	{
	  unsigned __int16 nCharHandle;
	  char bUse;
	};
`
}

//NC_CHAR_CLIENT_BASE_CMD
// idk why name is changed
//struct PROTO_NC_CHAR_BASE_CMD
//{
//  unsigned int chrregnum;
//  Name5 charid;
//  char slotno;
//  char Level;
//  unsigned __int64 Experience;
//  unsigned __int16 CurPwrStone;
//  unsigned __int16 CurGrdStone;
//  unsigned __int16 CurHPStone;
//  unsigned __int16 CurSPStone;
//  unsigned int CurHP;
//  unsigned int CurSP;
//  unsigned int CurLP;
//  unsigned int fame;
//  unsigned __int64 Cen;
//  PROTO_NC_CHAR_BASE_CMD::LoginLocation logininfo;
//  CHARSTATDISTSTR statdistribute;
//  char pkyellowtime;
//  unsigned int pkcount;
//  unsigned __int16 prisonmin;
//  char adminlevel;
//  PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags> flags;
//};
type NcCharClientBaseCmd struct {
	ChrRegNum uint32
	CharName Name5
	Slot byte
	Level byte
	Experience uint64
	PwrStone uint16
	GrdStone uint16
	HPStone uint16
	SPStone uint16
	CurHP uint32
	CurSP uint32
	CurLP uint32
	Unk  byte
	Fame uint32
	Cen uint64
	LoginInfo NcCharBaseCmdLoginLocation
	Stats CharStats
	IdleTime byte
	PkCount uint32
	PrisonMin uint16
	AdminLevel byte
	Flag NcCharBaseCmdFlag
	//MapName  Name3
}

func (nc *NcCharClientBaseCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharClientBaseCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_BASE_CMD
	{
	  unsigned int chrregnum;
	  Name5 charid;
	  char slotno;
	  char Level;
	  unsigned __int64 Experience;
	  unsigned __int16 CurPwrStone;
	  unsigned __int16 CurGrdStone;
	  unsigned __int16 CurHPStone;
	  unsigned __int16 CurSPStone;
	  unsigned int CurHP;
	  unsigned int CurSP;
	  unsigned int CurLP;
	  unsigned int fame;
	  unsigned __int64 Cen;
	  PROTO_NC_CHAR_BASE_CMD::LoginLocation logininfo;
	  CHARSTATDISTSTR statdistribute;
	  char pkyellowtime;
	  unsigned int pkcount;
	  unsigned __int16 prisonmin;
	  char adminlevel;
	  PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags> flags;
	};
`
}

//NC_CHAR_CLIENT_SHAPE_CMD
type NcCharClientShapeCmd ProtoAvatarShapeInfo