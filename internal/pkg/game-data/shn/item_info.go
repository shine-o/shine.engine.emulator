package shn

import (
	"reflect"
)

type ShineItemInfo struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn `struct:"sizefrom=ColumnCount"`
	ShineRow    []ItemInfo    `struct:"sizefrom=RowsCount"`
}

//struct ItemInfo
//{
//  unsigned __int16 ID;
//  char InxName[32];
//  char Name[64];
//  ItemTypeEnum Type;
//  ItemClassEnum Class;
//  unsigned int MaxLot;
//  ItemEquipEnum Equip;
//  AuctionGroup ItemAuctionGroup;
//  GradeType ItemGradeType;
//  char TwoHand;
//  unsigned int AtkSpeed;
//  unsigned int DemandLv;
//  unsigned int Grade;
//  unsigned int MinWC;
//  unsigned int MaxWC;
//  unsigned int AC;
//  unsigned int MinMA;
//  unsigned int MaxMA;
//  unsigned int MR;
//  unsigned int TH;
//  unsigned int TB;
//  unsigned int WCRate;
//  unsigned int MARate;
//  unsigned int ACRate;
//  unsigned int MRRate;
//  unsigned int CriRate;
//  unsigned int CriMinWc;
//  unsigned int CriMaxWc;
//  unsigned int CriMinMa;
//  unsigned int CriMaxMa;
//  unsigned int CrlTB;
//  UseClassType UseClass;
//  unsigned int BuyPrice;
//  unsigned int SellPrice;
//  char BuyDemandLv;
//  unsigned int BuyFame;
//  unsigned int BuyGToken;
//  unsigned int BuyGBCoin;
//  WeaponTypeEnum WeaponType;
//  ArmorTypeEnum ArmorType;
//  char UpLimit;
//  unsigned __int16 BasicUpInx;
//  unsigned __int16 UpSucRatio;
//  unsigned __int16 UpLuckRatio;
//  char UpResource;
//  unsigned __int16 AddUpInx;
//  unsigned int ShieldAC;
//  unsigned int HitRatePlus;
//  unsigned int EvaRatePlus;
//  unsigned int MACriPlus;
//  unsigned int CriDamPlus;
//  unsigned int MagCriDamPlus;
//  E_BelongType BT_Inx;
//  char TitleName[32];
//  char ItemUseSkill[32];
//  char SetItemIndex[32];
//  ItemFuncEnum ItemFunc;
//};
type ItemInfo struct {
	_                uint16
	ID               uint16
	InxName          string `struct:"[32]byte"`
	Name             string `struct:"[64]byte"`
	Type             ItemTypeEnum
	Class            ItemClassEnum
	MaxLot           uint32
	Equip            ItemEquipEnum
	ItemAuctionGroup AuctionGroup
	ItemGradeType    GradeType
	TwoHand          byte
	AtkSpeed         uint32
	DemandLv         uint32
	Grade            uint32
	MinWC            uint32
	MaxWC            uint32
	AC               uint32
	MinMA            uint32
	MaxMA            uint32
	MR               uint32
	TH               uint32
	TB               uint32
	WCRate           uint32
	MARate           uint32
	ACRate           uint32
	MRRate           uint32
	CriRate          uint32
	CriMinWc         uint32
	CriMaxWc         uint32
	CriMinMa         uint32
	CriMaxMa         uint32
	CrlTB            uint32
	UseClass         UseClassType
	BuyPrice         uint32
	SellPrice        uint32
	BuyDemandLv      byte
	BuyFame          uint32
	BuyGToken        uint32
	BuyGBCoin        uint32
	WeaponType       WeaponTypeEnum
	ArmorType        ArmorTypeEnum
	UpLimit          byte
	BasicUpInx       uint16
	UpSucRatio       uint16
	UpLuckRatio      uint16
	UpResource       byte
	AddUpInx         uint16
	ShieldAC         uint32
	HitRatePlus      uint32
	EvaRatePlus      uint32
	MACriPlus        uint32
	CriDamPlus       uint32
	MagCriDamPlus    uint32
	BTInx            EBelongType
	TitleName        string `struct:"[32]byte"`
	ItemUseSkill     string `struct:"[32]byte"`
	SetItemIndex     string `struct:"[32]byte"`
	ItemFunc         ItemFuncEnum
}

