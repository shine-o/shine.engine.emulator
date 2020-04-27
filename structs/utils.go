package structs

import (
	"encoding/binary"
	"fmt"
	"github.com/go-restruct/restruct"
	"github.com/google/logger"
	"io/ioutil"
	"reflect"
)

type NC interface {
	//String() string
	//PdbType() string
}

var log *logger.Logger

func init() {
	log = logger.Init("Structs Logger", true, false, ioutil.Discard)
	log.Info("structs logger init()")
	restruct.EnableExprBeta()
}

// ReadBinary data into a given struct
// if struct size is bigger than available data, fill with zeros
func Unpack(data []byte, nc interface{}) error {
	err := restruct.Unpack(data, binary.LittleEndian, nc)
	size, _ := restruct.SizeOf(nc)

	if err != nil {
		log.Errorf("size of %v: %v, size of data %v", reflect.TypeOf(nc).String(), size, len(data))
		log.Error(err)
		return err
	}

	if size != len(data) {
		return fmt.Errorf("size of %v: %v, size of data %v", reflect.TypeOf(nc).String(), size, len(data))
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
