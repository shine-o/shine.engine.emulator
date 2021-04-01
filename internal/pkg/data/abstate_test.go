package data

import "testing"

// load data path
//

func TestAbstateFile(t *testing.T) {
	var file ShineAbState
	err := Load("AbState.shn", &file)
	if err != nil {
		t.Error(err)
	}

	if len(file.ShineRow) == 0 || file.ShineRow == nil {
		t.Errorf("expected rows cound = %v, actual rows count= %v", file.RowsCount, len(file.ShineRow))
	}
}
