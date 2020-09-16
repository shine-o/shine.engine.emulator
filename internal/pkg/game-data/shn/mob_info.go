package shn

type ShineMobInfo struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn `struct:"sizefrom=ColumnCount"`
	ShineRow    []MobInfo   `struct:"sizefrom=RowsCount"`
}

//struct MobInfo
//{
//  unsigned __int16 ID;
//  char InxName[32];
//  char Name[32];
//  unsigned int Level;
//  unsigned int MaxHP;
//  unsigned int WalkSpeed;
//  unsigned int RunSpeed;
//  char IsNPC;
//  unsigned int Size;
//  WeaponTypeEnum WeaponType;
//  ArmorTypeEnum ArmorType;
//  MobGradeType GradeType;
//  MobType Type;
//  char IsPlayerSide;
//  unsigned int AbsoluteSize;
//};
type MobInfo struct {
	_         uint16
	ID        uint16
	InxName   string `struct:"[32]byte"`
	Name      string `struct:"[32]byte"`
	Level     uint32
	MaxHP     uint32
	WalkSpeed uint32
	RunSpeed  uint32
	IsNPC     byte
	Size      uint32
	WeaponTypeEnum
	ArmorTypeEnum
	MobGradeType
	MobType
	IsPlayerSide byte
	AbsoluteSize uint32
}

//enum MobGradeType
//{
//  MGT_NORMAL = 0x0,
//  MGT_CHIEF = 0x1,
//  MGT_BOSS = 0x2,
//  MGT_HERO = 0x3,
//  MGT_ELITE = 0x4,
//  MGT_NONE = 0x5,
//  MAX_MOBGRADETYPE = 0x6,
//};
type MobGradeType uint32

const (
	MGT_NORMAL MobGradeType = iota
	MGT_CHIEF
	MGT_BOSS
	MGT_HERO
	MGT_ELITE
	MGT_NONE
	MAX_MOBGRADETYPE
)

//enum MobType
//{
//  MT_HUMAN = 0x0,
//  MT_MAGICLIFE = 0x1,
//  MT_SPIRIT = 0x2,
//  MT_BEAST = 0x3,
//  MT_ELEMENTAL = 0x4,
//  MT_UNDEAD = 0x5,
//  MT_NPC = 0x6,
//  MT_OBJECT = 0x7,
//  MT_MINE = 0x8,
//  MT_HERB = 0x9,
//  MT_WOOD = 0xA,
//  MT_NONAME = 0xB,
//  MT_NOTARGET = 0xC,
//  MT_NOTARGET2 = 0xD,
//  MT_GLDITEM = 0xE,
//  MT_FLAG = 0xF,
//  MT_DEVIL = 0x10,
//  MT_META = 0x11,
//  MT_NODAMAGE = 0x12,
//  MT_NODAMAGE2 = 0x13,
//  MT_NONAMEGATE = 0x14,
//  MT_BOX_HERB = 0x15,
//  MT_BOX_MINE = 0x16,
//  MT_GB_DICE = 0x17,
//  MT_NODAMAGE3 = 0x18,
//  MT_FRIEND = 0x19,
//  MT_GB_SLOTMACHINE = 0x1A,
//  MT_FRIENDDMGABSORB = 0x1B,
//  MT_DEVILDOM = 0x1C,
//  MT_NOTARGET3 = 0x1D,
//  MT_META2 = 0x1E,
//  MT_DWARF = 0x1F,
//  MT_MACHINE = 0x20,
//  MAX_MOBTYPE = 0x21,
//};
type MobType uint32

const (
	MT_HUMAN MobType = iota
	MT_MAGICLIFE
	MT_SPIRIT
	MT_BEAST
	MT_ELEMENTAL
	MT_UNDEAD
	MT_NPC
	MT_OBJECT
	MT_MINE
	MT_HERB
	MT_WOOD
	MT_NONAME
	MT_NOTARGET
	MT_NOTARGET2
	MT_GLDITEM
	MT_FLAG
	MT_DEVIL
	MT_META
	MT_NODAMAGE
	MT_NODAMAGE2
	MT_NONAMEGATE
	MT_BOX_HERB
	MT_BOX_MINE
	MT_GB_DICE
	MT_NODAMAGE3
	MT_FRIEND
	MT_GB_SLOTMACHINE
	MT_FRIENDDMGABSORB
	MT_DEVILDOM
	MT_NOTARGET3
	MT_META2
	MT_DWARF
	MT_MACHINE
	MAX_MOBTYPE
)
