package data

type ShineRandomOption struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn  `struct:"sizefrom=ColumnCount"`
	ShineRow    []RandomOption `struct:"sizefrom=RowsCount"`
}

// stats
type RandomOption struct {
	_ uint16
	DropItemIndex    string `struct:"[33]byte"`
	RandomOptionType RandomOptionType
	Min              uint32
	Max              uint32
	TypeDropRate     uint32
}

//enum RandomOptionType
//{
//  ROT_STR = 0x0,
//  ROT_CON = 0x1,
//  ROT_DEX = 0x2,
//  ROT_INT = 0x3,
//  ROT_MEN = 0x4,
//  ROT_TH = 0x5,
//  ROT_CRI = 0x6,
//  ROT_WC = 0x7,
//  ROT_AC = 0x8,
//  ROT_MA = 0x9,
//  ROT_MR = 0xA,
//  ROT_TB = 0xB,
//  ROT_CRITICAL_TB = 0xC,
//  ROT_DEMANDLVDOWN = 0xD,
//  ROT_MAXHP = 0xE,
//  MAX_RANDOMOPTIONTYPE = 0xF,
//};

type RandomOptionType uint32

const (
	ROT_STR RandomOptionType = iota
	ROT_CON
	ROT_DEX
	ROT_INT
	ROT_MEN
	ROT_TH
	ROT_CRI
	ROT_WC
	ROT_AC
	ROT_MA
	ROT_MR
	ROT_TB
	ROT_CRITICAL_TB
	ROT_DEMANDLVDOWN
	ROT_MAXHP
	MAX_RANDOMOPTIONTYPE
)
