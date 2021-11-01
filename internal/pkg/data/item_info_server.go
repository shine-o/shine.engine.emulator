package data

type ShineItemInfoServer struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn    `struct:"sizefrom=ColumnCount"`
	ShineRow    []ItemInfoServer `struct:"sizefrom=RowsCount"`
}

type ItemInfoServer struct {
	_                      uint16
	ID                     uint32
	InxName                string `struct:"[32]byte"`
	MarketIndex            string `struct:"[20]byte"`
	City                   byte
	DropGroupA             string `struct:"[40]byte"`
	DropGroupB             string `struct:"[40]byte"`
	RandomOptionDropGroup  string `struct:"[33]byte"`
	Vanish                 uint32
	Looting                uint32
	DropRateKilledByMob    uint16
	DropRateKilledByPlayer uint16
	ISETIndex              ISEType
	ItemSortIndex          string `struct:"[32]byte"`
	KQITem                 byte
	PkKqUse                byte
	KqItemDrop             byte
	PreventAttack          byte
}

type ISEType uint32

const (
	ISET_NONEEQUIP ISEType = iota
	ISET_MINIMON
	ISET_MINIMON_R
	ISET_MINIMON_BOTH
	ISET_COS_TAIL
	ISET_COS_BACK
	ISET_COS_RIGHT
	ISET_COS_LEFT
	ISET_COS_TWOHAND
	ISET_COS_HEAD
	ISET_COS_EYE
	ISET_COS_3PIECE_AMOR
	ISET_COS_3PIECE_PANTS
	ISET_COS_3PIECE_BOOTS
	ISET_COS_2PIECE_PANTS
	ISET_COS_1PIECE
	ISET_NORMAL_BOOTS
	ISET_NORMAL_PANTS
	ISET_RING
	ISET_SHIELD
	ISET_NORMAL_AMOR
	ISET_WEAPON_RIGHT
	ISET_WEAPON_TWOHAND
	ISET_WEAPON_LEFT
	ISET_EARRING
	ISET_NORMAL_HAT
	ISET_NECK
	ISET_COS_MASK
	ISET_INVINCIBLEHAMMER
	ISET_COS_MASK_EYE
	ISET_COS_HIDE_HEAD
	ISET_COS_EFF
	ISET_COS_SHIELD
	ISET_BRACELET
	MAX_ISETYPE
)
