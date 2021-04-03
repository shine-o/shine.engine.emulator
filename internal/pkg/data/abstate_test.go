package data

import "testing"

func TestAbstateFile(t *testing.T) {
	var file ShineAbState
	err := Load(filesPath+"/shn/AbState.shn", &file)
	if err != nil {
		t.Error(err)
	}

	if len(file.ShineRow) == 0 || file.ShineRow == nil {
		t.Errorf("expected rows cound = %v, actual rows count= %v", file.RowsCount, len(file.ShineRow))
	}
}