//enum ItemTypeEnum
//{
//  ITEMTYPE_EQU = 0x0,
//  ITEMTYPE_EXH = 0x1,
//  ITEMTYPE_ETC = 0x2,
//  ITEMTYPE_QUEST = 0x3,
//  ITEMTYPE_STARTQUEST = 0x4,
//  ITEMTYPE_VIP = 0x5,
//  ITEMTYPE_CONFIRM = 0x6,
//  MAX_ITEMTYPEENUM = 0x7,
//};
type ItemTypeEnum uint32

const (
	ItemTypeEqu ItemTypeEnum = iota
	ItemTypeExh
	ItemTypeEtc
	ItemTypeQuest
	ItemTypeStartQuest
	ItemTypeVip
	ItemTypeConfirm
	MaxItemTypeEnum
)

//enum ItemClassEnum
//{
//  ITEMCLASS_BYTELOT = 0x0,
//  ITEMCLASS_WORDLOT = 0x1,
//  ITEMCLASS_DWRDLOT = 0x2,
//  ITEMCLASS_QUESTITEM = 0x3,
//  ITEMCLASS_AMULET = 0x4,
//  ITEMCLASS_WEAPON = 0x5,
//  ITEMCLASS_ARMOR = 0x6,
//  ITEMCLASS_SHIELD = 0x7,
//  ITEMCLASS_BOOT = 0x8,
//  ITEMCLASS_FURNITURE = 0x9,
//  ITEMCLASS_DECORATION = 0xA,
//  ITEMCLASS_SKILLSCROLL = 0xB,
//  ITEMCLASS_RECALLSCROLL = 0xC,
//  ITEMCLASS_BINDITEM = 0xD,
//  ITEMCLASS_UPSOURCE = 0xE,
//  ITEMCLASS_ITEMCHEST = 0xF,
//  ITEMCLASS_WTLICENCE = 0x10,
//  ITEMCLASS_KQ = 0x11,
//  ITEMCLASS_HOUSESKIN = 0x12,
//  ITEMCLASS_UPRED = 0x13,
//  ITEMCLASS_UPBLUE = 0x14,
//  ITEMCLASS_KQSTEP = 0x15,
//  ITEMCLASS_FEED = 0x16,
//  ITEMCLASS_RIDING = 0x17,
//  ITEMCLASS_AMOUNT = 0x18,
//  ITEMCLASS_UPGOLD = 0x19,
//  ITEMCLASS_COSWEAPON = 0x1A,
//  ITEMCLASS_ACTIONITEM = 0x1B,
//  ITEMCLASS_GBCOIN = 0x1C,
//  ITEMCLASS_CAPSULE = 0x1D,
//  ITEMCLASS_CLOSEDCARD = 0x1E,
//  ITEMCLASS_OPENCARD = 0x1F,
//  ITEMCLASS_MONEY = 0x20,
//  ITEMCLASS_NOEFFECT = 0x21,
//  ITEMCLASS_ENCHANT = 0x22,
//  ITEMCLASS_ACTIVESKILL = 0x23,
//  ITEMCLASS_PUP = 0x24,
//  ITEMCLASS_COSSHIELD = 0x25,
//  ITEMCLASS_BRACELET = 0x26,
//  MAX_ITEMCLASSENUM = 0x27,
//};
type ItemClassEnum uint32

const (
	ItemclassByteLot ItemClassEnum = iota
	ItemClassWordLot
	ItemClassDwrdLot
	ItemClassQuestItem
	ItemClassAmulet
	ItemClassWeapon
	ItemClassArmor
	ItemClassShield
	ItemClassBoot
	ItemClassFurniture
	ItemClassDecoration
	ItemClassSkillScroll
	ItemClassRecallScroll
	ItemClassBindItem
	ItemClassUpsource
	ItemClassItemChest
	ItemClassWtLicence
	ItemKq
	ItemHouseSkin
	ItemUpRed
	ItemUpBlue
	ItemKqStep
	ItemFeed
	ItemRiding
	ItemAmount
	ItemUpGold
	ItemCosWeapon
	ItemActionItem
	ItemGbCoin
	ItemCapsule
	ItemClosedCard
	ItemOpenCard
	ItemMoney
	ItemNoEffect
	ItemEnchant
	ItemActiveSkill
	ItemPup
	ItemCosShield
	ItemBracelet
	MaxItemEnum
)

