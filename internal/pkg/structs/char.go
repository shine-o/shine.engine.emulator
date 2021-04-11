package structs

//struct PROTO_NC_CHAR_CLIENT_SKILL_CMD
type NcCharClientSkillCmd struct {
	RestEmpower byte
	PartMark    PartMark
	MaxNum      uint16
	Skills      NcCharSkillClientCmd
}

//struct PROTO_NC_CHAR_SKILLCLIENT_CMD
type NcCharSkillClientCmd struct {
	ChrRegNum uint32
	Number    uint16
	Skills    []SkillReadBlockClient `struct:"sizefrom=Number"`
}

//struct PROTO_NC_CHAR_CLIENT_ITEM_CMD
type NcCharClientItemCmd struct {
	NumOfItem byte
	Box       byte
	Flag      ProtoNcCharClientItemCmdFlag
	Items     []ProtoItemPacketInformation `struct:"sizefrom=NumOfItem"`
}

//struct PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD
type NcClientCharTitleCmd struct {
	CurrentTitle        byte
	CurrentTitleElement byte
	CurrentTitleMobID   uint16
	NumOfTitle          uint16
	Titles              []CharTitleInfo `struct:"sizefrom=NumOfTitle"`
}

// struct PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ
type NcCharOptionGetShortcutSizeReq struct {
}

//struct PROTO_NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK
type NcCharOptionGetShortcutSizeAck struct {
	Success byte
	Data    NcCharOptionShortcutSize
}

//struct PROTO_NC_CHAR_GUILD_CMD
type NcCharGuildCmd struct {
	GuildNumber uint32
	Guilds      GuildClient
}

//struct PROTO_NC_CHAR_LOGIN_REQ
type NcCharLoginReq struct {
	Slot byte
}

//struct PROTO_NC_CHAR_MYSTERYVAULT_UI_STATE_CMD
type CharMysteryVaultUiStateCmd struct {
	MysteryVault byte
}

//struct PROTO_NC_CHAR_USEITEM_MINIMON_INFO_CLIENT_CMD
type CharUseItemMiniMonsterInfoClientCmd struct {
	MiniMonsterInfo UseItemMiniMonsterInfo
}

//struct PROTO_NC_CHARSAVE_UI_STATE_SAVE_REQ
type NcCharUiStateSaveReq struct {
	MysteryType byte
}

//struct PROTO_NC_CHAR_LOGIN_ACK
type NcCharLoginAck struct {
	ZoneIP   Name4
	ZonePort uint16
}

//struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD
type NcCharOptionImproveGetGameOptionCmd struct {
	OptionCount uint16
	GameOptions []GameOptionData `struct:"sizefrom=OptionCount"`
}

//struct PROTO_NC_CHAR_GUILD_ACADEMY_CMD
type NcCharGuildAcademyCmd struct {
	GuildAcademyNo uint32
	//GuildAcademyNo uint16
	IsGuildAcademyMember byte
	//GuildAcademyClients []GuildAcademyClient `struct:"sizefrom=GuildAcademyNo"`
	GuildAcademyClient GuildAcademyClient
}

//struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD
type NcCharGetShortcutDataCmd struct {
	Count     uint16
	Shortcuts []ShortCutData `struct:"sizefrom=Count"`
}

//struct PROTO_NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD
type NcCharGetKeyMapCmd struct {
	Count uint16
	Keys  []KeyMapData `struct:"sizefrom=Count"`
}

//struct PROTO_NC_CHAR_STAT_REMAINPOINT_CMD
type NcCharStatRemainPointCmd struct {
	Remain byte
}

//struct PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK
type NcCharOptionGetWindowPosAck struct {
	Success byte
	Data    NcCharOptionWindowPos
}

//struct PROTO_NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD
type NcCharNewbieGuideViewSetCmd struct {
	GuideView byte
}

//struct PROTO_NC_CHAR_CLIENT_AUTO_PICK_CMD
type NcCharClientAutoPickCmd struct {
	Player uint16
	Enable byte
}

