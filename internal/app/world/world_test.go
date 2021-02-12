package world

import (
	"encoding/hex"
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
			NcStruct: &structs.NcMiscSeedAck{},
		},
		networking.NC_USER_LOGINWORLD_REQ: {
			NcStruct: &structs.NcUserLoginWorldReq{},
		},
		networking.NC_MISC_GAMETIME_REQ:{
			NcStruct: &structs.NcMiscGameTimeAck{},
		},
		networking.NC_USER_WILL_WORLD_SELECT_REQ:{
			NcStruct: &structs.NcUserWorldSelectReq{},
		},
		networking.NC_AVATAR_CREATE_REQ:{
			NcStruct: &structs.NcAvatarCreateReq{},
		},
		networking.NC_AVATAR_ERASE_REQ:{
			NcStruct: &structs.NcAvatarEraseReq{},
		},
		networking.NC_CHAR_LOGIN_REQ:{
			NcStruct: &structs.NcCharLoginReq{},
		},
		networking.NC_CHAR_OPTION_GET_WINDOWPOS_REQ:{
			NcStruct: &structs.NcCharOptionWindowPos{},
		},
		//networking.NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ:{
		//	NcStruct: &structs.NcCharOptionGetShortcutSizeReq{},
		//},
		//networking.NC_PRISON_GET_REQ:{
		//	NcStruct: &structs.NcPrisonG
		//},
		networking.NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ:{
			NcStruct: &structs.NcCharOptionSetShortcutDataReq{},
		},

		//networking.NC_USER_AVATAR_LIST_REQ:{
		//	NcStruct: &structs.NcUserA
		//},

		networking.NC_USER_LOGINWORLD_ACK: {
			NcStruct: &structs.NcUserLoginWorldAck{},
		},
		networking.NC_CHAR_LOGIN_ACK: {
			NcStruct: &structs.NcCharLoginAck{},
		},
		networking.NC_AVATAR_CREATESUCC_ACK: {
			NcStruct: &structs.NcAvatarCreateSuccAck{},
		},
		networking.NC_AVATAR_ERASESUCC_ACK: {
			NcStruct: &structs.NcAvatarEraseSuccAck{},
		},
		networking.NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD: {
			NcStruct: &structs.NcCharOptionImproveGetGameOptionCmd{},
		},
		networking.NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD: {
			NcStruct: &structs.NcCharGetKeyMapCmd{},
		},
		networking.NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD: {
			NcStruct: &structs.NcCharGetShortcutDataCmd{},
		},
		networking.NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_ACK: {
			NcStruct: &structs.NcCharOptionImproveShortcutDataAck{},
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