//enum ItemEquipEnum
//{
//  ITEMEQUIP_NONE = 0x0,
//  ITEMEQUIP_HAT = 0x1,
//  ITEMEQUIP_NOUSE03 = 0x2,
//  ITEMEQUIP_NOUSE01 = 0x3,
//  ITEMEQUIP_NOUSE02 = 0x4,
//  ITEMEQUIP_FACETATTOO = 0x5,
//  ITEMEQUIP_NECKLACE = 0x6,
//  ITEMEQUIP_BODY = 0x7,
//  ITEMEQUIP_BODYACC = 0x8,
//  ITEMEQUIP_BACK = 0x9,
//  ITEMEQUIP_LEFTHAND = 0xA,
//  ITEMEQUIP_LEFTHANDACC = 0xB,
//  ITEMEQUIP_RIGHTHAND = 0xC,
//  ITEMEQUIP_RIGHTHANDACC = 0xD,
//  ITEMEQUIP_BRACELET = 0xE,
//  ITEMEQUIP_LEFTRING = 0xF,
//  ITEMEQUIP_RIGHTRING = 0x10,
//  ITEMEQUIP_COSEFF = 0x11,
//  ITEMEQUIP_TAIL = 0x12,
//  ITEMEQUIP_LEG = 0x13,
//  ITEMEQUIP_LEGACC = 0x14,
//  ITEMEQUIP_SHOES = 0x15,
//  ITEMEQUIP_SHOESACC = 0x16,
//  ITEMEQUIP_EARRING = 0x17,
//  ITEMEQUIP_MOUTH = 0x18,
//  ITEMEQUIP_MINIMON = 0x19,
//  ITEMEQUIP_EYE = 0x1A,
//  ITEMEQUIP_HATACC = 0x1B,
//  ITEMEQUIP_MINIMON_R = 0x1C,
//  ITEMEQUIP_SHIELDACC = 0x1D,
//  MAX_ITEMEQUIPENUM = 0x1E,
//};
type ItemEquipEnum uint32

const (
	ItemEquipNone ItemEquipEnum = iota
	ItemEquipHat
	ItemEquipNoUse03
	ItemEquipNoUse01
	ItemEquipNoUse02
	ItemEquipFaceTattoo
	ItemEquipNecklace
	ItemEquipBody
	ItemEquipBodyAcc
	ItemEquipBack
	ItemEquipLeftHand
	ItemEquipLeftHandAcc
	ItemEquipRightHand
	ItemEquipRightHandAcc
	ItemEquipBracelet
	ItemEquipLeftRing
	ItemEquipRightRing
	ItemEquipCosEff
	ItemEquipTail
	ItemEquipLeg
	ItemEquipLegAcc
	ItemEquipShoes
	ItemEquipShoesAcc
	ItemEquipEarRing
	ItemEquipMouth
	ItemEquipMinimon
	ItemEquipEye
	ItemEquipHatAcc
	ItemEquipMinimonR
	ItemEquipShieldAcc
	MaxItemEquipEnum
)

