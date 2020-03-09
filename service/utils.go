package service

import (
	"bytes"
	"encoding/binary"
	protocol "github.com/shine-o/shine.engine.protocol"
)

func readBinary(data []byte, nc interface{}) error {
	buf := bytes.NewBuffer(data)
	if err := binary.Read(buf, binary.LittleEndian, nc); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func writeBinary(nc interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, nc); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func logOutboundPacket(pc *protocol.Command) { // should be used in only one place :(
	log.Infof("Outbound packet %v", pc.Base.String())
}
