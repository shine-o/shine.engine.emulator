package data

import "testing"

const nilValue = "value cannot be nil"

func TestLoadItemData(t *testing.T) {
	data, err := LoadItemData(filesPath)
	if err != nil {
		t.Fatal(err)
	}

	if &data.ItemInfo == nil {
		t.Fatal(nilValue)
	}

	if &data.ItemInfoServer == nil {
		t.Fatal(nilValue)
	}

	if &data.GradeItemOptions == nil {
		t.Fatal(nilValue)
	}

	if &data.RandomOption == nil {
		t.Fatal(nilValue)
	}

	if &data.RandomOptionCount == nil {
		t.Fatal(nilValue)
	}

	//if data.ItemDropGroup == nil {
	//	t.Fail()
	//}
	//
	//if data.ItemDropTable == nil {
	//	t.Fail()
	//}
}