//enum AuctionGroup
//{
//  AG_ALL = 0x0,
//  AG_M_WEAPON = 0x1,
//  AG_M_ARMOR = 0x2,
//  AG_M_ACCESSORY = 0x3,
//  AG_M_PRODUCE = 0x4,
//  AG_M_ENCHANT = 0x5,
//  AG_M_RAW = 0x6,
//  AG_M_ETC = 0x7,
//  AG_S_ONEHANDSWORD = 0x8,
//  AG_S_TWOHANDSWORD = 0x9,
//  AG_S_AXE = 0xA,
//  AG_S_MACE = 0xB,
//  AG_S_HAMMER = 0xC,
//  AG_S_BOW = 0xD,
//  AG_S_CBOW = 0xE,
//  AG_S_STAFF = 0xF,
//  AG_S_WAND = 0x10,
//  AG_S_CLAW = 0x11,
//  AG_S_DSWORD = 0x12,
//  AG_S_FIGHTER = 0x13,
//  AG_S_CLERIC = 0x14,
//  AG_S_ARCHER = 0x15,
//  AG_S_MAGE = 0x16,
//  AG_S_JOKER = 0x17,
//  AG_S_NECK = 0x18,
//  AG_S_EARRING = 0x19,
//  AG_S_RING = 0x1A,
//  AG_S_SCROLL = 0x1B,
//  AG_S_POTION = 0x1C,
//  AG_S_STONE = 0x1D,
//  AG_S_FOOD = 0x1E,
//  AG_S_ENCHANT = 0x1F,
//  AG_S_PRODRAW = 0x20,
//  AG_S_FARM = 0x21,
//  AG_S_MOVER = 0x22,
//  AG_S_MINIHOUSE = 0x23,
//  AG_S_COSTUME = 0x24,
//  AG_S_ABILLITY = 0x25,
//  AG_S_EMOTION = 0x26,
//  AG_S_ETC = 0x27,
//  AG_S_BLADE = 0x28,
//  AG_S_SENTINEL = 0x29,
//  AG_S_BRACELET = 0x2A,
//  MAX_AUCTIONGROUP = 0x2B,
//};
type AuctionGroup uint32

const (
	AgAll AuctionGroup = iota
	AgMWeapon
	AgMArmor
	AgMAccessory
	AgMProduce
	AgMEnchant
	AgMRaw
	AgMEtc
	AgSOneHandSword
	AgSTwoHandSword
	AgSAxe
	AgSMace
	AgSHammer
	AgSBow
	AgSCBow
	AgSStaff
	AgSWand
	AgSClaw
	AgSDSword
	AgSFighter
	AgSCleric
	AgSArcher
	AgSMage
	AgSJoker
	AgSNeck
	AgSEarring
	AgSRing
	AgSScroll
	AgSPotion
	AgSStone
	AgSFood
	AgSEnchant
	AgSProdRaw
	AgSFarm
	AgSMover
	AgSMiniHouse
	AgSCostume
	AgSAbility
	AgSEmotion
	AgSEtc
	AgSBlade
	AgSSentinel
	AgSBracelet
	MaxAuctionGroup
)

//enum GradeType
//{
//  GT_NORMAL = 0x0,
//  GT_NAMED = 0x1,
//  GT_RARE = 0x2,
//  GT_UNIQUE = 0x3,
//  GT_CHARGE = 0x4,
//  GT_SET = 0x5,
//  GT_LEGENDARY = 0x6,
//  GT_MYTHIC = 0x7,
//  MAX_GRADETYPE = 0x8,
//};
type GradeType uint32

const (
	GtNormal GradeType = iota
	GtNamed
	GtRare
	GtUnique
	GtCharge
	GtSet
	GtLegendary
	GtMythic
	MaxGradeType
)

//enum UseClassType
//{
//  UCT_NONE = 0x0,
//  UCT_ALL = 0x1,
//  UCT_FIGHTER_ALL = 0x2,
//  UCT_CLEVERFIGHTER_AND_OVER = 0x3,
//  UCT_WARRIOR_AND_OVER = 0x4,
//  UCT_WARRIOR_OVER = 0x5,
//  UCT_GLADIATOR_ONLY = 0x6,
//  UCT_KNIGHT_ONLY = 0x7,
//  UCT_CLERIC_ALL = 0x8,
//  UCT_HIGHCLERIC_AND_OVER = 0x9,
//  UCT_PALADIN_AND_OVER = 0xA,
//  UCT_GUARDIAN_ONLY = 0xB,
//  UCT_HOLYKNIGHT_ONLY = 0xC,
//  UCT_PALADIN_OVER = 0xD,
//  UCT_ARCHER_ALL = 0xE,
//  UCT_HAWKARCHER_AND_OVER = 0xF,
//  UCT_SCOUT_AND_OVER = 0x10,
//  UCT_RANGER_ONLY = 0x11,
//  UCT_SHARPSHOOTER_ONLY = 0x12,
//  UCT_SCOUT_OVER = 0x13,
//  UCT_MAGE_ALL = 0x14,
//  UCT_WIZMAGE_AND_OVER = 0x15,
//  UCT_ENCHANTER_AND_OVER = 0x16,
//  UCT_WIZARD_ONLY = 0x17,
//  UCT_WARLOCK_ONLY = 0x18,
//  UCT_ENCHANTER_OVER = 0x19,
//  UCT_SENTINEL_EXCLUDE = 0x1A,
//  UCT_JOKER_ALL = 0x1B,
//  UCT_CHASER_AND_OVER = 0x1C,
//  UCT_CRUEL_AND_OVER = 0x1D,
//  UCT_ASSASSIN_ONLY = 0x1E,
//  UCT_CLOSER_ONLY = 0x1F,
//  UCT_CRUEL_OVER = 0x20,
//  UCT_SENTINEL_ALL = 0x21,
//  UCT_SAVIOR_ONLY = 0x22,
//  UCT_DEEPER_SKILL = 0x23,
//  UCT_SHIELD = 0x24,
//  UCT_CLASS_CHANGE = 0x25,
//  UCT_SHIELD_NOT_GLA = 0x26,
//  MAX_USECLASSTYPE = 0x27,
//};
type UseClassType uint32

