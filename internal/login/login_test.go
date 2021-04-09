package login

import (
	"encoding/hex"
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/utils"
	"reflect"
	"testing"
)

func netPackets() utils.TargetPackets {
	tp := utils.TargetPackets{
		networking.NC_MISC_SEED_ACK: {
			NcStruct: &structs.NcMiscSeedAck{},
		},
		networking.NC_USER_CLIENT_VERSION_CHECK_REQ: {
			NcStruct: &structs.NcUserClientVersionCheckReq{},
		},
		networking.NC_USER_US_LOGIN_REQ: {
			NcStruct: &structs.NewUserLoginReq{},
		},
		networking.NC_USER_WORLDSELECT_REQ: {
			NcStruct: &structs.NcUserWorldSelectReq{},
		},
		networking.NC_USER_NORMALLOGOUT_CMD:   {},
		networking.NC_USER_LOGIN_WITH_OTP_REQ: {},

		networking.NC_USER_LOGIN_ACK: {
			NcStruct: &structs.NcUserLoginAck{},
		},

		networking.NC_USER_LOGINFAIL_ACK: {
			NcStruct: &structs.NcUserLoginFailAck{},
		},

		networking.NC_USER_WORLDSELECT_ACK: {
			NcStruct: &structs.NcUserWorldSelectAck{},
		},
	}

	return tp

}

func TestPackets(t *testing.T) {
	netPackets := netPackets()

	files := []string{
		"../../../test-data/packets-1612910284-version-1.02.295.json",
		"../../../test-data/packets-1613170127-version-1.02.295.json",
	}

	for _, f := range files {
		packetData := utils.LoadPacketData(f)
		for opCode, packet := range netPackets {
			dataStrings, ok := packetData[uint16(opCode)]
			if ok {
				for _, dataString := range dataStrings {
					if dataString == "" {
						continue
					}

					data, err := hex.DecodeString(dataString)
					if err != nil {
						t.Error(err)
					}
					err = utils.TestPacket(packet, data)
					if err != nil {
						t.Error(err)
					}
					t.Log(fmt.Sprintf("ok, struct=%v data=%v", reflect.TypeOf(packet.NcStruct).String(), dataString))
				}
			}
		}
	}
}
