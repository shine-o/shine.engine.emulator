package zone

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// NC_MAP_LOGIN_REQ
func ncMapLoginReq(ctx context.Context, np *networking.Parameters) {
	var (
		nc   structs.NcMapLoginReq
		pmle playerMapLoginEvent
	)

	err := structs.Unpack(np.Command.Base.Data, &nc)
	if err != nil {
		log.Error(err)
		return
	}

	pmle = playerMapLoginEvent{
		nc: nc,
		np: np,
	}

	zoneEvents[playerMapLogin] <- &pmle
}

//NC_CHAR_CLIENT_BASE_CMD
func ncCharClientBaseCmd(p *player) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4152,
		},
		NcStruct: &structs.NcCharClientBaseCmd{
			ChrRegNum: uint32(p.char.ID),
			CharName: structs.Name5{
				Name: p.view.name,
			},
			Slot:       p.char.Slot,
			Level:      p.state.level,
			Experience: p.state.exp,
			PwrStone:   0,
			GrdStone:   0,
			HPStone:    p.stats.hpStones,
			SPStone:    p.stats.spStones,
			CurHP:      p.stats.hp,
			CurSP:      p.stats.sp,
			CurLP:      p.stats.lp,
			Unk:        1,
			Fame:       p.money.fame,
			Cen:        54983635, // ¿?
			LoginInfo: structs.NcCharBaseCmdLoginLocation{
				CurrentMap: structs.Name3{
					Name: p.current.mapName,
				},
				CurrentCoord: structs.ShineCoordType{
					XY: structs.ShineXYType{
						X: uint32(p.current.x),
						Y: uint32(p.current.y),
					},
					Direction: uint8(p.current.d),
				},
			},
			Stats: structs.CharStats{
				Strength:          p.stats.points.str,
				Constitute:        p.stats.points.end,
				Dexterity:         p.stats.points.dex,
				Intelligence:      p.stats.points.int,
				MentalPower:       p.stats.points.spr,
				RedistributePoint: p.stats.points.redistributionPoints,
			},
			IdleTime:   0,
			PkCount:    p.char.Attributes.KillPoints,
			PrisonMin:  0,
			AdminLevel: p.char.AdminLevel,
			Flag: structs.NcCharBaseCmdFlag{
				Val: 0,
			},
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_SHAPE_CMD
func ncCharClientShapeCmd(p *player) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4153,
		},
		NcStruct: p.view.protoAvatarShapeInfo(),
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_ITEM_CMD
func ncCharClientItemCmd(p *player, nc structs.NcCharClientItemCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4167,
		},
		NcStruct: &nc,
	}
	pc.Send(p.conn.outboundData)
}

