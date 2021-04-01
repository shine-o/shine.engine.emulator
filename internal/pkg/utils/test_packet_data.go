package utils

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"io/ioutil"
	"log"
	"reflect"
)

type TestablePacket struct {
	NcStruct interface{}
	Assert   func(interface{}) error
}

type TargetPackets map[networking.OperationCode]TestablePacket

func LoadPacketData(filePath string) map[uint16][]string {
	var packetData = make(map[uint16][]string)

	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &packetData)

	if err != nil {
		log.Fatal(err)
	}

	return packetData
}

func TestPacket(p TestablePacket, data []byte) error {
	// assert unpacking goes well
	err := structs.Unpack(data, p.NcStruct)
	if err != nil {
		return errors.New(fmt.Sprintf("packet=%v, data=%v", reflect.TypeOf(p.NcStruct).String(), hex.EncodeToString(data)))
	}
	if p.Assert != nil {
		err = p.Assert(p.NcStruct)
		if err != nil {
			return err
		}
	}
	return nil
	// assert the assigned function tests well
}
