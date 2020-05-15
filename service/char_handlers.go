package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NC_MAP_LOGIN_REQ
func ncMapLoginReq(ctx context.Context, np *networking.Parameters) {
	var (
		nc   structs.NcMapLoginReq
		cse  clientSHNEvent
		cde  playerDataEvent
		mqe  queryMapEvent
		rphe registerPlayerHandleEvent
	)

	err := structs.Unpack(np.Command.Base.Data, &nc)
	if err != nil {
		log.Error(err)
		return
	}

	cse = clientSHNEvent{
		inboundNC: nc,
		ok:        make(chan bool),
		err:       make(chan error),
	}

	zoneEvents[clientSHN] <- &cse

	cde = playerDataEvent{
		player:     make(chan *player),
		net:        np,
		playerName: nc.CharData.CharID.Name,
		err:        make(chan error),
	}

	zoneEvents[loadPlayerData] <- &cde

	select {
	case <-cse.ok:
		break
	case err := <-cse.err:
		log.Error(err)
		// fail ack with failure code
		// drop connection
		return
	}

	var p *player
	select {
	case p = <-cde.player:
		break
	case err := <-cde.err:
		log.Error(err)
		// fail ack with failure code
		// drop connection
		return
	}

	mqe = queryMapEvent{
		id:  p.location.mapID,
		zm:  make(chan *zoneMap),
		err: make(chan error),
	}

	zoneEvents[queryMap] <- &mqe

	var zm *zoneMap
	select {
	case zm = <-mqe.zm:
		break
	case err := <-mqe.err:
		log.Error(err)
		return
	}

	sv := ctx.Value(networking.ShineSession)

	session, ok := sv.(*session)

	if !ok {
		log.Errorf("no session available for player %v", p.view.name)
		return
	}

	rphe = registerPlayerHandleEvent{
		player:  p,
		session: session,
		done:    make(chan bool),
		err:     make(chan error),
	}

	zm.send[registerPlayerHandle] <- &rphe

	select {
	case <-rphe.done:
		ncCharClientBaseCmd(p)
		ncCharClientShapeCmd(p)
		// weird bug sometimes the client stucks in character select
		ncMapLoginAck(p)
	case err := <-rphe.err:
		log.Error(err)
	}
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
					Name: p.location.mapName,
				},
				CurrentCoord: structs.ShineCoordType{
					XY: structs.ShineXYType{
						X: p.location.x,
						Y: p.location.y,
					},
					Direction: p.location.d,
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
	pc.SendDirectly(p.conn.outboundData)
}

//NC_CHAR_CLIENT_SHAPE_CMD
func ncCharClientShapeCmd(p *player) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4153,
		},
		NcStruct: p.view.protoAvatarShapeInfo(),
	}
	pc.SendDirectly(p.conn.outboundData)
}

//NC_CHAR_CLIENT_ITEM_CMD
func ncCharClientItemCmd(p *player, nc structs.NcCharClientItemCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4167,
		},
		NcStruct: &nc,
	}
	pc.SendDirectly(p.conn.outboundData)
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
				X: p.location.x,
				Y: p.location.y,
			},
		},
	}
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
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
	pc.SendDirectly(p.conn.outboundData)
}

// NC_MAP_LOGINCOMPLETE_CMD
// 6147
func ncMapLoginCompleteCmd(ctx context.Context, np *networking.Parameters) {
	// fetch user session
	var (
		mqe queryMapEvent
		pae playerAppearedEvent
	)

	sv := ctx.Value(networking.ShineSession)

	session, ok := sv.(*session)

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
		playerHandle: session.handle,
		mapID:        session.mapID,
		err:          make(chan error),
	}

	zm.send[playerAppeared] <- &pae

	// there's the mapID, handle ID

	// start the heartbeat for this player

	// send info about surrounding mobs, players to this player

	// to all surrounding players, send info about this player
}

//NC_CHAR_LOGOUTREADY_CMD
func ncCharLogoutReadyCmd(ctx context.Context, np *networking.Parameters) {
	// start a ticker that in 10 seconds will close the connection
	// another packet can be received which will cancel that ticker
	np.NetVars.CloseConnection <- true
}