const (
	UctNone UseClassType = iota
	UctAll
	UctFighterAll
	UctCleverFighterAndOver
	UctWarriorAndOver
	UctWarriorOver
	UctGladiatorOnly
	UctKnightOnly
	UctClericAll
	UctHighClericAndOver
	UctPaladinAndOver
	UctGuardianOnly
	UctHolyKnightOnly
	UctPaladinOver
	UctArcherAll
	UctHawkArcherAndOver
	UctScoutAndOver
	UctRangerOnly
	UctSharpshooterOnly
	UctScoutOver
	UctMageAll
	UctWizMageAndOver
	UctEnchanterAndOver
	UctWizardOnly
	UctWarlockOnly
	UctEnchanterOver
	UctSentinelExclude
	UctJokerAll
	UctChaserAndOver
	UctCruelAndOver
	UctAssassinOnly
	UctCloserOnly
	UctCruelOver
	UctSentinelAll
	UctSaviorOnly
	UctDeeperSkill
	UctShield
	UctClassChange
	UctShieldNotGla
	MaxUseClassType
)

//enum WeaponTypeEnum
//{
//  WT_NONE = 0x0,
//  WT_SWORD = 0x1,
//  WT_BOW = 0x2,
//  WT_STAFF = 0x3,
//  WT_AXE = 0x4,
//  WT_MACE = 0x5,
//  WT_SPIKE = 0x6,
//  WT_FIST = 0x7,
//  WT_BODY = 0x8,
//  WT_STONE = 0x9,
//  WT_CROSSBOW = 0xA,
//  WT_WAND = 0xB,
//  WT_SPEAR = 0xC,
//  WT_HAMMER = 0xD,
//  WT_SPECIAL = 0xE,
//  WT_PRODUCTIONTOOL = 0xF,
//  WT_INVINCIBLEHAMMER = 0x10,
//  WT_DSWORD = 0x11,
//  WT_CLAW = 0x12,
//  WT_BLADE = 0x13,
//  WT_RANGE_PY = 0x14,
//  WT_TSWORD = 0x15,
//  MAX_WEAPONTYPEENUM = 0x16,
//};
type WeaponTypeEnum uint32

const (
	WtNone WeaponTypeEnum = iota
	WtSword
	WtBow
	WtStaff
	WtAxe
	WtMace
	WtSpike
	WtFist
	WtBody
	WtStone
	WtCrossbow
	WtWand
	WtSpear
	WtHammer
	WtSpecial
	WtProductionTool
	WtInvincibleHammer
	WtDSword
	WtClaw
	WtBlade
	WtRangePy
	WtTSword
	MaxWeaponTypeEnum
)

//enum ArmorTypeEnum
//{
//  AT_NONE = 0x0,
//  AT_CLOTH = 0x1,
//  AT_LEATHER = 0x2,
//  AT_SCALE = 0x3,
//  AT_PLATE = 0x4,
//  AT_BONE = 0x5,
//  AT_HARDSKIN = 0x6,
//  AT_WEAKSKIN = 0x7,
//  AT_BARTSKIN = 0x8,
//  AT_GELSKIN = 0x9,
//  AT_FURSKIN = 0xA,
//  AT_SPECIAL = 0xB,
//  MAX_ARMORTYPEENUM = 0xC,
//};
type ArmorTypeEnum uint32

