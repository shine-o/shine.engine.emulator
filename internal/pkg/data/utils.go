package data

import (
	"encoding/csv"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"io/ioutil"
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

func loadTxtFile(filePath string) ([][]string, error) {
	var data [][]string
	txtFile, err := os.Open(filePath)
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


func Load(filePath string, shn interface{}) error {
	data, err := loadRawData(filePath)
	if err != nil {
		return err
	}
	err = structs.Unpack(data, shn)
	if err != nil {
		return err
	}
	return nil
}

func loadRawData(filePath string) ([]byte, error) {
	var srf ShineRawFile
	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		return srf.Data, err
	}

	err = structs.Unpack(data, &srf)

	if err != nil {
		return srf.Data, err
	}

	decryptSHN(srf.Data, int(srf.FileSize)-36)

	return srf.Data, err
}

func decryptSHN(data []byte, length int) {
	if length < 1 {
		return
	}
	l := byte(length)
	for i := length - 1; i >= 0; i-- {
		var nl byte
		data[i] = data[i] ^ l
		nl = byte(i)
		nl = nl & byte(15)
		nl = nl + byte(85)
		nl = nl ^ (byte(i) * byte(11))
		nl = nl ^ l
		nl = nl ^ byte(170)
		l = nl
	}
}

