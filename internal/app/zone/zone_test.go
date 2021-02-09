package zone

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"io/ioutil"
	"reflect"
	"testing"
)

type packet struct {
	ncStruct interface{}
	assert func(interface{}) error
}

var targetPackets = make(map[uint16]packet)

var packetData = make(map[uint16][]byte)

func setNetPackets()  {
	targetPackets[networking.NC_MISC_SEED_ACK] = packet{
		ncStruct:   structs.NcMiscSeedAck{},
		assert: func(i interface{}) error {
			ncS, ok := i.(structs.NcMiscSeedAck)
			if !ok {
				return errors.New(fmt.Sprintf("bad struct %v", reflect.TypeOf(ncS).String()))
			}
			if ncS.Seed > 499 {
				return errors.New(fmt.Sprintf("bad seed key %v", ncS.Seed))
			}
			return nil
		},
	}
}

func loadPacketData(filePath string) {
	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(data, &packetData)

	if err != nil {
		log.Fatal(err)
	}
}

func testPackets(t * testing.T)  {
	for opCode, data := range packetData {
		p, ok :=  targetPackets[opCode]
		if ok {
			err := testPacket(p, data)
			if err != nil {
				t.Error(err)
			}
		}
	}
}

func testPacket(p packet, data []byte) error {
	// assert unpacking goes well
	err := structs.Unpack(data, &p.ncStruct)
	if err != nil {
		return err
	}
	err = p.assert(p.ncStruct)
	if err != nil {
		return err
	}
	return nil
	// assert the assigned function tests well
}