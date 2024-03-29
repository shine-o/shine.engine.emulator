package data

type ShineChargedEffect struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn   `struct:"sizefrom=ColumnCount"`
	ShineRow    []ChargedEffect `struct:"sizefrom=RowsCount"`
}

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

type ChargedEffect struct {
	_            uint16
	Handle       uint16
	ItemIndex    string `struct:"[32]byte"`
	KeepTimeHour uint16
	EffectEnum   EffectEnumerate
	EffectValue  uint16
	StaStrength  byte
}
