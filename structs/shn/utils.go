package shn

import (
	"github.com/shine-o/shine.engine.core/structs"
	"io/ioutil"
)

func Load(filePath string, shn interface{}) error {
	data, err := LoadRawData(filePath)
	if err != nil {
		return err
	}
	err = structs.Unpack(data, shn)
	if err != nil {
		return err
	}
	return nil
}

func LoadRawData(filePath string) ([]byte, error) {
	var srf ShineRawFile
	var content []byte
	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		return content, err
	}

	err = structs.Unpack(data, &srf)

	if err != nil {
		return content, err
	}

	content = make([]byte, len(data[36:]))

	copy(content, data[36:])

	DecryptSHN(content, len(content))

	return content, err
}

func DecryptSHN(data []byte, length int) {
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
