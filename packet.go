package networking

import (
	"bytes"
	"encoding/binary"
	"github.com/google/logger"
	"io/ioutil"
	dl "log"
)

var (
	log         *logger.Logger
	xorKey      []byte
	xorLimit    uint16
	commandList *PCList
)

func init() {
	logger.SetFlags(dl.Llongfile)

	log = logger.Init("NetworkingLogger", true, true, ioutil.Discard)

	log.Info("Networking Logger init()")
}

// find out if big or small packet
// return length and type
func PacketBoundary(offset int, data []byte) (int, string) {
	if data[offset] == 0 {
		var (
			pLen uint16
			d    []byte
		)

		d = append(d, data[offset:]...)

		br := bytes.NewReader(d)

		if _, err := br.ReadAt(d, 1); err != nil {
			log.Error(err)
		}

		if err := binary.Read(br, binary.LittleEndian, &pLen); err != nil {
			log.Error(err)
		}

		return int(pLen), "big"
	} else {
		var pLen uint8
		pLen = data[offset]
		return int(pLen), "small"
	}
}

// read packet data
func DecodePacket(pType string, pLen int, data []byte) (Command, error) {
	var (
		opCode     uint16
		department uint16
		command    uint16
		pc         Command
	)

	br := bytes.NewReader(data)

	if err := binary.Read(br, binary.LittleEndian, &opCode); err != nil {
		log.Error(err)
		return pc, err
	}

	department = opCode >> 10
	command = opCode & 1023

	pc.Base = CommandBase{
		PacketType:    pType,
		Length:        pLen,
		Department:    department,
		Command:       command,
		OperationCode: opCode,
		Data:          data[2:], // omit operationCode bytes
	}

	return pc, nil
}
