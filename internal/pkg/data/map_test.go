package data

import "testing"

func TestLoadMapData(t *testing.T) {
	data, err := LoadMapData(filesPath)

	if err != nil {
		t.Fatal(err)
	}

	if data.Maps == nil {
		t.Fatal("value cannot be nil")
	}
}