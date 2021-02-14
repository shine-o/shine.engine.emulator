package zone

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/utils"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"reflect"
	"testing"
)


func netPackets() utils.TargetPackets  {
	tp := utils.TargetPackets{
		networking.NC_MISC_SEED_ACK: {
			NcStruct:   &structs.NcMiscSeedAck{},
			Assert: func(i interface{}) error {
				ncS, ok := i.(*structs.NcMiscSeedAck)
				if !ok {
					return errors.New(fmt.Sprintf("bad struct %v", reflect.TypeOf(ncS).String()))
				}
				if ncS.Seed > 499 {
					return errors.New(fmt.Sprintf("bad seed key %v", ncS.Seed))
				}
				return nil
			},
		},

		// inbound
		networking.NC_MISC_HEARTBEAT_ACK: {
			NcStruct: &structs.NcMiscHeartBeatAck{},
		},
		networking.NC_MAP_LOGIN_REQ: {
			NcStruct: &structs.NcMapLoginReq{},
		},
		networking.NC_MAP_LOGINCOMPLETE_CMD: {
			NcStruct: &structs.NcMapLoginCompleteCmd{},
		},
		//networking.NC_CHAR_LOGOUTREADY_CMD: {
		//	NcStruct: &structs.NcLogou
		//},
		//networking.NC_CHAR_LOGOUTCANCEL_CMD: {
		//	NcStruct: &structs.NcCharLogou
		//},
		networking.NC_ACT_MOVEWALK_CMD: {
			NcStruct: &structs.NcActMoveWalkCmd{},
		},
		networking.NC_ACT_MOVERUN_CMD: {
			NcStruct: &structs.NcActMoveRunCmd{},
		},
		//networking.NC_ACT_JUMP_CMD: {
		//},
		networking.NC_ACT_STOP_REQ: {
			NcStruct: &structs.NcActStopReq{},
		},
		networking.NC_BRIEFINFO_INFORM_CMD: {
			NcStruct: &structs.NcBriefInfoInformCmd{},
		},
		//networking.NC_BAT_TARGETTING_REQ: {
		//},
		networking.NC_BAT_UNTARGET_REQ: {
			NcStruct: &structs.NcBatUnTargetReq{},
		},
		//networking.NC_USER_NORMALLOGOUT_CMD: {
		//	NcStruct: &structs.N
		//},
		networking.NC_ACT_NPCCLICK_CMD: {
			NcStruct: &structs.NcActNpcClickCmd{},
		},
		networking.NC_MENU_SERVERMENU_ACK: {
			NcStruct: &structs.NcServerMenuAck{},
		},

		// outbound
		networking.NC_MENU_SERVERMENU_REQ: {
			NcStruct: &structs.NcServerMenuReq{},
		},
		networking.NC_BAT_TARGETINFO_CMD: {
			NcStruct: &structs.NcBatTargetInfoCmd{},
		},
		networking.NC_ACT_SOMEONEMOVEWALK_CMD: {
			NcStruct: &structs.NcActSomeoneMoveWalkCmd{},
		},
		networking.NC_ACT_SOMEONEMOVERUN_CMD: {
			NcStruct: &structs.NcActSomeoneMoveRunCmd{},
		},
		networking.NC_BRIEFINFO_BRIEFINFODELETE_CMD: {
			NcStruct: &structs.NcBriefInfoDeleteCmd{},
		},
		networking.NC_ACT_SOMEONESTOP_CMD: {
			NcStruct: &structs.NcActSomeoneStopCmd{},
		},
		networking.NC_ACT_SOMEEONEJUMP_CMD: {
			NcStruct: &structs.NcActSomeoneJumpCmd{},
		},
		networking.NC_BRIEFINFO_LOGINCHARACTER_CMD: {
			NcStruct: &structs.NcBriefInfoLoginCharacterCmd{},
		},
		networking.NC_BRIEFINFO_REGENMOB_CMD: {
			NcStruct: &structs.NcBriefInfoRegenMobCmd{},
		},
		networking.NC_BRIEFINFO_CHARACTER_CMD: {
			NcStruct: &structs.NcBriefInfoCharacterCmd{},
		},
		networking.NC_BRIEFINFO_MOB_CMD: {
			NcStruct: &structs.NcBriefInfoMobCmd{},
		},
		networking.NC_MAP_LINKSAME_CMD: {
			NcStruct: &structs.NcMapLinkSameCmd{},
		},
		networking.NC_CHAR_CLIENT_BASE_CMD: {
			NcStruct: &structs.NcCharClientBaseCmd{},
		},
		networking.NC_CHAR_CLIENT_SHAPE_CMD: {
			NcStruct: &structs.NcCharClientShapeCmd{},
		},
		networking.NC_MAP_LOGIN_ACK: {
			NcStruct: &structs.NcMapLoginAck{},
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
			dataStrings, ok :=  packetData[uint16(opCode)]
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
