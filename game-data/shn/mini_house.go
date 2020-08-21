package shn

type ShineMiniHouse struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn `struct:"sizefrom=ColumnCount"`
	ShineRow    []MiniHouse    `struct:"sizefrom=RowsCount"`
}

//struct MiniHouse
//{
//  unsigned __int16 Handle;
//  char ItemID[32];
//  char DummyType[32];
//  char Backimage[32];
//  unsigned __int16 KeepTime_Hour;
//  unsigned __int16 HPTick;
//  unsigned __int16 SPTick;
//  unsigned __int16 HPRecovery;
//  unsigned __int16 SPRecovery;
//  unsigned __int16 Casting;
//  char Slot;
//};
type MiniHouse struct {
	_ uint16
	Handle uint16
	ItemID string `struct:"[32]byte"`
	DummyType string `struct:"[32]byte"`
	BackImage string `struct:"[32]byte"`
	KeepTimeHour uint16
	HPTick uint16
	SPTick uint16
	HPRecovery uint16
	SPRecovery uint16
	Casting uint16
	Slot byte
}

