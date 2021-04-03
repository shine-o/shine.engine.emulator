package data

import "testing"

func TestLoadNpcData(t *testing.T) {
	data, err := LoadNPCData(filesPath)

	if err != nil {
		t.Fatal(err)
	}

	if data.MapNPCs == nil {
		t.Fatal("value cannot be nil")
	}

	//if data.VendorNPCs == nil {
	//	t.Fail()
	//}
}

