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
