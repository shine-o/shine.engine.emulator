package data

import (
	"testing"
)

var filesPath = "../../../files"

func TestLoadItemData(t *testing.T) {
	data, err := LoadItemData(filesPath)

	if err != nil {
		t.Fatal(err)
	}

	if data.ItemInfo == nil {
		t.Fail()
	}

	if data.ItemInfoServer == nil {
		t.Fail()
	}

	if data.GradeItemOptions == nil {
		t.Fail()
	}

	if data.RandomOption == nil {
		t.Fail()
	}

	if data.RandomOptionCount == nil {
		t.Fail()
	}

	//
	//if data.ItemDropGroup == nil {
	//	t.Fail()
	//}
	//
	//if data.ItemDropTable == nil {
	//	t.Fail()
	//}
}

type ItemData struct {
	ItemInfo          *ShineItemInfo
	ItemInfoServer    *ShineItemInfoServer
	GradeItemOptions  *ShineGradeItemOption
	RandomOption      *ShineRandomOption
	RandomOptionCount *ShineRandomOptionCount
}

func LoadItemData(path string) (ItemData, error) {
	var (
		itemData           ItemData
		itemInfo          = &ShineItemInfo{}
		itemInfoServer    = &ShineItemInfoServer{}
		gradeItemOptions  = &ShineGradeItemOption{}
		randomOption      = &ShineRandomOption{}
		randomOptionCount = &ShineRandomOptionCount{}
	)

	err := Load(path+"/shn/ItemInfo.shn", itemInfo)

	if err != nil {
		return itemData, err
	}

	err = Load(path+"/shn/ItemInfoServer.shn", itemInfoServer)

	if err != nil {
		return itemData, err
	}

	err = Load(path+"/shn/GradeItemOption.shn", gradeItemOptions)

	if err != nil {
		return itemData, err
	}

	err = Load(path+"/shn/RandomOption.shn", randomOption)

	if err != nil {
		return itemData, err
	}

	err = Load(path+"/shn/RandomOptionCount.shn", randomOptionCount)

	if err != nil {
		return itemData, err
	}

	itemData.ItemInfo = itemInfo
	itemData.ItemInfoServer = itemInfoServer
	itemData.GradeItemOptions = gradeItemOptions
	itemData.RandomOption = randomOption
	itemData.RandomOptionCount = randomOptionCount

	return itemData, nil
}

func TestLoadMobData(t *testing.T) {

}

func TestLoadNpcData(t *testing.T) {

}

func TestLoadMapData(t *testing.T) {

}
