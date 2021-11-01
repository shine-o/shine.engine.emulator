package data

type ShineGradeItemOption struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn     `struct:"sizefrom=ColumnCount"`
	ShineRow    []GradeItemOption `struct:"sizefrom=RowsCount"`
}

type GradeItemOption struct {
	_                  uint16
	ItemIndex          string `struct:"[32]byte"`
	Strength           uint16
	Endurance          uint16
	Dexterity          uint16
	Intelligence       uint16
	Spirit             uint16
	PoisonResistance   uint16
	DiseaseResistance  uint16
	CurseResistance    uint16
	MobilityResistance uint16
	AimRate            uint16
	EvasionRate        uint16
	MaxHP              uint16
	MaxSP              uint16
	PDamageIncrease    uint16
	MDamageIncrease    uint16
}