//struct PROTO_NC_CHAR_USEITEM_MINIMON_USE_BROAD_CMD
type NcCharUseItemMinimonUseBroadCmd struct {
	CharHandle uint16
	Use        byte
}

// NC_CHAR_CLIENT_BASE_CMD
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
}

//NC_CHAR_CLIENT_SHAPE_CMD
type NcCharClientShapeCmd ProtoAvatarShapeInfo

//struct PROTO_NC_CHAR_MAPLOGIN_ACK
type NcCharMapLoginAck struct {
	Handle     uint16
	Params     CharParameterData
	LoginCoord ShineXYType
}

//struct PROTO_NC_CHAR_REVIVEOTHER_CMD
type NcCharReviveOtherCmd struct {
	Link    NcCharReviveSameCmd
	Socket  NcCharLoginAck
	WorldID uint16
}

//struct PROTO_NC_CHAR_REVIVESAME_CMD
type NcCharReviveSameCmd struct {
	MapID    uint16
	Location ShineXYType
}

//NC_MAP_LINKOTHER_CMD
type NcMapLinkOtherCmd NcCharReviveOtherCmd

// NC_MAP_LINKSAME_CMD
type NcMapLinkSameCmd NcCharReviveSameCmd

//struct PROTO_NC_CHAR_SKILL_PASSIVE_CMD
type NcCharSkillPassiveCmd struct {
	Number   uint16
	Passives []uint16 `struct:"sizefrom=Number"`
}

//NC_CHAR_CLIENT_PASSIVE_CMD
type NcCharClientPassiveCmd NcCharSkillPassiveCmd

//struct PROTO_NC_CHAR_QUEST_READ_CMD
type NcCharQuestReadCmd struct {
	CharID          uint32
	NumOfReadQuests uint16
	Quests          []uint16 `struct:"sizefrom=NumOfReadQuests"`
}

//NC_CHAR_CLIENT_QUEST_READ_CMD
type NcCharClientQuestReadCmd NcCharQuestReadCmd

//struct PROTO_NC_CHAR_QUEST_DOING_CMD
type NcCharQuestDoingCmd struct {
	CharID          uint32
	NeedClear       byte
	NumOfDoingQuest byte
	Quests          []PlayerQuestInfo `struct:"sizefrom=NumOfDoingQuest"`
}

//NC_CHAR_CLIENT_QUEST_DOING_CMD
type NcCharClientQuestDoingCmd NcCharQuestDoingCmd

//struct PROTO_NC_CHAR_QUEST_DONE_CMD
type NcCharQuestDoneCmd struct {
	CharID             uint32
	TotalDoneQuest     uint16
	TotalDoneQuestSize uint16
	Count              uint16
	Index              uint16
	Quests             []PlayerQuestDoneInfo `struct:"sizefrom=Count"`
}

//NC_CHAR_CLIENT_QUEST_DONE_CMD
type NcCharClientQuestDoneCmd NcCharQuestDoneCmd

//struct PROTO_NC_CHAR_QUEST_REPEAT_CMD
type NcCharQuestRepeatCmd struct {
	CharID uint32
	Count  uint16
	Quests []PlayerQuestInfo `struct:"sizefrom=Count"`
}

//NC_CHAR_CLIENT_QUEST_REPEAT_CMD
type NcCharClientQuestRepeatCmd NcCharQuestRepeatCmd

//struct PROTO_NC_CHAR_CHARGEDBUFF_CMD
type NcCharChargedBuffCmd struct {
	Count uint16
	Buffs []ChargedBuffInfo `struct:"sizefrom=Count"`
}

//NC_CHAR_CLIENT_CHARGEDBUFF_CMD
type NcCharClientChargedBuffCmd NcCharChargedBuffCmd

//struct PROTO_NC_CHAR_COININFO_CMD
type NcCharCoinInfoCmd struct {
	Coin          uint64
	ExchangedCoin uint64
}

//NC_CHAR_CLIENT_COININFO_CMD
type NcCharClientCoinInfoCmd NcCharCoinInfoCmd

//NC_CHAR_CLIENT_GAME_CMD
type NcCharClientGameCmd struct {
	Filler0 uint16
	Filler1 uint16
}
