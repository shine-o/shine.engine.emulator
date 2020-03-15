package networking

import (
	"bytes"
	"encoding/binary"
)

func ReadBinary(data []byte, nc interface{}) error {
	structSize := binary.Size(nc)
	buffer := make([]byte, structSize)
	copy(buffer, data)
	buf := bytes.NewBuffer(buffer)
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