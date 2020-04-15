package structs

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

// struct PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ
//{
//  char dummy[1];
//};
type NcCharOptionGetShortcutSizeReq struct {
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

//struct PROTO_NC_CHAR_GUILD_CMD
//{
//  unsigned int nGuildNo;
//  GUILD_CLIENT Guild[];
//};
type NcCharGuildCmd struct {
	GuildNumber uint32
	Guilds      GuildClient
}

//struct PROTO_NC_CHAR_LOGIN_REQ
//{
//  char slot;
//};
type NcCharLoginReq struct {
	Slot byte
}

//struct PROTO_NC_CHAR_MYSTERYVAULT_UI_STATE_CMD
//{
//  char mystery_vault;
//};
type CharMysteryVaultUiStateCmd struct {
	MysteryVault byte
}

//struct PROTO_NC_CHAR_USEITEM_MINIMON_INFO_CLIENT_CMD
//{
//  USEITEM_MINIMON_INFO UseItemMinimonInfo;
//};
type CharUseItemMiniMonsterInfoClientCmd struct {
	MiniMonsterInfo UseItemMiniMonsterInfo
}

//struct PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ
//{
//  char btMysteryType;
//};
type NcCharUiStateSaveReq struct {
	MysteryType byte
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

//struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD
//{
//  unsigned __int16 nGameOptionDataCnt;
//  GAME_OPTION_DATA GameOptionData[];
//};
type NcCharOptionImproveGetGameOptionCmd struct {
	OptionCount uint16
	GameOptions []GameOptionData `struct:"sizefrom=OptionCount"`
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

//struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD
//{
//unsigned __int16 nShortCutDataCnt;
//SHORT_CUT_DATA ShortCutData[];
//};
type NcCharGetShortcutDataCmd struct {
	Count     uint16
	Shortcuts []ShortCutData `struct:"sizefrom=Count"`
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

//struct PROTO_NC_CHAR_STAT_REMAINPOINT_CMD
//{
//  char remain;
//};
type NcCharStatRemainPointCmd struct {
	Remain byte
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

//struct PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD
//{
//  char nGuideView;
//};
type NcCharNewbieGuideViewSetCmd struct {
	GuideView byte
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

//struct PROTO_NC_CHAR_USEITEM_MINIMON_USE_BROAD_CMD
//{
//  unsigned __int16 nCharHandle;
//  char bUse;
//};
type NcCharUseItemMinimonUseBroadCmd struct {
	CharHandle uint16
	Use        byte
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
	ChrRegNum  uint32
	CharName   Name5
	Slot       byte
	Level      byte
	Experience uint64
	PwrStone   uint16
	GrdStone   uint16
	HPStone    uint16
	SPStone    uint16
	CurHP      uint32
	CurSP      uint32
	CurLP      uint32
	Unk        byte
	Fame       uint32
	Cen        uint64
	LoginInfo  NcCharBaseCmdLoginLocation
	Stats      CharStats
	IdleTime   byte
	PkCount    uint32
	PrisonMin  uint16
	AdminLevel byte
	Flag       NcCharBaseCmdFlag
	//MapName  Name3
}

//NC_CHAR_CLIENT_SHAPE_CMD
type NcCharClientShapeCmd ProtoAvatarShapeInfo

//struct PROTO_NC_CHAR_MAPLOGIN_ACK
//{
//  unsigned __int16 charhandle;
//  CHAR_PARAMETER_DATA param;
//  SHINE_XY_TYPE logincoord;
//};
type NcCharMapLoginAck struct {
	Handle uint16
	Param
}