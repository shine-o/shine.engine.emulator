package structs

// struct PROTO_SKILLREADBLOCKCLIENT
type SkillReadBlockClient struct {
	SkillID  uint16
	CoolTime uint32
	Empower  SkillReadBlockClientEmpower
	Mastery  uint32
}

// struct PROTO_SKILLREADBLOCKCLIENT::<unnamed-type-empow>
type SkillReadBlockClientEmpower struct {
	Gap0 byte
	BF1  byte
}

// struct PARTMARK
type PartMark struct {
	BF0 byte
}

// struct PROTO_NC_CHAR_CLIENT_ITEM_CMD::<unnamed-type-flag>
type ProtoNcCharClientItemCmdFlag struct {
	BF0 byte
}

// struct PROTO_ITEMPACKET_INFORM
type ProtoItemPacketInformation struct {
	DataSize byte
	Location ItemInventory
	ItemID   uint16
	ItemAttr []byte `struct-size:"DataSize - 4"`
}

// struct ShineItemAttr_ByteLot
type ShineItemAttrByteLot byte

type ShineItemAttrWordLot uint16

type ShineItemAttrDwrdLot uint32

type ShineItemAttrQuestItem uint16

// struct ShineItemAttr_Amulet
type ShineItemAttrAmulet struct {
	DeleteTime       int32
	IsBound          uint32
	Upgrade          byte
	Strengthen       byte
	UpgradeFailCount byte
	// this is set
	UpgradeOption            UpgradeItemOptionStorage
	RandomOptionChangedCount byte
	// this is dynamic
	Option ItemOptionStorage
}

type ItemOptionStorageFixedInfo struct {
	// this might be the Count of stats
	// statType >> 1
	StatType byte
}

type ItemOptionStorageElement struct {
	ItemOptionType  byte
	ItemOptionValue uint16
}

type ItemOptionStorage struct {
	AmountBit byte
	Elements  []ItemOptionStorageElement `struct-size:"AmountBit >> 1"`
}

type UpgradeItemOptionStorage struct {
	// if its fixed stats or generated

	FixedStat ItemOptionStorageFixedInfo
	// ItemOptionStorageElement max 8
	Elements [8]ItemOptionStorageElement
}

// struct ShineItemAttr_Weapon
type ShineItemAttrWeapon struct {
	Upgrade            byte
	Strengthen         byte
	UpgradeFailCount   byte
	IsBound            uint32
	Licences           [3]ShineItemWeaponLicence
	WeaponLicenceTitle uint16
	UserTitle          [21]byte
	GemSockets         [3]ShineItemWeaponGemSocket
	MaxSocketCount     byte
	CreatedSocketCount byte
	DeleteTime         int32
	// Bijou hammer usages
	RandomOptionChangedCount byte
	Option                   ItemOptionStorage
}

type ShineItemWeaponLicence struct {
	MobID uint16
	BF2   int32
}

type ShineItemWeaponGemSocket struct {
	GemID     uint16
	RestCount byte
}

// struct ShineItemAttr_Armor
type ShineItemAttrArmor struct {
	Upgrade                  byte
	Strengthen               byte
	UpgradeFailCount         byte
	IsBound                  uint32
	DeleteTime               int32
	RandomOptionChangedCount byte
	Option                   ItemOptionStorage
}

// struct ShineItemAttr_Shield
type ShineItemAttrShield struct {
	Upgrade                  byte
	Strengthen               byte
	UpgradeFailCount         byte
	IsBound                  uint32
	DeleteTime               int32
	RandomOptionChangedCount byte
	Option                   ItemOptionStorage
}

// struct ShineItemAttr_Boot
type ShineItemAttrBoot struct {
	Upgrade                  byte
	Strengthen               byte
	UpgradeFailCount         byte
	IsBound                  uint32
	DeleteTime               int32
	RandomOptionChangedCount byte
	Option                   ItemOptionStorage
}

// struct ShineItemAttr_Furniture
type ShineItemAttrFurniture struct {
	Flag           byte
	FurnitureID    uint16
	DeleteTime     int32
	LocX           float32
	LocY           float32
	LocZ           float32
	Direction      float32
	ExpirationTime int32
	EndureGrade    byte
	RewardMoney    uint64
}

// struct ShineItemAttr_Decoration
type ShineItemAttrDecoration struct {
	IsBound    uint32
	DeleteTime int32
}

// struct ShineItemAttr_BindItem
type ShineItemAttrBindItem struct {
	PortalNum byte
	Portals   [10]Bind
}

type Bind struct {
	MapID uint16
	X     uint32
	Y     uint32
}

// struct ShineItemAttr_ItemChest
type ShineItemAttrItemChest struct {
	Type    byte
	Content [8][8]byte
}

// struct ShineItemAttr_MiniHouseSkin
type ShineItemAttrMiniHouseSkin struct {
	DeleteTime int32
}

