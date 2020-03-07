package lib

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/spf13/viper"
)

// Read packet data from segments
func handleSegments(ctx context.Context, segment <-chan []byte, xorOffset *uint16) {
	var (
		data   []byte
		offset int
	)
	offset = 0

	for {
		select {
		case <-ctx.Done():
			return
		case b := <-segment:
			data = append(data, b...)

			if offset > len(data) {
				break
			}

			if offset != len(data) {
				var skipBytes int
				var pLen int
				var pType string
				var pd []byte

				pLen, pType = packetBoundary(offset, data)

				if pType == "small" {
					skipBytes = 1
				} else {
					skipBytes = 3
				}

				nextOffset := offset + skipBytes + pLen
				if nextOffset > len(data) {
					break
				}

				pd = append(pd, data[offset+skipBytes:nextOffset]...)

				xorCipher(pd, xorOffset)
				pc := processPacket(pType, pLen, pd)

				//go pc.handle(ctx)

				log.Infof("Got one %v", pc.pcb.String())
				offset += skipBytes + pLen
			}
		}
	}
}

// read packet data
func processPacket(pType string, pLen int, packetData []byte) ProtocolCommand {
	var opCode, department, command uint16
	br := bytes.NewReader(packetData)
	binary.Read(br, binary.LittleEndian, &opCode)

	department = opCode >> 10
	command = opCode & 1023

	return ProtocolCommand{
		pcb: ProtocolCommandBase{
			packetType:    pType,
			length:        pLen,
			department:    department,
			command:       command,
			operationCode: opCode,
			data:          packetData[2:],
		},
	}
}

// find out if big or small packet
// return length and type
func packetBoundary(offset int, b []byte) (int, string) {
	if b[offset] == 0 {
		var pLen uint16
		var tempB []byte
		tempB = append(tempB, b[offset:]...)
		br := bytes.NewReader(tempB)
		br.ReadAt(tempB, 1)
		binary.Read(br, binary.LittleEndian, &pLen)
		return int(pLen), "big"
	} else {
		var pLen uint8
		pLen = b[offset]
		return int(pLen), "small"
	}
}

// decrypt encrypted bytes using captured xorKey and xorTable
func xorCipher(eb []byte, xorPos *uint16) {
	xorLimit := uint16(viper.GetInt("crypt.xorLimit"))
	for i, _ := range eb {
		eb[i] ^= xorKey[*xorPos]
		*xorPos++
		//log.Info(*xorPos)
		if *xorPos >= xorLimit {
			*xorPos = 0
		}
	}
}
