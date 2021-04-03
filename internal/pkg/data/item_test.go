package data

import "testing"

func TestLoadItemData(t *testing.T) {
	data, err := LoadItemData(filesPath)

	if err != nil {
		t.Fatal(err)
	}

	if &data.ItemInfo == nil {
		t.Fatal("value cannot be nil")
	}

	if &data.ItemInfoServer == nil {
		t.Fatal("value cannot be nil")
	}

	if &data.GradeItemOptions == nil {
		t.Fatal("value cannot be nil")
	}

	if &data.RandomOption == nil {
		t.Fatal("value cannot be nil")
	}

	if &data.RandomOptionCount == nil {
		t.Fatal("value cannot be nil")
	}

	//if data.ItemDropGroup == nil {
	//	t.Fail()
	//}
	//
	//if data.ItemDropTable == nil {
	//	t.Fail()
	//}
}