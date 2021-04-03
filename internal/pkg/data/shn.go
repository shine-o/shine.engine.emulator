package data

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"io/ioutil"
)

type ShineDataType uint32

const (
	TypeListEnd ShineDataType = iota
	TypeListByte
	TypeListWord
	TypeListDword
	TypeListQWord
	TypeListFloat
	TypeListFilename
	TypeListFileAuto
	TypeListRemark
	TypeListStr
	TypeListStrAuto
	TypeListInx
	TypeListInxByte
	TypeListInxWord
	TypeListInxDword
	TypeListInxQWord
	TypeListByteBit
	TypeListWordBit
	TypeListDwordBit
	TypeListQWordBit
	TypeListByteArray
	TypeListWordArray
	TypeListDWordArray
	TypeListQWordArray
	TypeListStrArray
	TypeListStrAutoArray
	TypeListVarStr
	TypeListInxStr
	TypeListUnknownEd
	TypeListTwoInx
)

type ShineDataMode uint32

const (
	DataModeNormal ShineDataMode = iota
	DataModeEncryption
)

type ShineRawFile struct {
	VersionKey uint32
	Version    [20]byte
	Reserved   uint32
	DataMode   ShineDataMode
	FileSize   int32
	// column, row data
	Data []byte `struct-while:"!_eof"`
}

type ShineColumn struct {
	Name string `struct:"[48]byte"`
	Type ShineDataType
	Size uint32
}


func Load(filesPath string, shn interface{}) error {
	data, err := loadRawData(filesPath)
	if err != nil {
		return err
	}
	err = structs.Unpack(data, shn)
	if err != nil {
		return err
	}
	return nil
}


func loadRawData(filesPath string) ([]byte, error) {
	var srf ShineRawFile
	data, err := ioutil.ReadFile(filesPath)

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
