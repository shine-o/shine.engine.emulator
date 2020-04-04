package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.networking/structs"
	"gopkg.in/restruct.v1"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sync"
)

var ocs *OpCodeStructs

type OpCodeStructs struct {
	structs map[uint16]string
	mu      sync.Mutex
}

type NcRepresentation struct {
	Pdb          string `json:"pdb_analog"`
	UnpackedData string `json:"unpacked_data"`
}

func generateOpCodeSwitch() {
	type processedStructs struct {
		List map[uint16]bool `json:"processedStructs"`
	}
	// load processed structs
	filePath, err := filepath.Abs("defaults/processed-structs.json")
	if err != nil {
		log.Fatal(err)
	}

	d, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	var ps processedStructs
	err = json.Unmarshal(d, &ps)
	if err != nil {
		log.Fatal(err)
	}
	log.Error(ps)
	start := `
	package generated
	func ncStructRepresentation(opCode uint16, data []byte) {
	switch opCode {` + "\n"
	ocs.mu.Lock()
	fmt.Println(ocs.structs)
	for k, v := range ocs.structs {
		if _, processed := ps.List[k]; processed {
			continue
		}
		caseStmt := fmt.Sprintf("\t"+`case %v:`+"\n", k)
		caseStmt += fmt.Sprintf("\t"+"// %v\n", v)
		caseStmt += "\t" + "// return ncStructData(&nc, data)\n"
		caseStmt += "\t" + "break\n"
		start += caseStmt
	}
	ocs.mu.Unlock()
	end := "}}"
	//log.Error(start+end)

	pathName, err := filepath.Abs("output/opcodes-switch.go")
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(pathName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Write([]byte(start + end))
	if err != nil {
		log.Fatal(err)
	}
}

func ncStructRepresentation(opCode uint16, data []byte) (NcRepresentation, error) {
	switch opCode {
	case 3173:
		nc := structs.NcUserClientVersionCheckReq{}
		return ncStructData(&nc, data)
	case 3084:
		nc := structs.NcUserWorldSelectAck{}
		return ncStructData(&nc, data)
	case 3162:
		nc := structs.NcUserLoginWorldReq{}
		return ncStructData(&nc, data)
	case 4153:
		// NC_CHAR_CLIENT_SHAPE_CMD
		// return ncStructData(&nc, data)
		break
	case 4157:
		// NC_CHAR_CLIENT_SKILL_CMD
		nc := structs.NcCharClientSkillCmd{}
		return ncStructData(&nc, data)
	case 4168:
		// NC_CHAR_CLIENT_GAME_CMD
		// return ncStructData(&nc, data)
		break
	case 8193:
		// NC_ACT_CHAT_REQ
		nc := structs.NcActChatReq{}
		return ncStructData(&nc, data)
	case 2055:
		// NC_MISC_SEED_ACK
		nc := structs.NcMiscSeedAck{}
		return ncStructData(&nc, data)
	case 4154:
		// NC_CHAR_CLIENT_QUEST_DOING_CMD
		// return ncStructData(&nc, data)
		break
	case 4302:
		// NC_CHAR_CLIENT_QUEST_READ_CMD
		// return ncStructData(&nc, data)
		break
	case 4311:
		// NC_CHAR_CLIENT_QUEST_REPEAT_CMD
		// return ncStructData(&nc, data)
		break
	case 4167:
		// NC_CHAR_CLIENT_ITEM_CMD
		nc := structs.NcCharClientItemCmd{}
		return ncStructData(&nc, data)
	case 6147:
		// NC_MAP_LOGINCOMPLETE_CMD
		nc := structs.NcMapLoginCompleteCmd{}
		return ncStructData(&nc, data)
		break
	case 8203:
		// NC_ACT_ENDOFTRADE_CMD
		// return ncStructData(&nc, data)
		break
	case 4318:
		// NC_CHAR_CLIENT_COININFO_CMD
		// return ncStructData(&nc, data)
		break
	case 4155:
		// NC_CHAR_CLIENT_QUEST_DONE_CMD
		// return ncStructData(&nc, data)
		break
	case 4158:
		// NC_CHAR_CLIENT_PASSIVE_CMD
		// return ncStructData(&nc, data)
		break
	case 6145:
		// NC_MAP_LOGIN_REQ
		nc := structs.NcMapLoginReq{}
		return ncStructData(&nc, data)
	case 4170:
		// NC_CHAR_CLIENT_CHARGEDBUFF_CMD
		// return ncStructData(&nc, data)
		break
	case 17438:
		// NC_QUEST_RESET_TIME_CLIENT_CMD
		// return ncStructData(&nc, data)
		break
	case 6146:
		// NC_MAP_LOGIN_ACK
		// return ncStructData(&nc, data)
		break
	case 4152:
		// NC_CHAR_CLIENT_BASE_CMD
		// return ncStructData(&nc, data)
		break
	case 4169:
		// NC_CHAR_CLIENT_CHARTITLE_CMD
		nc := structs.NcClientCharTitleCmd{}
		return ncStructData(&nc, data)
	case 7192:
		// NC_BRIEFINFO_ABSTATE_CHANGE_CMD
		// return ncStructData(&nc, data)
		break
	case 12303:
		// NC_ITEM_EQUIP_REQ
		// return ncStructData(&nc, data)
		break
	case 20487:
		// NC_SOULSTONE_HP_USE_REQ
		// return ncStructData(&nc, data)
		break
	case 7182:
		// NC_BRIEFINFO_BRIEFINFODELETE_CMD
		// return ncStructData(&nc, data)
		break
	case 28676:
		// NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ
		// return ncStructData(&nc, data)
		break
	case 2053:
		// NC_MISC_HEARTBEAT_ACK
		// return ncStructData(&nc, data)
		break
	case 12296:
		// NC_ITEM_DROP_ACK
		// return ncStructData(&nc, data)
		break
	case 3175:
		// NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK
		// return ncStructData(&nc, data)
		break
	case 28677:
		// NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK
		// return ncStructData(&nc, data)
		break
	case 4324:
		// NC_CHAR_CLIENT_CARDCOLLECT_CMD
		// return ncStructData(&nc, data)
		break
	case 12297:
		// NC_ITEM_PICK_REQ
		// return ncStructData(&nc, data)
		break
	case 4294:
		// NC_CHAR_ADMIN_LEVEL_INFORM_CMD
		// return ncStructData(&nc, data)
		break
	case 4114:
		// NC_CHAR_GUILD_CMD
		// return ncStructData(&nc, data)
		break
	case 37908:
		// NC_HOLY_PROMISE_LIST_CMD
		// return ncStructData(&nc, data)
		break
	case 4097:
		// NC_CHAR_LOGIN_REQ
		// return ncStructData(&nc, data)
		break
	case 12295:
		// NC_ITEM_DROP_REQ
		// return ncStructData(&nc, data)
		break
	case 7178:
		// NC_BRIEFINFO_DROPEDITEM_CMD
		// return ncStructData(&nc, data)
		break
	case 7170:
		// NC_BRIEFINFO_CHANGEDECORATE_CMD
		// return ncStructData(&nc, data)
		break
	case 12321:
		// NC_ITEM_CHARGEDINVENOPEN_ACK
		// return ncStructData(&nc, data)
		break
	case 9256:
		// NC_BAT_ABSTATERESET_CMD
		// return ncStructData(&nc, data)
		break
	case 12332:
		// NC_ITEM_REWARDINVENOPEN_REQ
		// return ncStructData(&nc, data)
		break
	case 31751:
		// NC_PRISON_GET_ACK
		// return ncStructData(&nc, data)
		break
	case 18476:
		// NC_SKILL_ITEMACTIONCOOLTIME_CMD
		// return ncStructData(&nc, data)
		break
	case 4308:
		// NC_CHAR_MYSTERYVAULT_UI_STATE_CMD
		// return ncStructData(&nc, data)
		break
	case 15361:
		// NC_MENU_SERVERMENU_REQ
		// return ncStructData(&nc, data)
		break
	case 6149:
		// NC_MAP_LOGOUT_CMD
		// return ncStructData(&nc, data)
		break
	case 3092:
		// NC_USER_LOGINWORLD_ACK
		// return ncStructData(&nc, data)
		break
	case 4387:
		// NC_CHAR_USEITEM_MINIMON_INFO_CLIENT_CMD
		// return ncStructData(&nc, data)
		break
	case 9231:
		// NC_BAT_SPCHANGE_CMD
		// return ncStructData(&nc, data)
		break
	case 20489:
		// NC_SOULSTONE_SP_USE_REQ
		// return ncStructData(&nc, data)
		break
	case 16421:
		// NC_CHARSAVE_UI_STATE_SAVE_REQ
		// return ncStructData(&nc, data)
		break
	case 3082:
		// NC_USER_LOGIN_ACK
		// return ncStructData(&nc, data)
		break
	case 4099:
		// NC_CHAR_LOGIN_ACK
		// return ncStructData(&nc, data)
		break
	case 28724:
		// NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD
		// return ncStructData(&nc, data)
		break
	case 4247:
		// NC_CHAR_GUILD_ACADEMY_CMD
		// return ncStructData(&nc, data)
		break
	case 8217:
		// NC_ACT_MOVERUN_CMD
		// return ncStructData(&nc, data)
		break
	case 7177:
		// NC_BRIEFINFO_MOB_CMD
		// return ncStructData(&nc, data)
		break
	case 28722:
		// NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD
		// return ncStructData(&nc, data)
		break
	case 22586:
		// NC_KQ_TEAM_TYPE_CMD
		// return ncStructData(&nc, data)
		break
	case 4149:
		// NC_CHAR_CHANGEPARAMCHANGE_CMD
		// return ncStructData(&nc, data)
		break
	case 8218:
		// NC_ACT_SOMEONEMOVERUN_CMD
		// return ncStructData(&nc, data)
		break
	case 12289:
		// NC_ITEM_CELLCHANGE_CMD
		// return ncStructData(&nc, data)
		break
	case 28684:
		// NC_CHAR_OPTION_GET_WINDOWPOS_REQ
		// return ncStructData(&nc, data)
		break
	case 28723:
		// NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD
		// return ncStructData(&nc, data)
		break
	case 8210:
		// NC_ACT_STOP_REQ
		// return ncStructData(&nc, data)
		break
	case 8228:
		// NC_ACT_JUMP_CMD
		// return ncStructData(&nc, data)
		break
	case 4187:
		// NC_CHAR_STAT_REMAINPOINT_CMD
		// return ncStructData(&nc, data)
		break
	case 6183:
		// NC_MAP_FIELD_ATTRIBUTE_CMD
		// return ncStructData(&nc, data)
		break
	case 9311:
		// NC_BAT_LPCHANGE_CMD
		// return ncStructData(&nc, data)
		break
	case 7172:
		// NC_BRIEFINFO_UNEQUIP_CMD
		// return ncStructData(&nc, data)
		break
	case 9258:
		// NC_BAT_ABSTATEINFORM_NOEFFECT_CMD
		// return ncStructData(&nc, data)
		break
	case 8254:
		// NC_ACT_MOVESPEED_CMD
		// return ncStructData(&nc, data)
		break
	case 12333:
		// NC_ITEM_REWARDINVENOPEN_ACK
		// return ncStructData(&nc, data)
		break
	case 3087:
		// NC_USER_LOGINWORLD_REQ
		// return ncStructData(&nc, data)
		break
	case 12320:
		// NC_ITEM_CHARGEDINVENOPEN_REQ
		// return ncStructData(&nc, data)
		break
	case 4327:
		// NC_CHAR_CLIENT_CARDCOLLECT_BOOKMARK_CMD
		// return ncStructData(&nc, data)
		break
	case 8216:
		// NC_ACT_SOMEONEMOVEWALK_CMD
		// return ncStructData(&nc, data)
		break
	case 2062:
		// NC_MISC_GAMETIME_ACK
		// return ncStructData(&nc, data)
		break
	case 2061:
		// NC_MISC_GAMETIME_REQ
		// return ncStructData(&nc, data)
		break
	case 22556:
		// NC_KQ_LIST_TIME_ACK
		// return ncStructData(&nc, data)
		break
	case 28685:
		// NC_CHAR_OPTION_GET_WINDOWPOS_ACK
		// return ncStructData(&nc, data)
		break
	case 3076:
		// NC_USER_XTRAP_REQ
		// return ncStructData(&nc, data)
		break
	case 22557:
		// NC_KQ_LIST_ADD_ACK
		// return ncStructData(&nc, data)
		break
	case 4330:
		// NC_CHAR_CLIENT_CARDCOLLECT_REWARD_CMD
		// return ncStructData(&nc, data)
		break
	case 12309:
		// NC_ITEM_USE_REQ
		// return ncStructData(&nc, data)
		break
	case 12298:
		// NC_ITEM_PICK_ACK
		// return ncStructData(&nc, data)
		break
	case 31750:
		// NC_PRISON_GET_REQ
		// return ncStructData(&nc, data)
		break
	case 15362:
		// NC_MENU_SERVERMENU_ACK
		// return ncStructData(&nc, data)
		break
	case 22555:
		// NC_KQ_LIST_REFRESH_REQ
		// return ncStructData(&nc, data)
		break
	case 6187:
		// NC_MAP_CAN_USE_REVIVEITEM_CMD
		// return ncStructData(&nc, data)
		break
	case 8209:
		// NC_ACT_NOTICE_CMD
		// return ncStructData(&nc, data)
		break
	case 9230:
		// NC_BAT_HPCHANGE_CMD
		// return ncStructData(&nc, data)
		break
	case 7176:
		// NC_BRIEFINFO_REGENMOB_CMD
		// return ncStructData(&nc, data)
		break
	case 3077:
		// NC_USER_XTRAP_ACK
		// return ncStructData(&nc, data)
		break
	case 4314:
		// NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD
		// return ncStructData(&nc, data)
		break
	case 36880:
		// NC_CHARGED_BOOTHSLOTSIZE_CMD
		// return ncStructData(&nc, data)
		break
	default:
		return NcRepresentation{}, errors.New(fmt.Sprintf("no struct assigned to this operation code %v", opCode))
	}
	return NcRepresentation{}, errors.New(fmt.Sprintf("no struct assigned to this operation code %v", opCode))
}

func ncStructData(nc structs.NC, data []byte) (NcRepresentation, error) {
	err := nc.Unpack(data)
	if err != nil {
		n, err := restruct.SizeOf(nc)
		if err != nil {
			log.Error(err)
		}
		log.Errorf("struct: %v, size: %v", reflect.TypeOf(nc).String(), n)
		return NcRepresentation{}, err
	}

	nr := NcRepresentation{
		Pdb:          nc.PdbAnalog(),
		UnpackedData: nc.String(),
	}

	return nr, nil
}
