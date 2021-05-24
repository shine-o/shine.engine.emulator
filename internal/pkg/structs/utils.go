package structs

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"reflect"

	"github.com/go-restruct/restruct"
	"github.com/google/logger"
)

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
	if err != nil {
		size, _ := restruct.SizeOf(nc)
		return fmt.Errorf("size of %v: %v, size of data %v, %v", reflect.TypeOf(nc).String(), size, len(data), err)
	}
	return nil
}

// WriteBinary data into a given struct and return bytes
func Pack(nc interface{}) ([]byte, error) {
	data, err := restruct.Pack(binary.LittleEndian, nc)
	if err != nil {
		return data, err
	}
	return data, nil
}
