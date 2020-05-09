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
	Location ItemInventory
	ItemID   uint16
	ItemAttr []byte `struct-size:"DataSize - 4"`
}

//struct SHINE_ITEM_STRUCT
//{
//  unsigned __int16 itemid;
//  SHINE_ITEM_ATTRIBUTE itemattr;
//};
type ShineItem struct {
	ItemID uint16
	//Attr [101]byte
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
type GuildClient struct { // WRONG, 2020 uses different struct
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

//struct PROTO_NC_CHAR_BASE_CMD::LoginLocation
//{
//  Name3 currentmap;
//  SHINE_COORD_TYPE currentcoord;
//};
type NcCharBaseCmdLoginLocation struct {
	CurrentMap   Name3
	CurrentCoord ShineCoordType
}

//struct CHARSTATDISTSTR
//{
//  char Strength;
//  char Constitute;
//  char Dexterity;
//  char Intelligence;
//  char MentalPower;
//  char RedistributePoint;
//};
type CharStats struct {
	Strength          byte
	Constitute        byte
	Dexterity         byte
	Intelligence      byte
	MentalPower       byte
	RedistributePoint byte
}

//struct PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>::<unnamed-type-str>
//{
//  int _bf0;
//};
//
//union PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>
//{
//  unsigned int bin;
//  PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>::<unnamed-type-str> str;
//};
type NcCharBaseCmdFlag struct {
	Val int32
}

//struct CHAR_PARAMETER_DATA
//{
//  unsigned __int64 PrevExp;
//  unsigned __int64 NextExp;
//  SHINE_CHAR_STATVAR Strength;
//  SHINE_CHAR_STATVAR Constitute;
//  SHINE_CHAR_STATVAR Dexterity;
//  SHINE_CHAR_STATVAR Intelligence;
//  SHINE_CHAR_STATVAR Wizdom;
//  SHINE_CHAR_STATVAR MentalPower;
//  SHINE_CHAR_STATVAR WClow;
//  SHINE_CHAR_STATVAR WChigh;
//  SHINE_CHAR_STATVAR AC;
//  SHINE_CHAR_STATVAR TH;
//  SHINE_CHAR_STATVAR TB;
//  SHINE_CHAR_STATVAR MAlow;
//  SHINE_CHAR_STATVAR MAhigh;
//  SHINE_CHAR_STATVAR MR;
//  SHINE_CHAR_STATVAR MH;
//  SHINE_CHAR_STATVAR MB;
//  unsigned int MaxHp;
//  unsigned int MaxSp;
//  unsigned int MaxLp;
//  unsigned int MaxAp;
//  unsigned int MaxHPStone;
//  unsigned int MaxSPStone;
//  CHAR_PARAMETER_DATA::<unnamed-type-PwrStone> PwrStone;
//  CHAR_PARAMETER_DATA::<unnamed-type-PwrStone> GrdStone;
//  SHINE_CHAR_STATVAR PainRes;
//  SHINE_CHAR_STATVAR RestraintRes;
//  SHINE_CHAR_STATVAR CurseRes;
//  SHINE_CHAR_STATVAR ShockRes;
//};
type CharParameterData struct {
	// i'll have to rename these fields later when I can identify exactly what each field is for x.x
	PrevExp      uint64
	NextExp      uint64
	Strength     ShineCharStatVar
	Constitute   ShineCharStatVar
	Dexterity    ShineCharStatVar
	Intelligence ShineCharStatVar
	Wisdom       ShineCharStatVar
	MentalPower  ShineCharStatVar
	WCLow        ShineCharStatVar // min physical dmg
	WCHigh       ShineCharStatVar // max physical dmg
	AC           ShineCharStatVar // physical defense
	TH           ShineCharStatVar //aim
	TB           ShineCharStatVar //evasion
	MALow        ShineCharStatVar // min magical dmg
	MAHigh       ShineCharStatVar // max magical dmg
	MR           ShineCharStatVar // magical defense
	MH           ShineCharStatVar
	MB           ShineCharStatVar
	MaxHP        uint32
	MaxSP        uint32
	MaxLP        uint32
	MaxAP        uint32
	MaxHPStone   uint32
	MaxSPStone   uint32
	PwrStone     CharParameterDataPwrStone
	GrdStone     CharParameterDataPwrStone
	PainRes      ShineCharStatVar
	RestraintRes ShineCharStatVar
	CurseRes     ShineCharStatVar
	ShockRes     ShineCharStatVar
}

//struct CHAR_PARAMETER_DATA::<unnamed-type-PwrStone>
//{
//  unsigned int flag;
//  unsigned int EPPysic;
//  unsigned int EPMagic;
//  unsigned int MaxStone;
//};
type CharParameterDataPwrStone struct {
	Flag      uint32
	EPPPhysic uint32
	EPMagic   uint32
	MaxStone  uint32
}

//struct SHINE_CHAR_STATVAR
//{
//  unsigned int base;
//  unsigned int change;
//};
type ShineCharStatVar struct {
	Base   uint32
	Change uint32
}

//struct PLAYER_QUEST_INFO
//{
//  unsigned __int16 ID;
//  char Status;
//  PLAYER_QUEST_DATA Data;
//};
type PlayerQuestInfo struct {
	ID     uint16
	Status byte
	Data   PlayerQuestData
}

//struct PLAYER_QUEST_DATA
//{
//  __int64 StartTime;
//  __int64 EndTime;
//  unsigned int RepeatCount;
//  char ProgressStep;
//  char End_NPCMobCount[5];
//  char _bf26;
//  unsigned __int16 End_RunningTimeSec;
//};
type PlayerQuestData struct {
	StartTime         int64
	EndTime           int64
	RepeatCount       uint32
	ProgressStep      byte
	EndNpcMobCount    [5]byte
	BF26              byte
	EndRunningTimeSec uint16
}

//struct PLAYER_QUEST_DONE_INFO
//{
//  unsigned __int16 ID;
//  __int64 tEndTime;
//};
type PlayerQuestDoneInfo struct {
	ID      uint16
	EndTime int64
}

//struct PROTO_CHARGEDBUFF_INFO
//{
//  unsigned int ChargedBuffKey;
//  unsigned __int16 ChargedBuffID;
//  ShineDateTime UseTime;
//  ShineDateTime EndTime;
//};
type ChargedBuffInfo struct {
	BuffKey uint32
	BuffID  uint16
	UseTime ShineDateTime
	EndTime ShineDateTime
}
