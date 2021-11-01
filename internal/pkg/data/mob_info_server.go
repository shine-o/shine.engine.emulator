package data

type ShineMobInfoServer struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn   `struct:"sizefrom=ColumnCount"`
	ShineRow    []MobInfoServer `struct:"sizefrom=RowsCount"`
}

type MobInfoServer struct {
	_       uint16
	ID      uint32
	InxName string `struct:"[32]byte"`
	Visible byte
	AC      uint16
	TB      uint16
	MR      uint16
	MB      uint16
	EnemyDetect
	MobKillType
	MonExp        uint32
	ExpRange      uint16
	DetectCha     uint16
	ResetInterval byte
	CutInterval   uint16
	CutNonAT      uint32
	FollowCha     uint32
	PceHPRcvDly   uint16
	PceHPRcv      uint16
	AtkHPRcvDly   uint16
	AtkHPRcv      uint16
	Str           uint16
	Dex           uint16
	End           uint16
	Int           uint16
	Spr           uint16
	MobRace
	Rank            byte
	FamilyArea      uint32
	FamilyRescArea  uint32
	FamilyRescCount byte
	BloodingResi    uint16
	StunResi        uint16
	MoveSpeedResi   uint16
	FearResi        uint16
	ResIndex        string `struct:"[32]byte"`
	KQKillPoint     uint16
	Return2Regen    byte
	IsRoaming       byte
	RoamingNumber   byte
	RoamingDistance uint16
	RoamingRestTime uint16
	MaxSP           uint16
	BroadAtDead     byte
	TurnSpeed       uint16
	WalkChase       uint16
	AllCanLoot      byte
	DmgByHealMin    uint16
	DmgByHealMax    uint16
	RegenInterval   uint16
}

type EnemyDetect uint32

const (
	ED_BOUT EnemyDetect = iota
	ED_AGGRESSIVE
	ED_NOBRAIN
	ED_AGGRESSIVE2
	ED_AGGREESIVEALL
	ED_ENEMYALLDETECT
	MAX_ENEMYDETECT
)

type MobKillType uint32

const (
	MKT_DEFAULT MobKillType = iota
	MKT_ONLYSCRIPT
	MKT_MOB
	MAX_MOBKILLTYPE
)

type MobRace uint32

const (
	MR_NONE MobRace = iota
	MR_PHINO
	MR_MARA
	MR_MARLONE
	MR_SKEL
	MR_CEM
	MR_GOBLIN
	MR_KARA
	MR_KEEPER
	MR_PI
	MR_LIZARD
	MR_TRUMPY
	MR_ORC
	MR_SLIME
	MR_BOAR
	MR_STAFF
	MR_ARCHON
	MR_STONIE
	MR_INCUBUS
	MR_TREE
	MR_IMP
	MR_VIVI
	MR_KEBING
	MR_GUARDIAN
	MR_MINER
	MR_BELLOW
	MR_CAIMAN
	MR_RHINOCE
	MR_MUD
	MR_SLUG
	MR_SHADOW
	MR_CHAR
	MR_STATUE
	MR_HELGA
	MR_SPIRIT
	MR_MAGRITE
	MR_WOLF
	MR_BEAR
	MR_SPIDER
	MR_MAND
	MR_LICH
	MR_POON
	MR_DEPRAVITY
	MR_WIND
	MR_SELF
	MR_ELF
	MR_HONEYING
	MR_BOOGY
	MR_CRAB
	MR_GUARD_NORMAL
	MR_DEVILDOM
	MR_SLAYER
	MR_DARKARMY
	MR_BKNIGHTS
	MAX_MOBRACE
)
