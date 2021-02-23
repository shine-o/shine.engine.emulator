package shn

import (
	"fmt"
	"reflect"
)

type ShineItemInfoServer struct {
	DataSize    uint32
	RowsCount   uint32
	FieldSize   uint32
	ColumnCount uint32
	Columns     []ShineColumn    `struct:"sizefrom=ColumnCount"`
	ShineRow    []ItemInfoServer `struct:"sizefrom=RowsCount"`
}

//struct __unaligned __declspec(align(2)) ItemInfoServer
//{
//  unsigned int ID;
//  char InxName[32];
//  char MarketIndex[20];
//  char City[1];
//  char DropGroupA[40];
//  char DropGroupB[40];
//  char RandomOptionDropGroup[33];
//  unsigned int Vanish;
//  unsigned int looting;
//  unsigned __int16 DropRateKilledByMob;
//  unsigned __int16 DropRateKilledByPlayer;
//  ISEType ISET_Index;
//  char ItemSort_Index[32];
//  char KQItem;
//  char PK_KQ_USE;
//  char KQ_Item_Drop;
//  char PreventAttack;
//};
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

//enum ISEType
//{
//  ISET_NONEEQUIP = 0x0,
//  ISET_MINIMON = 0x1,
//  ISET_MINIMON_R = 0x2,
//  ISET_MINIMON_BOTH = 0x3,
//  ISET_COS_TAIL = 0x4,
//  ISET_COS_BACK = 0x5,
//  ISET_COS_RIGHT = 0x6,
//  ISET_COS_LEFT = 0x7,
//  ISET_COS_TWOHAND = 0x8,
//  ISET_COS_HEAD = 0x9,
//  ISET_COS_EYE = 0xA,
//  ISET_COS_3PIECE_AMOR = 0xB,
//  ISET_COS_3PIECE_PANTS = 0xC,
//  ISET_COS_3PIECE_BOOTS = 0xD,
//  ISET_COS_2PIECE_PANTS = 0xE,
//  ISET_COS_1PIECE = 0xF,
//  ISET_NORMAL_BOOTS = 0x10,
//  ISET_NORMAL_PANTS = 0x11,
//  ISET_RING = 0x12,
//  ISET_SHIELD = 0x13,
//  ISET_NORMAL_AMOR = 0x14,
//  ISET_WEAPON_RIGHT = 0x15,
//  ISET_WEAPON_TWOHAND = 0x16,
//  ISET_WEAPON_LEFT = 0x17,
//  ISET_EARRING = 0x18,
//  ISET_NORMAL_HAT = 0x19,
//  ISET_NECK = 0x1A,
//  ISET_COS_MASK = 0x1B,
//  ISET_INVINCIBLEHAMMER = 0x1C,
//  ISET_COS_MASK_EYE = 0x1D,
//  ISET_COS_HIDE_HEAD = 0x1E,
//  ISET_COS_EFF = 0x1F,
//  ISET_COS_SHIELD = 0x20,
//  ISET_BRACELET = 0x21,
//  MAX_ISETYPE = 0x22,
//};
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

func (s * ShineItemInfoServer) MissingIndexes(filePath string) (map[string][]string, error) {
	// have a function for each dependent file separately
	// ItemInfoServer
	var res = make(map[string][]string)

	var iis ShineItemInfo
	err := Load(filePath + "/shn/ItemInfo.shn", &iis)
	if err != nil {
		return res, err
	}

	res[reflect.TypeOf(iis).String()] = s.missingItemInfoIndex(&iis)

	return res, nil
}

func (s * ShineItemInfoServer) MissingIDs(filePath string) ( map[string][]uint16, error) {
	var res = make(map[string][]uint16)
	var iis ShineItemInfo
	err := Load(filePath + "/shn/ItemInfo.shn", &iis)
	if err != nil {
		return res, err
	}
	res[reflect.TypeOf(iis).String()] = s.missingItemInfoIDs(&iis)
	return res, nil
}

func (s * ShineItemInfoServer) MismatchedIndexAndID(filePath string) (map[string][]string, error){
	var res = make(map[string][]string)

	var iis ShineItemInfo
	err := Load(filePath + "/shn/ItemInfo.shn", &iis)
	if err != nil {
		return res, err
	}

	res[reflect.TypeOf(iis).String()] = s.itemInfoServerMismatchedIndexID(&iis)
	return res, nil
}

func (s * ShineItemInfoServer) itemInfoServerMismatchedIndexID(iis * ShineItemInfo) []string {
	var res []string

	for _, i := range s.ShineRow {

		//var (
		//	id uint16
		//	index string
		//)
		match := false
		for _, j := range iis.ShineRow {
			if i.InxName == j.InxName && uint16(i.ID) == j.ID {
				match = true
				break
			}

			if i.InxName == j.InxName && uint16(i.ID) != j.ID {
				break
			}

			if i.InxName != j.InxName && uint16(i.ID) == j.ID {
				break
			}
		}

		if !match {
			res = append(res, fmt.Sprintf("%v %v", i.ID, i.InxName))
		}
	}

	return res
}


func (s * ShineItemInfoServer) missingItemInfoIndex(iis * ShineItemInfo) []string {
	var res []string

	for _, i := range s.ShineRow {

		hasIndex := false
		for _, j := range iis.ShineRow {
			if i.InxName == j.InxName {
				hasIndex = true
				break
			}
		}

		if !hasIndex {
			res = append(res, i.InxName)
		}
	}
	return res
}

func (s * ShineItemInfoServer) missingItemInfoIDs(iis * ShineItemInfo) []uint16 {
	var res []uint16

	for _, i := range s.ShineRow {

		hasID := false
		for _, j := range iis.ShineRow {
			if uint16(i.ID) == j.ID {
				hasID = true
				break
			}
		}

		if !hasID {
			res = append(res, uint16(i.ID))
		}
	}
	return res
}