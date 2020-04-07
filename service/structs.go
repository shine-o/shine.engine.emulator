package service

import (
	"encoding/hex"
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
		nc := structs.NcBriefInfoAbstateChangeCmd{}
		return ncStructData(&nc, data)
	case 12303:
		// NC_ITEM_EQUIP_REQ
		nc := structs.NcItemEquipReq{}
		return ncStructData(&nc, data)
	case 20487:
		// NC_SOULSTONE_HP_USE_REQ
		// return ncStructData(&nc, data)
		break
	case 7182:
		// NC_BRIEFINFO_BRIEFINFODELETE_CMD
		nc := structs.NcBriefInfoDeleteCmd{}
		return ncStructData(&nc, data)
	case 28676:
		// NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ
		nc := structs.NcCharOptionGetShortcutSizeReq{}
		return ncStructData(&nc, data)
	case 2053:
		// NC_MISC_HEARTBEAT_ACK
		nc := structs.NcMiscHeartBeatAck{}
		return ncStructData(&nc, data)
	case 12296:
		// NC_ITEM_DROP_ACK
		nc := new(structs.NcItemDropAck)
		return ncStructData(nc, data)
	case 3175:
		// NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK
		// return ncStructData(&nc, data)
		break
	case 28677:
		// NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK
		nc := structs.NcCharOptionGetShortcutSizeAck{}
		return ncStructData(&nc, data)
	case 4324:
		// NC_CHAR_CLIENT_CARDCOLLECT_CMD
		// return ncStructData(&nc, data)
		break
	case 12297:
		// NC_ITEM_PICK_REQ
		nc := structs.NcItemPickReq{}
		return ncStructData(&nc, data)
	case 4294:
		// NC_CHAR_ADMIN_LEVEL_INFORM_CMD
		nc := structs.NcCharAdminLevelInformCmd{}
		return ncStructData(&nc, data)
	case 4114:
		// NC_CHAR_GUILD_CMD
		nc := structs.NcCharGuildCmd{}
		return ncStructData(&nc, data)
	case 37908:
		// NC_HOLY_PROMISE_LIST_CMD
		nc := structs.NcHolyPromiseListCmd{}
		return ncStructData(&nc, data)
	case 4097:
		// NC_CHAR_LOGIN_REQ
		nc := structs.NcCharLoginReq{}
		return ncStructData(&nc, data)
	case 12295:
		// NC_ITEM_DROP_REQ
		nc := structs.NcItemDropReq{}
		return ncStructData(&nc, data)
	case 7178:
		// NC_BRIEFINFO_DROPEDITEM_CMD
		nc := structs.NcBriefInfoDroppedItemCmd{}
		return ncStructData(&nc, data)
	case 7170:
		// NC_BRIEFINFO_CHANGEDECORATE_CMD
		nc := structs.NcBriefInfoChangeDecorateCmd{}
		return ncStructData(&nc, data)
	case 12321:
		// NC_ITEM_CHARGEDINVENOPEN_ACK
		nc := structs.NcItemChangedInventoryOpenAck{}
		return ncStructData(&nc, data)
	case 9256:
		// NC_BAT_ABSTATERESET_CMD
		nc := structs.NcBatAbstateResetCmd{}
		return ncStructData(&nc, data)
	case 12332:
		// NC_ITEM_REWARDINVENOPEN_REQ
		nc := structs.NcItemRewardInventoryOpenReq{}
		return ncStructData(&nc, data)
	case 31751:
		// NC_PRISON_GET_ACK
		nc := structs.NcPrisonGetAck{}
		return ncStructData(&nc, data)
	case 18476:
		// NC_SKILL_ITEMACTIONCOOLTIME_CMD
		nc := structs.SkillItemActionCoolTimeCmd{}
		return ncStructData(&nc, data)
	case 4308:
		// NC_CHAR_MYSTERYVAULT_UI_STATE_CMD
		nc := structs.CharMysteryVaultUiStateCmd{}
		return ncStructData(&nc, data)
	case 15361:
		// NC_MENU_SERVERMENU_REQ
		nc := structs.NcServerMenuReq{}
		return ncStructData(&nc, data)
	case 6149:
		// NC_MAP_LOGOUT_CMD
		nc := structs.MapLogoutCmd{}
		return ncStructData(&nc, data)
	case 3092:
		// NC_USER_LOGINWORLD_ACK
		nc := structs.NcUserLoginWorldAck{}
		return ncStructData(&nc, data)
	case 4387:
		// NC_CHAR_USEITEM_MINIMON_INFO_CLIENT_CMD
		nc := structs.CharUseItemMiniMonsterInfoClientCmd{}
		return ncStructData(&nc, data)
	case 9231:
		// NC_BAT_SPCHANGE_CMD
		nc := structs.NcBatSpChangeCmd{}
		return ncStructData(&nc, data)
	case 20489:
		// NC_SOULSTONE_SP_USE_REQ
		// return ncStructData(&nc, data)
		break
	case 16421:
		// NC_CHARSAVE_UI_STATE_SAVE_REQ
		nc := structs.NcCharUiStateSaveReq{}
		return ncStructData(&nc, data)
	case 3082:
		// NC_USER_LOGIN_ACK
		nc := structs.NcUserLoginAck{}
		return ncStructData(&nc, data)
	case 4099:
		// NC_CHAR_LOGIN_ACK
		nc := structs.NcCharLoginAck{}
		return ncStructData(&nc, data)
	case 28724:
		// NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD
		nc := structs.NcCharOptionImproveGetGameOptionCmd{}
		return ncStructData(&nc, data)
	case 4247:
		// NC_CHAR_GUILD_ACADEMY_CMD
		nc := structs.NcCharGuildAcademyCmd{}
		return ncStructData(&nc, data)
	case 8217:
		// NC_ACT_MOVERUN_CMD
		// return ncStructData(&nc, data)
		break
	case 7177:
		// NC_BRIEFINFO_MOB_CMD
		// struct values vary between clients
		nc := structs.NcBriefInfoMobCmd{}
		return ncStructData(&nc, data)
	case 28722:
		// NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD
		nc := structs.NcCharGetShortcutDataCmd{}
		return ncStructData(&nc, data)
	case 22586:
		// NC_KQ_TEAM_TYPE_CMD
		nc := structs.NcKqTeamTypeCmd{}
		return ncStructData(&nc, data)
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
		nc := structs.NcItemCellChangeCmd{}
		return ncStructData(&nc, data)
	case 28684:
		// NC_CHAR_OPTION_GET_WINDOWPOS_REQ
		// return ncStructData(&nc, data)
		break
	case 28723:
		// NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD
		nc := structs.NcCharGetKeyMapCmd{}
		return ncStructData(&nc, data)
	case 8210:
		// NC_ACT_STOP_REQ
		nc := structs.NcActStopReq{}
		return ncStructData(&nc, data)
	case 8228:
		// NC_ACT_JUMP_CMD
		// return ncStructData(&nc, data)
		break
	case 4187:
		// NC_CHAR_STAT_REMAINPOINT_CMD
		nc := structs.NcCharStatRemainPointCmd{}
		return ncStructData(&nc, data)
	case 6183:
		// NC_MAP_FIELD_ATTRIBUTE_CMD
		nc := structs.NcMapFieldAttributeCmd{}
		return ncStructData(&nc, data)
	case 9311:
		// NC_BAT_LPCHANGE_CMD
		nc := structs.NcBatLpChangeCmd{}
		return ncStructData(&nc, data)
	case 7172:
		// NC_BRIEFINFO_UNEQUIP_CMD
		nc := structs.NcBriefInfoUnequipCmd{}
		return ncStructData(&nc, data)
	case 9258:
		// NC_BAT_ABSTATEINFORM_NOEFFECT_CMD
		nc := structs.NcBatAbstateInformNoEffectCmd{}
		return ncStructData(&nc, data)
	case 8254:
		// NC_ACT_MOVESPEED_CMD
		nc := structs.NcActMoveSpeedCmd{}
		return ncStructData(&nc, data)
	case 12333:
		// NC_ITEM_REWARDINVENOPEN_ACK
		nc := structs.NcItemRewardInventoryOpenAck{}
		return ncStructData(&nc, data)
	case 3087:
		// NC_USER_LOGINWORLD_REQ
		nc := structs.NcUserLoginWorldReq{}
		return ncStructData(&nc, data)
	case 12320:
		// NC_ITEM_CHARGEDINVENOPEN_REQ
		nc := structs.NcITemChargedInventoryOpenReq{}
		return ncStructData(&nc, data)
	case 4327:
		// NC_CHAR_CLIENT_CARDCOLLECT_BOOKMARK_CMD
		// return ncStructData(&nc, data)
		break
	case 8216:
		// NC_ACT_SOMEONEMOVEWALK_CMD
		nc := structs.NcActSomeoneMoveWalkCmd{}
		return ncStructData(&nc, data)
	case 2062:
		// NC_MISC_GAMETIME_ACK
		nc := structs.NcMiscGameTimeAck{}
		return ncStructData(&nc, data)
	case 2061:
		// NC_MISC_GAMETIME_REQ
		// return ncStructData(&nc, data)
		break
	case 22556:
		// NC_KQ_LIST_TIME_ACK
		nc := structs.NcKqListTimeAck{}
		return ncStructData(&nc, data)
	case 28685:
		// NC_CHAR_OPTION_GET_WINDOWPOS_ACK
		nc := structs.NcCharOptionGetWindowPosAck{}
		return ncStructData(&nc, data)
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
		nc := structs.NcItemUseReq{}
		return ncStructData(&nc, data)
	case 12298:
		// NC_ITEM_PICK_ACK
		nc := structs.NcItemPickAck{}
		return ncStructData(&nc, data)
	case 31750:
		// NC_PRISON_GET_REQ
		// return ncStructData(&nc, data)
		break
	case 15362:
		// NC_MENU_SERVERMENU_ACK
		nc := structs.NcServerMenuAck{}
		return ncStructData(&nc, data)
	case 22555:
		// NC_KQ_LIST_REFRESH_REQ
		// return ncStructData(&nc, data)
		break
	case 6187:
		// NC_MAP_CAN_USE_REVIVEITEM_CMD
		nc := structs.NcMapCanUseReviveItemCmd{}
		return ncStructData(&nc, data)
	case 8209:
		// NC_ACT_NOTICE_CMD
		// return ncStructData(&nc, data)
		break
	case 9230:
		// NC_BAT_HPCHANGE_CMD
		nc := structs.NcBatHpChangeCmd{}
		return ncStructData(&nc, data)
	case 7176:
		// NC_BRIEFINFO_REGENMOB_CMD
		nc := structs.NcBriefInfoRegenMobCmd{}
		return ncStructData(&nc, data)
	case 3077:
		// NC_USER_XTRAP_ACK
		// return ncStructData(&nc, data)
		break
	case 4314:
		// NC_CHAR_NEWBIE_GUIDE_VIEW_SET_CMD
		nc := structs.NcCharNewbieGuideViewSetCmd{}
		return ncStructData(&nc, data)
	case 36880:
		// NC_CHARGED_BOOTHSLOTSIZE_CMD
		nc := structs.NcChargedBoothSlotSizeCmd{}
		return ncStructData(&nc, data)
	case 9277:
		// NC_BAT_CEASE_FIRE_CMD
		nc := structs.NcBatCeaseFireCmd{}
		return ncStructData(&nc, data)
	case 9259:
		// NC_BAT_BASHSTART_CMD
		// return ncStructData(&nc, data)
		break
	case 7193:
		// NC_BRIEFINFO_ABSTATE_CHANGE_LIST_CMD
		nc := structs.NcBriefInfoAbstateChangeListCmd{}
		return ncStructData(&nc, data)
	case 9257:
		// NC_BAT_ABSTATEINFORM_CMD
		nc := structs.NcBatAbstateInformCmd{}
		return ncStructData(&nc, data)
	case 49169:
		// oh my
		// return ncStructData(&nc, data)
		break
	case 26631:
		// NC_BOOTH_ENTRY_REQ
		nc := structs.NcBoothEntryReq{}
		return ncStructData(&nc, data)
	case 12306:
		// NC_ITEM_UNEQUIP_REQ
		nc := structs.NcItemUnequipReq{}
		return ncStructData(&nc, data)
	case 9224:
		// NC_BAT_UNTARGET_REQ
		// return ncStructData(&nc, data)
		break
	case 52237:
		// NC_MOVER_MOVESPEED_CMD
		nc := structs.NcMoverMoveSpeedCmd{}
		return ncStructData(&nc, data)
	case 12365:
		// NC_ITEM_ACCOUNT_STORAGE_CLOSE_CMD

		// return ncStructData(&nc, data)
		break
	case 9280:
		// NC_BAT_SKILLBASH_OBJ_CAST_REQ
		nc := structs.NcBatSkillBashObjCastReq{}
		return ncStructData(&nc, data)
	case 26627:
		// NC_BOOTH_SOMEONEOPEN_CMD
		nc := structs.NcBoothSomeoneOpenCmd{}
		return ncStructData(&nc, data)
	case 49168:
		//
		// return ncStructData(&nc, data)
		break
	case 50184:
		// NC_COLLECT_CARDREGIST_REQ
		nc := structs.NcCollectCardRegisterReq{}
		return ncStructData(&nc, data)
	case 17428:
		// NC_QUEST_START_REQ
		nc := structs.NcQuestStartReq{}
		return ncStructData(&nc, data)
	case 7175:
		// NC_BRIEFINFO_CHARACTER_CMD
		nc := structs.NcBriefInfoCharacterCmd{}
		return ncStructData(&nc, data)
	case 7194:
		// NC_BRIEFINFO_REGENMOVER_CMD
		nc := structs.NcBriefInfoRegenMoverCmd{}
		return ncStructData(&nc, data)
	case 8237:
		// NC_ACT_GATHERSTART_REQ
		nc := structs.NcActGatherStartReq{}
		return ncStructData(&nc, data)
	case 12299:
		// NC_ITEM_RELOC_REQ
		nc := structs.NcitemRelocateReq{}
		return ncStructData(&nc, data)
	case 8223:
		// NC_ACT_SOMEONESHOUT_CMD
		nc := structs.NcActSomeoneShoutCmd{}
		return ncStructData(&nc, data)
	case 9217:
		// NC_BAT_TARGETTING_REQ
		// return ncStructData(&nc, data)
		break
	case 52234:
		// NC_MOVER_HUNGRY_CMD
		nc := structs.NcMoverHungryCmd{}
		return ncStructData(&nc, data)
	case 26634:
		// NC_BOOTH_REFRESH_REQ
		nc := structs.NcBoothRefreshReq{}
		return ncStructData(&nc, data)
	case 17410:
		// NC_QUEST_SCRIPT_CMD_ACK
		nc := structs.NcQuestScriptCmdAck{}
		return ncStructData(&nc, data)
	case 7171:
		// NC_BRIEFINFO_CHANGEUPGRADE_CMD
		nc := structs.NcBriefInfoChangeUpgradeCmd{}
		return ncStructData(&nc, data)
	case 7195:
		// NC_BRIEFINFO_MOVER_CMD
		nc := structs.NcBriefInfoMoverCmd{}
		return ncStructData(&nc, data)
	case 9255:
		// NC_BAT_ABSTATESET_CMD
		nc := structs.NcBatAbstateSetCmd{}
		return ncStructData(&nc, data)
	case 7174:
		// NC_BRIEFINFO_LOGINCHARACTER_CMD
		nc := structs.NcBriefInfoLoginCharacterCmd{}
		return ncStructData(&nc, data)
	case 8211:
		// NC_ACT_SOMEONESTOP_CMD
		nc := structs.NcActSomeoneStopCmd{}
		return ncStructData(&nc, data)
	case 52226:
		// NC_MOVER_RIDE_ON_CMD
		nc := structs.NcMoverRideOnCmd{}
		return ncStructData(&nc, data)
	case 8242:
		// NC_ACT_GATHERCOMPLETE_REQ
		// return ncStructData(&nc, data)
		break
	case 8202:
		// NC_ACT_NPCCLICK_CMD
		nc := structs.NcActNpcClickCmd{}
		return ncStructData(&nc, data)
	case 7173:
		// NC_BRIEFINFO_CHANGEWEAPON_CMD
		nc := structs.NcBriefInfoChangeWeaponCmd{}
		return ncStructData(&nc, data)
	case 52228:
		// NC_MOVER_SOMEONE_RIDE_ON_CMD
		nc := structs.NcMoverSomeoneRideOnCmd{}
		return ncStructData(&nc, data)
	case 9276:
		// NC_BAT_DOTDAMAGE_CMD
		nc := structs.NcBatDotDamageCmd{}
		return ncStructData(&nc, data)
	case 8200:
		// NC_ACT_CHANGEMODE_REQ
		nc := structs.NcActChangeModeReq{}
		return ncStructData(&nc, data)
	case 26648:
		//
		// return ncStructData(&nc, data)
		break
	case 8221:
		// NC_ACT_NPCMENUOPEN_ACK
		// return ncStructData(&nc, data)
		break
	case 9218:
		// NC_BAT_TARGETINFO_CMD
		nc := structs.NcBatTargetInfoCmd{}
		return ncStructData(&nc, data)
	case 9266:
		// NC_BAT_BASHSTOP_CMD
		// return ncStructData(&nc, data)
		break
	case 2052:
		// NC_MISC_HEARTBEAT_REQ
		// return ncStructData(&nc, data)
		break
	case 4286:
		// NC_CHAR_CLIENT_AUTO_PICK_CMD
		nc := structs.NcCharClientAutoPickCmd{}
		return ncStructData(&nc, data)
	case 8248:
		// NC_ACT_SOMEONEPRODUCE_CAST_CMD
		nc := structs.NcActSomeoneProduceCastCmd{}
		return ncStructData(&nc, data)
	case 20492:
		// NC_SOULSTONE_SP_SOMEONEUSE_CMD
		nc := structs.NcSoulStoneSpSomeoneUseCmd{}
		return ncStructData(&nc, data)
	case 26632:
		// NC_BOOTH_ENTRY_SELL_ACK
		nc := structs.NcBoothEntrySellAck{}
		return ncStructData(&nc, data)
	case 9298:
		// NC_BAT_SKILLBASH_HIT_DAMAGE_CMD
		nc := structs.NcBatSkillBashHitDamageCmd{}
		return ncStructData(&nc, data)
	case 8236:
		// NC_ACT_SOMEONEFOLDTENT_CMD
		nc := structs.NcActSomeoneFoldTentCmd{}
		return ncStructData(&nc, data)
	case 9295:
		// NC_BAT_SOMEONESKILLBASH_HIT_OBJ_START_CMD
		nc := structs.NcBatSomeoneSkillBashHitObjStartCmd{}
		return ncStructData(&nc, data)
	case 52232:
		// NC_MOVER_SOMEONE_RIDE_OFF_CMD
		nc := structs.NcMoverSomeoneRideOffCmd{}
		return ncStructData(&nc, data)
	case 26635:
		// NC_BOOTH_REFRESH_SELL_ACK
		// return ncStructData(&nc, data)
		break
	case 2054:
		// NC_MISC_SEED_REQ
		// return ncStructData(&nc, data)
		break
	case 26647:
		// NC_BOOTH_SEARCH_BOOTH_CLOSED_CMD
		nc := structs.NcBoothSearchBoothClosedCmd{}
		return ncStructData(&nc, data)
	case 4396:
		// NC_CHAR_USEITEM_MINIMON_USE_BROAD_CMD
		nc := structs.NcCharUseItemMinimonUseBroadCmd{}
		return ncStructData(&nc, data)
	case 9303:
		// NC_BAT_SKILLBASH_HIT_BLAST_CMD
		nc := structs.NcBatSkillBashHitBlastCmd{}
		return ncStructData(&nc, data)
	case 20491:
		// NC_SOULSTONE_HP_SOMEONEUSE_CMD
		nc := structs.NcSoulStoneHpSomeoneUseCmd{}
		return ncStructData(&nc, data)
	case 7169:
		// NC_BRIEFINFO_INFORM_CMD
		nc := structs.NcBriefInfoInformCmd{}
		return ncStructData(&nc, data)
	case 8201:
		// NC_ACT_SOMEONECHANGEMODE_CMD
		nc := structs.NcActSomeoneChangeModeCmd{}
		return ncStructData(&nc, data)
	case 8252:
		// NC_ACT_SOMEONEPRODUCE_MAKE_CMD
		nc := structs.NcActSomeoneProduceMakeCmd{}
		return ncStructData(&nc, data)
	default:
		return NcRepresentation{}, errors.New(fmt.Sprintf("no struct assigned to this operation code %v", opCode))
	}
	return NcRepresentation{}, errors.New(fmt.Sprintf("no struct assigned to this operation code %v", opCode))
}

func ncStructData(nc structs.NC, data []byte) (NcRepresentation, error) {
	err := structs.Unpack(data, nc)
	if err != nil {
		n, err := restruct.SizeOf(nc)
		if err != nil {
			log.Error(err)
		}
		hexString := hex.EncodeToString(data)
		log.Error(hexString)
		log.Errorf("struct: %v, size: %v", reflect.TypeOf(nc).String(), n)
		return NcRepresentation{}, err
	}
	//n, err := restruct.SizeOf(nc)
	//if err != nil {
	//	log.Error(err)
	//}
	//log.Warningf("struct: %v, size: %v", reflect.TypeOf(nc).String(), n)
	nr := NcRepresentation{
		Pdb:          nc.PdbType(),
		UnpackedData: nc.String(),
	}

	return nr, nil
}