const (
	AtCloth ArmorTypeEnum = iota
	AtLeather
	AtScale
	AtPlate
	AtBone
	AtHardSkin
	AtWeakSkin
	AtBartSkin
	AtGelSkin
	AtFurSkin
	AtSpecial
	MaxArmorTypeEnum
)

//enum E_BelongType
//{
//  BT_COMMON = 0x0,
//  BT_NO_SELL = 0x1,
//  BT_NO_DROP = 0x2,
//  BT_NO_SELL_DROP = 0x3,
//  BT_ACC = 0x4,
//  BT_CHR = 0x5,
//  BT_ONLY_DEL = 0x6,
//  BT_NO_DEL = 0x7,
//  BT_PUTON_ACC = 0x8,
//  BT_PUTON_CHR = 0x9,
//  BT_NO_STORAGE = 0xA,
//  MAX_E_BELONGTYPE = 0xB,
//};
type EBelongType uint32

const (
	BtCommon EBelongType = iota
	BtNoSell
	BtNoDrop
	BtNoSellDrop
	BtAcc
	BtChr
	BtOnlyDel
	BtNoDel
	BtPutOnAcc
	BtPutOnChr
	BtNoStorage
	MaxEBelongType
)

//enum ItemFuncEnum
//{
//  ITEMFUNC_NONE = 0x0,
//  ITEMFUNC_ENDUREKIT_WC = 0x1,
//  ITEMFUNC_ENDUREKIT_F = 0x2,
//  ITEMFUNC_JUSTREVIVAL = 0x3,
//  ITEMFUNC_CHANGE_NAME = 0x4,
//  ITEMFUNC_CHANGE_RELATION = 0x5,
//  ITEMFUNC_PUTON_CLEAR = 0x6,
//  MAX_ITEMFUNCENUM = 0x7,
//};
type ItemFuncEnum uint32

const (
	ItemFuncNone ItemFuncEnum = iota
	ItemFuncEndureKitWc
	ItemFuncEndureKitF
	ItemFuncJustRevival
	ItemFuncChangeName
	ItemFuncChangeRelation
	ItemFuncPutOnClear
	MaxItemFuncEnum
)


func (s * ShineItemInfo) MissingIndexes(filePath string) (map[string][]string, error) {
	// have a function for each dependent file separately
	// ItemInfoServer
	var res = make(map[string][]string)

	var iis ShineItemInfoServer
	err := Load(filePath + "/shn/ItemInfoServer.shn", &iis)
	if err != nil {
		return res, err
	}

	res[reflect.TypeOf(iis).String()] = s.missingItemInfoServerIndex(&iis)

	return res, nil
}

func (s * ShineItemInfo) MissingIDs(filePath string) ( map[string][]uint16, error) {
	var res = make(map[string][]uint16)
	var iis ShineItemInfoServer
	err := Load(filePath + "/shn/ItemInfoServer.shn", &iis)
	if err != nil {
		return res, err
	}
	res[reflect.TypeOf(iis).String()] = s.missingItemInfoServerIDs(&iis)
	return res, nil
}

func (s * ShineItemInfo) missingItemInfoServerIndex(iis * ShineItemInfoServer) []string {
	var res []string

	for _, i := range s.ShineRow {

		hasIndex := false
		for _, j := range iis.ShineRow {
			if i.InxName == j.InxName {
				hasIndex = true
				break
			}
		}

		if !hasIndex {
			res = append(res, i.InxName)
		}
	}
	return res
}

func (s * ShineItemInfo) missingItemInfoServerIDs(iis * ShineItemInfoServer) []uint16 {
	var res []uint16

	for _, i := range s.ShineRow {

		hasID := false
		for _, j := range iis.ShineRow {
			if i.ID == uint16(j.ID) {
				hasID = true
				break
			}
		}

		if !hasID {
			res = append(res, i.ID)
		}
	}
	return res
}