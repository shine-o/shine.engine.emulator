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
	Guilds      []GuildClient `struct:"sizefrom=GuildNumber"`
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
	ZoneIP Name4
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
	IsGuildAcademyMember byte
	GuildAcademyClients []GuildAcademyClient `struct:"sizefrom=GuildAcademyNo"`
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