//NC_MAP_LOGIN_ACK
func ncMapLoginAck(p *player) {
	// handle ID
	// character complete Parameters (resultant stats from base + items )
	// todo: character stat calculation given the assigned stats, title, equipped items, abstate
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 6146,
		},
		NcStruct: &structs.NcMapLoginAck{
			Handle: p.handle, // id of the entity inside this map
			//Handle: 2222, // id of the entity inside this map
			Params: p.charParameterData(),
			LoginCoord: structs.ShineXYType{
				X: uint32(p.current.x),
				Y: uint32(p.current.y),
			},
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_QUEST_READ_CMD
//4302
func ncCharClientQuestReadCmd(p *player) {
	// todo: quest logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4302,
		},
		NcStruct: &structs.NcCharClientQuestReadCmd{
			CharID:          uint32(p.char.ID),
			NumOfReadQuests: 0,
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_QUEST_DOING_CMD
//4154
func ncCharClientQuestDoingCmd(p *player) {
	// todo: quest logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4154,
		},
		NcStruct: &structs.NcCharClientQuestDoingCmd{
			CharID:          uint32(p.char.ID),
			NeedClear:       0,
			NumOfDoingQuest: 0,
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_QUEST_DONE_CMD
//4155
func ncCharClientQuestDoneCmd(p *player) {
	// todo: quest logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4155,
		},
		NcStruct: &structs.NcCharClientQuestDoneCmd{
			CharID:             uint32(p.char.ID),
			TotalDoneQuest:     0,
			TotalDoneQuestSize: 0,
			Count:              0,
			Index:              0,
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_QUEST_REPEAT_CMD
//4311
func ncCharClientQuestRepeatCmd(p *player) {
	// todo: quest logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4311,
		},
		NcStruct: &structs.NcCharClientQuestRepeatCmd{
			CharID: uint32(p.char.ID),
			Count:  0,
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_PASSIVE_CMD
//4158
func ncCharClientPassiveCmd(p *player) {
	// todo: skill logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4158,
		},
		NcStruct: &structs.NcCharClientPassiveCmd{
			Number: 0,
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_SKILL_CMD
//4157
func ncCharClientSkillCmd(p *player) {
	// todo: skill logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4157,
		},
		NcStruct: &structs.NcCharClientSkillCmd{
			RestEmpower: 0,
			PartMark: structs.PartMark{
				BF0: 0,
			},
			MaxNum: 0,
			Skills: structs.NcCharSkillClientCmd{
				ChrRegNum: uint32(p.char.ID),
				Number:    0,
			},
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_CHARTITLE_CMD
//4169
func ncCharClientCharTitleCmd(p *player) {
	// todo: character title logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4169,
		},
		NcStruct: &structs.NcClientCharTitleCmd{
			CurrentTitle:        0,
			CurrentTitleElement: 0,
			CurrentTitleMobID:   0,
			NumOfTitle:          0,
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_GAME_CMD
//4168
func ncCharClientGameCmd(p *player) {
	// no idea what this is f	or
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4168,
		},
		NcStruct: &structs.NcCharClientGameCmd{
			Filler0: 65535,
			Filler1: 65535,
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_CHARGEDBUFF_CMD
//4170
func ncCharClientChargedBuffCmd(p *player) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4170,
		},
		NcStruct: &structs.NcCharClientChargedBuffCmd{
			Count: 0,
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_CHAR_CLIENT_COININFO_CMD
//4318
func ncCharClientCoinInfoCmd(p *player) {
	// todo: money && fame logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4318,
		},
		NcStruct: &structs.NcCharClientCoinInfoCmd{
			Coin:          73560,
			ExchangedCoin: 100000,
		},
	}
	pc.Send(p.conn.outboundData)
}

//NC_QUEST_RESET_TIME_CLIENT_CMD
//17438
func ncQuestResetTimeClientCmd(p *player) {
	// todo: quest logic
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 17438,
		},
		NcStruct: &structs.NcQuestResetTimeClientCmd{
			ResetYear:  1577862000,
			ResetMonth: 1585724400,
			ResetWeek:  1586761200,
			ResetDay:   1587279600,
		},
	}
	pc.Send(p.conn.outboundData)
}

// NC_MAP_LOGINCOMPLETE_CMD
// 6147
func ncMapLoginCompleteCmd(ctx context.Context, np *networking.Parameters) {
	// fetch user session
	var (
		mqe queryMapEvent
		pae playerAppearedEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	mqe = queryMapEvent{
		id:  session.mapID,
		zm:  make(chan *zoneMap),
		err: make(chan error),
	}

	zoneEvents[queryMap] <- &mqe

	var zm *zoneMap
	select {
	case zm = <-mqe.zm:
		break
	case e := <-mqe.err:
		log.Error(e)
		return
	}

	pae = playerAppearedEvent{
		handle: session.handle,
	}

	zm.send[playerAppeared] <- &pae

}

//4210
func ncCharLogoutCancelCmd(ctx context.Context, np *networking.Parameters) {
	var (
		plce playerLogoutCancelEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	plce = playerLogoutCancelEvent{
		sessionID: session.id,
	}

	zoneEvents[playerLogoutCancel] <- &plce
}

//NC_CHAR_LOGOUTREADY_CMD
func ncCharLogoutReadyCmd(ctx context.Context, np *networking.Parameters) {
	var (
		plse playerLogoutStartEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	plse = playerLogoutStartEvent{
		sessionID: session.id,
		mapID:     session.mapID,
		handle:    session.handle,
	}

	zoneEvents[playerLogoutStart] <- &plse
}
