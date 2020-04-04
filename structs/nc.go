package structs

import (
	"encoding/binary"
	"github.com/google/logger"
	"gopkg.in/restruct.v1"
	"io/ioutil"
)

type NC interface {
	String() string
	PdbType() string
}

var log *logger.Logger

func init() {
	log = logger.Init("Structs Logger", true, false, ioutil.Discard)
	log.Info("structs logger init()")
}

// ReadBinary data into a given struct
// if struct size is bigger than available data, fill with zeros
func Unpack(data []byte, nc interface{}) error {
	err := restruct.Unpack(data, binary.LittleEndian, nc)

	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// WriteBinary data into a given struct and return bytes
// todo: check again if it can be done like this: https://github.com/go-restruct/restruct/issues/39
func Pack(nc interface{}) ([]byte, error) {
	data, err := restruct.Pack(binary.LittleEndian, nc)
	if err != nil {
		log.Error(err)
		return data, err
	}
	return data, nil
}
