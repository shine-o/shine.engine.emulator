package networking

import (
	"bytes"
	"encoding/binary"
	"fmt"
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
		log.Error(err)
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

	if (&PCList{}) != commandList {// should be commented out on production to increase performance
		commandList.mu.Lock()
		if dpt, ok := commandList.Departments[uint8(department)]; ok {
			pc.Base.ClientStructName = dpt.ProcessedCommands[fmt.Sprintf("%X", command)]
		} else {
			log.Warningf("Missing friendly name for command with: operationCode %v,  department %v, command %v, ", opCode, department, fmt.Sprintf("%X", command))
		}
		commandList.mu.Unlock()
	}

	return pc, nil
}
