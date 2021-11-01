package data

type ShineMobInfo struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn `struct:"sizefrom=ColumnCount"`
	ShineRow    []MobInfo     `struct:"sizefrom=RowsCount"`
}

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
