package networking

import (
	"bytes"
	"encoding/binary"
)

func ReadBinary(data []byte, nc interface{}) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.LittleEndian, nc); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func WriteBinary(nc interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, nc); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return buf.Bytes(), nil
}
