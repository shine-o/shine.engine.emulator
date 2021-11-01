package data

import (
	"encoding/csv"
	"os"
	"path/filepath"
)

// check if path is correct and return absolute path
func ValidPath(path string) (string, error) {
	var absPath string
	absPath, err := filepath.Abs(path)
	if err != nil {
		return absPath, err
	}
	if _, err := os.Stat(path); err == os.ErrNotExist {
		return absPath, err
	}
	return absPath, nil
}

func loadTxtFile(filesPath string) ([][]string, error) {
	var data [][]string
	txtFile, err := os.Open(filesPath)
	if err != nil {
		return data, err
	}
	reader := csv.NewReader(txtFile)

	reader.Comma = '\t'
	reader.FieldsPerRecord = -1

	data, err = reader.ReadAll()
	if err != nil {
		return data, err
	}
	return data, err
}
