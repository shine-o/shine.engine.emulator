package data

import "testing"

func TestLoadMobData(t *testing.T) {
	data, err := LoadMonsterData(filesPath)

	if err != nil {
		t.Fatal(err)
	}

	if &data.MapRegens == nil {
		t.Fatal("value cannot be nil")
	}

	if &data.MobInfo == nil {
		t.Fatal("value cannot be nil")
	}

	if &data.MobInfoServer == nil {
		t.Fatal("value cannot be nil")
	}
}
