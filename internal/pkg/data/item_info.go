package data

const (
	itemInfoServer = "ItemInfoServer.shn"
)

type ShineItemInfo struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn `struct:"sizefrom=ColumnCount"`
	ShineRow    []ItemInfo    `struct:"sizefrom=RowsCount"`
}

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

type ItemClassEnum uint32

const (
	ItemClassByteLot ItemClassEnum = iota
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

// E_BelongType
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

func (s *ShineItemInfo) MissingIdentifiers(filesPath string) (Files, error) {
	var (
		res = Files{}
		iis = &ShineItemInfoServer{}
	)

	res[itemInfoServer] = Identifiers{}
	err := Load(filesPath+"/shn/ItemInfoServer.shn", iis)
	if err != nil {
		return res, err
	}

	itemInfoServerDeps(s, iis, res)

	return res, nil
}

func itemInfoServerDeps(s *ShineItemInfo, iis *ShineItemInfoServer, res Files) {
	for _, i := range s.ShineRow {

		ok := false
		for _, j := range iis.ShineRow {
			if i.InxName == j.InxName && i.ID == uint16(j.ID) {
				ok = true
				break
			}
		}

		if !ok {
			res[itemInfoServer]["ID"] = append(res[itemInfoServer]["ID"], i.ID)
			res[itemInfoServer]["InxName"] = append(res[itemInfoServer]["InxName"], i.InxName)
		}
	}
}
