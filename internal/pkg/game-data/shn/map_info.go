package shn

import (
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

type ShineMapInfo struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn `struct:"sizefrom=ColumnCount"`
	Rows        []MapInfo     `struct:"sizefrom=RowsCount"`
}

//struct MapInfo
//{
//  unsigned __int16 ID;
//  Name3 MapName;
//  char Name[32];
//  WM_Link IsWMLink;
//  unsigned int RegenX;
//  unsigned int RegenY;
//  char KingdomMap;
//  Name3 MapFolderName;
//  char InSide;
//  unsigned int Sight;
//};
type MapInfo struct {
	_             uint16
	ID            uint16
	MapName       structs.Name3
	Name          string `struct:"[32]byte"`
	IsWMLink      WMLink
	RegenX        uint32
	RegenY        uint32
	KingdomMap    byte
	MapFolderName structs.Name3
	InSide        byte
	Sight         uint32
}

//enum WM_Link
//{
//  WM_NONE = 0x0,
//  WM_ROUN = 0x1,
//  WM_ROUCOS01 = 0x2,
//  WM_ROUCOS02 = 0x3,
//  WM_ROUCOS03 = 0x4,
//  WM_ROUVAL01 = 0x5,
//  WM_ROUVAL02 = 0x6,
//  WM_ELD = 0x7,
//  WM_ELDCEM01 = 0x8,
//  WM_ELDCEM02 = 0x9,
//  WM_ELDGBL01 = 0xA,
//  WM_ELDPRI01 = 0xB,
//  WM_ELDFOR01 = 0xC,
//  WM_ELDSLEEP01 = 0xD,
//  WM_URG = 0xE,
//  WM_ECHOCAVE = 0xF,
//  WM_WINDYCAVE = 0x10,
//  WM_GOLDCAVE = 0x11,
//  WM_URGFIRE01 = 0x12,
//  WM_URGSWA01 = 0x13,
//  WM_ELDGBL02 = 0x14,
//  WM_ELDPRI02 = 0x15,
//  WM_LINKFIELD01 = 0x16,
//  WM_LINKFIELD02 = 0x17,
//  WM_URG_ALRUIN = 0x18,
//  WM_ADLTHORN01 = 0x19,
//  WM_URGDARK01 = 0x1A,
//  WM_BERKAL01 = 0x1B,
//  WM_BERA_ = 0x1C,
//  WM_ADL = 0x1D,
//  WM_BERFRZ01 = 0x1E,
//  WM_BERVALE01 = 0x1F,
//  WM_ADLVAL01 = 0x20,
//  WM_TEVAL = 0x21,
//  WM_BATTLEFIELD = 0x22,
//  WM_TCAVE = 0x23,
//  WM_SER = 0x24,
//  MAX_WM_LINK = 0x25,
//};
type WMLink uint32

const (
	WmNone WMLink = iota
	WmRouN
	WmRouCos01
	WmRouCos02
	WmRouCos03
	WmRouVal01
	WmRouVal02
	WmEld
	WmEldCem01
	WmEldCem02
	WmEldGbl01
	WmEldPri01
	WmEldFor01
	WmEldSleep01
	WmUrg
	WmEchoCave
	WmWindyCave
	WmGoldCave
	WmUrgFire01
	WmUrgSwa01
	WmEldGbl02
	WmEldPri02
	WmLinkField01
	WmLinkField02
	WmUrgAlruin
	WmAdlThorn01
	WmUrgDark01
	WmBerKal01
	WmBera
	WmAdl
	WmBerFrz01
	WmBerVale01
	WmAdlVal01
	WmTeVal
	WmBattlefield
	WmTCave
	WmSer
	MaxWmLink
)
