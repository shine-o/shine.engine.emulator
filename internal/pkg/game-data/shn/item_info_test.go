package shn

import "testing"


var filePath = "../../../../data"

func TestItemInfo(t *testing.T)  {
	var file ShineItemInfo
	err := Load(filePath + "/ItemInfo.shn", &file)
	if err != nil {
		t.Error(err)
	}

	if len(file.ShineRow) == 0 || file.ShineRow == nil {
		t.Errorf("expected rows cound = %v, actual rows count= %v", file.RowsCount, len(file.ShineRow))
	}
}

func TestLinkedIndexes(t *testing.T) {
	var file ShineItemInfo
	err := Load(filePath + "/ItemInfo.shn", &file)
	if err != nil {
		t.Error(err)
	}

	res, err := file.MissingIDs(filePath)

	if err != nil {
		t.Error(err)
	}

	count := 0
	if len(res) > 0 {
		for k, v1 := range res {
			for _, v2 := range v1 {
				t.Logf("%v not in %v \n", v2, k)
				count++
			}
		}
		t.Errorf("missing IDs, total %v", count)
	}

}
