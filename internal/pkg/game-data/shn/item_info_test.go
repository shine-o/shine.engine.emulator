package shn

import (
	"testing"
)

func TestItemInfo(t *testing.T)  {
	var file ShineItemInfo
	err := Load(filePath + "/shn/ItemInfo.shn", &file)
	if err != nil {
		t.Error(err)
	}

	if len(file.ShineRow) == 0 || file.ShineRow == nil {
		t.Errorf("expected rows cound = %v, actual rows count= %v", file.RowsCount, len(file.ShineRow))
	}
}
