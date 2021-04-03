package data


type ItemData struct {
	ItemInfo          ShineItemInfo
	ItemInfoServer    ShineItemInfoServer
	GradeItemOptions  ShineGradeItemOption
	RandomOption      ShineRandomOption
	RandomOptionCount ShineRandomOptionCount
}

// TODO: make loader with routines and channels, its too slow right now
func LoadItemData(path string) (*ItemData, error) {
	var (
		itemData          = &ItemData{}
		itemInfo          = ShineItemInfo{}
		itemInfoServer    = ShineItemInfoServer{}
		gradeItemOptions  = ShineGradeItemOption{}
		randomOption      = ShineRandomOption{}
		randomOptionCount = ShineRandomOptionCount{}
	)

	err := Load(path+"/shn/ItemInfo.shn", &itemInfo)

	if err != nil {
		return itemData, err
	}

	err = Load(path+"/shn/ItemInfoServer.shn", &itemInfoServer)

	if err != nil {
		return itemData, err
	}

	err = Load(path+"/shn/GradeItemOption.shn", &gradeItemOptions)

	if err != nil {
		return itemData, err
	}

	err = Load(path+"/shn/RandomOption.shn", &randomOption)

	if err != nil {
		return itemData, err
	}

	err = Load(path+"/shn/RandomOptionCount.shn", &randomOptionCount)

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
