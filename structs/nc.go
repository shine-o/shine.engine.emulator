package structs

import (
	"encoding/binary"
	"github.com/google/logger"
	"gopkg.in/restruct.v1"
	"io/ioutil"
)

type Aux interface {
	String()    string
	PdbAnalog() string
}

type NC interface {
	Aux
	Pack() ([]byte, error)
	Unpack([]byte) error
}

var log *logger.Logger

func init() {
	log = logger.Init("Structs Logger", true, false, ioutil.Discard)
	log.Info("structs logger init()")
}

// ReadBinary data into a given struct
// if struct size is bigger than available data, fill with zeros
func Unpack(data []byte, nc interface{}) error {
	var buffer []byte
	structSize , err := restruct.SizeOf(nc)

	if err != nil {
		log.Error(err)
		return err
	}

	if structSize == -1 {
		buffer = make([]byte, 0)
	} else {
		buffer = make([]byte, structSize)
	}
	copy(buffer, data)

	err = restruct.Unpack(buffer, binary.LittleEndian, nc)

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