package data

type ShineGradeItemOption struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn      `struct:"sizefrom=ColumnCount"`
	ShineRow    []GradeItemOption `struct:"sizefrom=RowsCount"`
}

//struct GradeItemOption
//{
//  char ItemIndex[32];
//  unsigned __int16 STR;
//  unsigned __int16 CON;
//  unsigned __int16 DEX;
//  unsigned __int16 INT;
//  unsigned __int16 MEN;
//  unsigned __int16 ResistPoison;
//  unsigned __int16 ResistDeaseas;
//  unsigned __int16 ResistCurse;
//  unsigned __int16 ResistMoveSpdDown;
//  unsigned __int16 ToHitRate;
//  unsigned __int16 ToBlockRate;
//  unsigned __int16 MaxHP;
//  unsigned __int16 MaxSP;
//  unsigned __int16 WCPlus;
//  unsigned __int16 MAPlus;
//};
type GradeItemOption struct {
	_         uint16
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
