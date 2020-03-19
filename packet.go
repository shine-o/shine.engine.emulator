package networking

import (
	"bytes"
	"encoding/binary"
	"github.com/google/logger"
)

var (
	log         *logger.Logger
	xorKey      []byte
	xorLimit    uint16
	commandList *PCList
)

// PacketBoundary of a packet in a data segment
// Shine packets can be small or big
// small packets specify the length in the first byte
// big packets ignore the first byte (that is, is 0 value), and set the length in the next 2 bytes
// the data segment is usually a local variable in a goroutine
func PacketBoundary(offset int, data []byte) (int, string) {
	if data[offset] == 0 { // big packet
		var (
			pLen uint16
			d    []byte
		)

		d = append(d, data[offset+1:]...)
		br := bytes.NewReader(d)

		if err := binary.Read(br, binary.LittleEndian, &pLen); err != nil {
			log.Error(err)
		}

		return int(pLen), "big"

	}
	var pLen uint8
	pLen = data[offset]
	return int(pLen), "small"
}

// DecodePacket data into a struct
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
