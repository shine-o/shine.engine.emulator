package networking

import (
	"bytes"
	"encoding/binary"
)

var (
	xorKey      []byte
	xorLimit    uint16
	commandList *PCList
)

// PacketBoundary of a packet in a data segment
// Shine packets can be small or big
// small packets specify the length in the first byte
// big packets ignore the first byte (that is, is 0 value), and set the length in the next 2 bytes
// the data segment is usually a local variable in a goroutine
func PacketBoundary(offset int, data []byte) (uint16, int) {
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
		return pLen, 3
	}
	var pLen uint8
	pLen = data[offset]
	return uint16(pLen), 1
}

// todo: cleaner way to create packets
//type RawPacket struct {
//	OpCode uint16
//	Data []byte `struct-while:"!_eof"`
//}

// DecodePacket data into a struct
func DecodePacket(data []byte) (Command, error) {

	var (
		opCode     uint16
		department uint16
		command    uint16
		pc         Command
	)

	br := bytes.NewReader(data)

	if err := binary.Read(br, binary.LittleEndian, &opCode); err != nil {
		log.Errorf("data stream %v %v", data, err)
		return pc, err
	}

	department = opCode >> 10
	command = opCode & 1023

	pc.Base = CommandBase{
		Department:    department,
		Command:       command,
		OperationCode: opCode,
		Data:          data[2:], // omit operationCode bytes
	}

	return pc, nil
}