// struct ShineItemAttr_Riding
type ShineItemAttrRiding struct {
	HungryPoints  uint16
	DeleteTime    int32
	RidingFlag    uint16
	IsBound       uint32
	HP            uint32
	Grade         byte
	RareFailCount uint16
}

// ShineItemAttr_CostumWeapon
type ShineItemAttrCostumeWeapon struct {
	Durability uint32
}

// struct ShineItemAttr_ActionItem
type ShineItemAttrActionItem struct {
	DeleteTime int32
}

// struct ShineItemAttr_Capsule
type ShineItemAttrCapsule struct {
	Content     [8]byte
	UseAbleTime int32
}

// struct ShineItemAttr_MobCardCollect_Unident
type ShineItemAttrMobCardCollectClosed struct {
	SerialNumber uint32
	CardID       uint16
	Star         byte
	Group        uint16
}

// ShineItemAttr_MobCardCollect
type ShineItemAttrMobCardCollect struct {
	SerialNumber uint32
	Start        byte
}

// struct ShineItemAttr_Amount
type ShineItemAttrAmount struct {
	Amount uint32
}

// struct ShineItemAttr_Pet
type ShineItemAttrPet struct {
	PetRegNum uint32
	PetID     uint32
	Name      [17]byte
	Summoning byte
}

// struct ShineItemAttr_Bracelet
type ShineItemAttrBracelet struct {
	DeleteTime               int32
	IsBound                  uint32
	Upgrade                  byte
	Strengthen               byte
	UpgradeFailCount         byte
	RandomOptionChangedCount byte
	Option                   ItemOptionStorage
}

// struct ShineItemAttrCostumeShield
type ShineItemAttrCostumeShield struct {
	Durability uint32
}

type ShineItem struct {
	ItemID uint16
	// Attr [101]byte
}

// struct CT_INFO
type CharTitleInfo struct {
	Type byte
	BF1  byte
}

// struct PROTO_NC_CHAR_OPTION_SHORTCUTSIZE
type NcCharOptionShortcutSize struct {
	Data [24]byte
}

// struct GUILD_CLIENT
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

// struct SHINE_GUILD_SCORE
type ShineGuildScore struct {
	KillCount [7]uint16
	KillScore [7]uint32
}

// struct tm
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

// struct USEITEM_MINIMON_INFO
type UseItemMiniMonsterInfo struct {
	NormalItem      byte
	ChargedItem     byte
	NormalItemList  [12]uint16
	ChargedItemList [12]uint16
}

// struct GAME_OPTION_DATA
type GameOptionData struct {
	OptionNo uint16
	Value    byte
}

// struct GUILD_ACADEMY_CLIENT
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

// struct SHORT_CUT_DATA
type ShortCutData struct {
	SlotNo byte
	// 0 = remove item
	// 4 = add item
	CodeNo uint16
	// action index
	Value int32
}

// struct KEY_MAP_DATA
type KeyMapData struct {
	FunctionNo uint16
	ExtendKey  byte
	AsciiCode  byte
}

// struct PROTO_NC_CHAR_OPTION_WINDOWPOS
type NcCharOptionWindowPos struct {
	Data [392]byte
}

// struct PROTO_NC_CHAR_BASE_CMD::LoginLocation
type NcCharBaseCmdLoginLocation struct {
	CurrentMap   Name3
	CurrentCoord ShineCoordType
}

// struct CHARSTATDISTSTR
type CharStats struct {
	Strength          byte
	Constitute        byte
	Dexterity         byte
	Intelligence      byte
	MentalPower       byte
	RedistributePoint byte
}

// struct PROTO_NC_CHAR_BASE_CMD::<unnamed-type-flags>::<unnamed-type-str>
type NcCharBaseCmdFlag struct {
	Val int32
}

// struct CHAR_PARAMETER_DATA
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
	TH           ShineCharStatVar // aim
	TB           ShineCharStatVar // evasion
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

// struct CHAR_PARAMETER_DATA::<unnamed-type-PwrStone>
type CharParameterDataPwrStone struct {
	Flag      uint32
	EPPPhysic uint32
	EPMagic   uint32
	MaxStone  uint32
}

// struct SHINE_CHAR_STATVAR
type ShineCharStatVar struct {
	Base   uint32
	Change uint32
}

// struct PLAYER_QUEST_INFO
type PlayerQuestInfo struct {
	ID     uint16
	Status byte
	Data   PlayerQuestData
}

// struct PLAYER_QUEST_DATA
type PlayerQuestData struct {
	StartTime         int64
	EndTime           int64
	RepeatCount       uint32
	ProgressStep      byte
	EndNpcMobCount    [5]byte
	BF26              byte
	EndRunningTimeSec uint16
}

// struct PLAYER_QUEST_DONE_INFO
type PlayerQuestDoneInfo struct {
	ID      uint16
	EndTime int64
}

// struct PROTO_CHARGEDBUFF_INFO
type ChargedBuffInfo struct {
	BuffKey uint32
	BuffID  uint16
	UseTime ShineDateTime
	EndTime ShineDateTime
}
