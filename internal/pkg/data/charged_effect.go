package data

type ShineChargedEffect struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn   `struct:"sizefrom=ColumnCount"`
	ShineRow    []ChargedEffect `struct:"sizefrom=RowsCount"`
}

//enum EffectEnumerate
//{
//  EE_SILVERWINGCOOLTIME = 0x0,
//  EE_NOLOSTINKILLED = 0x1,
//  EE_MOREINVENTORY = 0x2,
//  EE_MORESTORAGE = 0x3,
//  EE_MOREBOOTHSLOT = 0x4,
//  EE_FASTMINING = 0x5,
//  EE_MOREHPSTONE = 0x6,
//  EE_MORESPSTONE = 0x7,
//  EE_MOREHSPSTONE = 0x8,
//  EE_COSTUM = 0x9,
//  EE_AKPOWER = 0xA,
//  EE_DPPOWER = 0xB,
//  EE_ALLPOWER = 0xC,
//  EE_HPINCREASE = 0xD,
//  EE_SPINCREASE = 0xE,
//  EE_ALLINCREASE = 0xF,
//  EE_DROP_RATE = 0x10,
//  EE_FEED = 0x11,
//  EE_EXP_RATE = 0x12,
//  EE_STATUS = 0x13,
//  EE_ITEMAT_RATE = 0x14,
//  EE_ITENDF_RATE = 0x15,
//  EE_ITEMAL_RATE = 0x16,
//  EE_PRODUCTSPEEDRATE = 0x17,
//  EE_PRODUCTMASTERYRATE = 0x18,
//  EE_PRODUCTALLRATE = 0x19,
//  EE_JUSTREVIVE_HP = 0x1A,
//  EE_ENDURE_KIT = 0x1B,
//  EE_WEAPON_MAXENDURE = 0x1C,
//  EE_ADDPRODSKILL = 0x1D,
//  EE_JUSTRREVIVAL = 0x1E,
//  EE_SETABSTATE = 0x1F,
//  EE_GBCOIN = 0x20,
//  EE_STORAGE_ANYWHERE = 0x21,
//  EE_LPINCREASE = 0x22,
//  EE_AUTOPATHFIND = 0x23,
//  EE_CHATCOLOR = 0x24,
//  EE_TERMEXTEND = 0x25,
//  EE_EXPHOLD = 0x26,
//  EE_CLASSCHANGE = 0x27,
//  EE_QEXP_RATE = 0x28,
//  MAX_EFFECTENUMERATE = 0x29,
//};
type EffectEnumerate uint32

const (
	EE_SILVERWINGCOOLTIME EffectEnumerate = iota
	EE_NOLOSTINKILLED
	EE_MOREINVENTORY
	EE_MORESTORAGE
	EE_MOREBOOTHSLOT
	EE_FASTMINING
	EE_MOREHPSTONE
	EE_MORESPSTONE
	EE_MOREHSPSTONE
	EE_COSTUM
	EE_AKPOWER
	EE_DPPOWER
	EE_ALLPOWER
	EE_HPINCREASE
	EE_SPINCREASE
	EE_ALLINCREASE
	EE_DROP_RATE
	EE_FEED
	EE_EXP_RATE
	EE_STATUS
	EE_ITEMAT_RATE
	EE_ITENDF_RATE
	EE_ITEMAL_RATE
	EE_PRODUCTSPEEDRATE
	EE_PRODUCTMASTERYRATE
	EE_PRODUCTALLRATE
	EE_JUSTREVIVE_HP
	EE_ENDURE_KIT
	EE_WEAPON_MAXENDURE
	EE_ADDPRODSKILL
	EE_JUSTRREVIVAL
	EE_SETABSTATE
	EE_GBCOIN
	EE_STORAGE_ANYWHERE
	EE_LPINCREASE
	EE_AUTOPATHFIND
	EE_CHATCOLOR
	EE_TERMEXTEND
	EE_EXPHOLD
	EE_CLASSCHANGE
	EE_QEXP_RATE
	MAX_EFFECTENUMERATE
)

//struct ChargedItemEffect
//{
//  unsigned __int16 Handle;
//  char ItemID[32];
//  unsigned __int16 KeepTime_Hour;
//  EffectEnumerate EffectEnum;
//  unsigned __int16 EffectValue;
//  char StaStrength;
//};

type ChargedEffect struct {
	_            uint16
	Handle       uint16
	ItemIndex    string `struct:"[32]byte"`
	KeepTimeHour uint16
	EffectEnum   EffectEnumerate
	EffectValue  uint16
	StaStrength  byte
}
