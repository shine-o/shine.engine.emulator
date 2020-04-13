package structs

//struct PROTO_SKILLREADBLOCKCLIENT
//{
//	unsigned __int16 skillid;
//	unsigned int cooltime;
//	PROTO_SKILLREADBLOCKCLIENT::<unnamed-type-empow> empow;
//	unsigned int mastery;
//};
type SkillReadBlockClient struct {
	SkillID  uint16
	CoolTime uint32
	Empower  SkillReadBlockClientEmpower
	Mastery  uint32
}

//struct PROTO_SKILLREADBLOCKCLIENT::<unnamed-type-empow>
//{
//  _BYTE gap0[1];
//  char _bf1;
//};
type SkillReadBlockClientEmpower struct {
	Gap0 byte
	BF1  byte
}

//struct PARTMARK
//{
//	char _bf0;
//};
type PartMark struct {
	BF0 byte
}

//struct PROTO_NC_CHAR_CLIENT_ITEM_CMD::<unnamed-type-flag>
//{
//	char _bf0;
//};
type ProtoNcCharClientItemCmdFlag struct {
	BF0 byte
}

//struct PROTO_ITEMPACKET_INFORM
//{
//	char datasize;
//	ITEM_INVEN location;
//	SHINE_ITEM_STRUCT info;
//};
type ProtoItemPacketInformation struct {
	DataSize byte
	// can't be done like this, since data size also covers Location and Info and there's no way to use sizefrom with operators -+ :(
	// at the handler level, i would have to read the fields manually.
	ItemData []byte `struct:"sizefrom=DataSize"`
}

//struct CT_INFO
//{
//  char Type;
//  char _bf1;
//};
type CharTitleInfo struct {
	Type byte
	BF1  byte
}

//struct PROTO_NC_CHAR_OPTION_SHORTCUTSIZE
//{
//  char Data[24];
//};
type NcCharOptionShortcutSize struct {
	Data [24]byte
}

//struct GUILD_CLIENT
//{
//  unsigned int nNo;
//  Name4 sName;
//  unsigned __int64 nMoney;
//  char nType;
//  char nGrade;
//  unsigned int nFame;
//  unsigned __int16 nStoneLevel;
//  unsigned __int64 nExp;
//  int dCreateDate;
//  tm tm_dCreateDate;
//  unsigned __int16 nNumMembers;
//  unsigned __int16 nMaxMembers;
//  char nWarStatus;
//  int dWarRequestDate;
//  int dWarStartDate;
//  int dWarEndDate;
//  tm tm_dWarRequestDate;
//  tm tm_dWarStartDate;
//  tm tm_dWarEndDate;
//  unsigned int nWarEnemy;
//  Name4 sWarEnemyName;
//  char nWarEnemyGrade;
//  SHINE_GUILD_SCORE MyScore;
//  SHINE_GUILD_SCORE EnemyScore;
//  unsigned int nWarWinCount;
//  unsigned int nWarLoseCount;
//  unsigned int nWarDrawCount;
//  char nDismissStatus;
//  int dDismissDate;
//  tm tm_dDismissDate;
//  char sIntro[128];
//  int dNotifyDate;
//  tm tm_dNotifyDate;
//  Name5 sNotifyCharID;
//  char sNotify[512];
//};
type GuildClient struct {
	Number           uint32
	Name             Name4
	Money            uint64
	Type             byte
	Grade            byte
	Fame             uint32
	StoneLevel       uint16
	Exp              uint64
	CreatedDate      int32
	NumMembers       uint16
	MaxMembers       uint16
	WarStatus        byte
	WarRequestDate   int32
	WarStartDate     int32
	WarEndDate       int32
	TmWarRequestDate TM
	TmWarStartDate   TM
	TmWarEndDate     TM
	WarEnemy         uint32
	WarEnemyName     Name4
	WarEnemyGrade    byte
	MyScore          ShineGuildScore
	EnemyScore       ShineGuildScore
	WarWinCount      uint32
	WarLoseCount     uint32
	WarDrawCount     uint32
	DismissStatus    byte
	DismissDate      int32
	TmDismissDate    TM
	Intro            [128]byte
	NotifyDate       int32
	TmNotifyDate     TM
	NotifyCharID     Name5
	Notify           [512]byte
}

//struct SHINE_GUILD_SCORE
//{
//  unsigned __int16 nKillCount[7];
//  unsigned int nKillScore[7];
//};
type ShineGuildScore struct {
	KillCount [7]uint16
	KillScore [7]uint32
}

//struct tm
//{
//  int tm_sec;
//  int tm_min;
//  int tm_hour;
//  int tm_mday;
//  int tm_mon;
//  int tm_year;
//  int tm_wday;
//  int tm_yday;
//  int tm_isdst;
//};
type TM struct {
	Seconds  int32
	Minutes  int32
	Hour     int32
	MonthDay int32
	Month    int32
	Year     int32
	WeekDay  int32
	YearDay  int32
	IsDst    int32
}

//struct USEITEM_MINIMON_INFO
//{
//  char bNormalItem;
//  char bChargedItem;
//  unsigned __int16 NormalItemList[12];
//  unsigned __int16 ChargedItemList[12];
//};
type UseItemMiniMonsterInfo struct {
	NormalItem      byte
	ChargedItem     byte
	NormalItemList  [12]uint16
	ChargedItemList [12]uint16
}

//struct GAME_OPTION_DATA
//{
//  unsigned __int16 nOptionNo;
//  char nValue;
//};
type GameOptionData struct {
	OptionNo uint16
	Value    byte
}

//struct GUILD_ACADEMY_CLIENT
//{
//  Name5 sAcademyMasterName;
//  unsigned __int16 nNumAcademyMembers;
//  unsigned __int16 nMaxAcademyMembers;
//  unsigned int nAcademyPoint;
//  unsigned int nAcademyRank;
//  int dAcademyBuffUntilTime;
//  char sIntro[128];
//  int dNotifyDate;
//  tm tm_dNotifyDate;
//  Name5 sNotifyCharID;
//  char sNotify[512];
//};
type GuildAcademyClient struct {
	AcademyMasterName Name5
	NumAcademyMembers uint16
	MaxAcademyMembers uint16
	AcademyPoint      uint32
	AcademyRank       uint32
	AcademyBuffTime   int32
	Intro             [128]byte
	NotifyDate        int32
	TmNotifyDate      TM
	NotifyCharID      Name5
	Notify            [512]byte
}

//struct SHORT_CUT_DATA
//{
//  char nSlotNo;
//  unsigned __int16 nCodeNo;
//  int nValue;
//};
type ShortCutData struct {
	SlotNo byte
	CodeNo uint16
	Value  int32
}

//struct KEY_MAP_DATA
//{
//  unsigned __int16 nFunctionNo;
//  char nExtendKey;
//  char nASCIICode;
//};
type KeyMapData struct {
	FunctionNo uint16
	ExtendKey  byte
	AsciiCode  byte
}

//struct PROTO_NC_CHAR_OPTION_WINDOWPOS
//{
//  char Data[392];
//};
type NcCharOptionWindowPos struct {
	Data [392]byte
}
