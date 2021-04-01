package data

type ShineMobInfoServer struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn   `struct:"sizefrom=ColumnCount"`
	ShineRow    []MobInfoServer `struct:"sizefrom=RowsCount"`
}

//struct __unaligned __declspec(align(1)) MobInfoServer
//{
//  unsigned int ID;
//  char InxName[32];
//  char Visible;
//  unsigned __int16 AC;
//  unsigned __int16 TB;
//  unsigned __int16 MR;
//  unsigned __int16 MB;
//  EnemyDetect EnemyDetectType;
//  MobKillType MobKillInx;
//  unsigned int MonEXP;
//  unsigned __int16 EXPRange;
//  unsigned __int16 DetectCha;
//  char ResetInterval;
//  unsigned __int16 CutInterval;
//  unsigned int CutNonAT;
//  unsigned int FollowCha;
//  unsigned __int16 PceHPRcvDly;
//  unsigned __int16 PceHPRcv;
//  unsigned __int16 AtkHPRcvDly;
//  unsigned __int16 AtkHPRcv;
//  unsigned __int16 Str;
//  unsigned __int16 Dex;
//  unsigned __int16 Con;
//  unsigned __int16 Int;
//  unsigned __int16 Men;
//  MobRace MobRaceType;
//  char Rank;
//  unsigned int FamilyArea;
//  unsigned int FamilyRescArea;
//  char FamilyRescCount;
//  unsigned __int16 BloodingResi;
//  unsigned __int16 StunResi;
//  unsigned __int16 MoveSpeedResi;
//  unsigned __int16 FearResi;
//  char ResIndex[32];
//  unsigned __int16 KQKillPoint;
//  char Return2Regen;
//  char IsRoaming;
//  char RoamingNumber;
//  unsigned __int16 RoamingDistance;
//  unsigned __int16 RoamingRestTime;
//  unsigned __int16 MaxSP;
//  char BroadAtDead;
//  unsigned __int16 TurnSpeed;
//  unsigned __int16 WalkChase;
//  char AllCanLoot;
//  unsigned __int16 DmgByHealMin;
//  unsigned __int16 DmgByHealMax;
//  unsigned __int16 RegenInterval;
//};
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

//enum EnemyDetect
//{
//  ED_BOUT = 0x0,
//  ED_AGGRESSIVE = 0x1,
//  ED_NOBRAIN = 0x2,
//  ED_AGGRESSIVE2 = 0x3,
//  ED_AGGREESIVEALL = 0x4,
//  ED_ENEMYALLDETECT = 0x5,
//  MAX_ENEMYDETECT = 0x6,
//};
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

//enum MobKillType
//{
//  MKT_DEFAULT = 0x0,
//  MKT_ONLYSCRIPT = 0x1,
//  MKT_MOB = 0x2,
//  MAX_MOBKILLTYPE = 0x3,
//};
type MobKillType uint32

const (
	MKT_DEFAULT MobKillType = iota
	MKT_ONLYSCRIPT
	MKT_MOB
	MAX_MOBKILLTYPE
)

//enum MobRace
//{
//  MR_NONE = 0x0,
//  MR_PHINO = 0x1,
//  MR_MARA = 0x2,
//  MR_MARLONE = 0x3,
//  MR_SKEL = 0x4,
//  MR_CEM = 0x5,
//  MR_GOBLIN = 0x6,
//  MR_KARA = 0x7,
//  MR_KEEPER = 0x8,
//  MR_PI = 0x9,
//  MR_LIZARD = 0xA,
//  MR_TRUMPY = 0xB,
//  MR_ORC = 0xC,
//  MR_SLIME = 0xD,
//  MR_BOAR = 0xE,
//  MR_STAFF = 0xF,
//  MR_ARCHON = 0x10,
//  MR_STONIE = 0x11,
//  MR_INCUBUS = 0x12,
//  MR_TREE = 0x13,
//  MR_IMP = 0x14,
//  MR_VIVI = 0x15,
//  MR_KEBING = 0x16,
//  MR_GUARDIAN = 0x17,
//  MR_MINER = 0x18,
//  MR_BELLOW = 0x19,
//  MR_CAIMAN = 0x1A,
//  MR_RHINOCE = 0x1B,
//  MR_MUD = 0x1C,
//  MR_SLUG = 0x1D,
//  MR_SHADOW = 0x1E,
//  MR_CHAR = 0x1F,
//  MR_STATUE = 0x20,
//  MR_HELGA = 0x21,
//  MR_SPIRIT = 0x22,
//  MR_MAGRITE = 0x23,
//  MR_WOLF = 0x24,
//  MR_BEAR = 0x25,
//  MR_SPIDER = 0x26,
//  MR_MAND = 0x27,
//  MR_LICH = 0x28,
//  MR_POON = 0x29,
//  MR_DEPRAVITY = 0x2A,
//  MR_WIND = 0x2B,
//  MR_SELF = 0x2C,
//  MR_ELF = 0x2D,
//  MR_HONEYING = 0x2E,
//  MR_BOOGY = 0x2F,
//  MR_CRAB = 0x30,
//  MR_GUARD_NORMAL = 0x31,
//  MR_DEVILDOM = 0x32,
//  MR_SLAYER = 0x33,
//  MR_DARKARMY = 0x34,
//  MR_BKNIGHTS = 0x35,
//  MAX_MOBRACE = 0x36,
//};

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
