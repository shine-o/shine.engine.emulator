package data

type ShineRandomOptionCount struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn       `struct:"sizefrom=ColumnCount"`
	ShineRow    []RandomOptionCount `struct:"sizefrom=RowsCount"`
}

// Amount of stats that can drop and the drop rate
type RandomOptionCount struct {
	_             uint16
	DropItemIndex string `struct:"[33]byte"`
	LimitCount    uint16
	LimitDropRate uint16
}
