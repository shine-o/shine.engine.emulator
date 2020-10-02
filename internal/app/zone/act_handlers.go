package zone

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// NC_ACT_MOVEWALK_CMD
// 8215
// walk
func ncActMoveWalkCmd(ctx context.Context, np *networking.Parameters) {
	var (
		pwe playerWalksEvent
		mqe queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	pwe = playerWalksEvent{
		nc:     &structs.NcActMoveRunCmd{},
		handle: session.handle,
	}

	err := structs.Unpack(np.Command.Base.Data, pwe.nc)
	if err != nil {
		log.Error(err)
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

	zm.send[playerWalks] <- &pwe
}

// NC_ACT_MOVERUN_CMD
// 8217
// run
func ncActMoveRunCmd(ctx context.Context, np *networking.Parameters) {
	var (
		pre playerRunsEvent
		mqe queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	pre = playerRunsEvent{
		nc:     &structs.NcActMoveRunCmd{},
		handle: session.handle,
	}

	err := structs.Unpack(np.Command.Base.Data, pre.nc)
	if err != nil {
		log.Error(err)
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

	zm.send[playerRuns] <- &pre
}

// NC_ACT_JUMP_CMD
// 8228
// jump
func ncActJumpCmd(ctx context.Context, np *networking.Parameters) {
	var (
		pje playerJumpedEvent
		mqe queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	pje = playerJumpedEvent{
		handle: session.handle,
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

	zm.send[playerJumped] <- &pje
}

// NC_ACT_STOP_REQ
// 8210
// stop walk/run, a.k.a last known position of entity
func ncActStopReq(ctx context.Context, np *networking.Parameters) {
	var (
		pse playerStoppedEvent
		mqe queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	pse = playerStoppedEvent{
		nc:     &structs.NcActStopReq{},
		handle: session.handle,
	}

	err := structs.Unpack(np.Command.Base.Data, pse.nc)
	if err != nil {
		log.Error(err)
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

	zm.send[playerStopped] <- &pse
}

// NC_ACT_SOMEONEMOVEWALK_CMD
// 8216
// someone walked
func ncActSomeoneMoveWalkCmd(p *player, nc *structs.NcActSomeoneMoveWalkCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 8216,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}

// NC_ACT_SOMEONEMOVERUN_CMD
// 8218
// someone has run
func ncActSomeoneMoveRunCmd(p *player, nc *structs.NcActSomeoneMoveRunCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 8218,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}

// NC_ACT_SOMEONEJUMP_CMD
// 8229
func ncActSomeoneJumpCmd(p *player, nc *structs.NcActSomeoneJumpCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 8229,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}

// NC_ACT_SOMEONESTOP_CMD
// 8211
// someone stopped
func ncActSomeoneStopCmd(p *player, nc *structs.NcActSomeoneStopCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 8211,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}

// NC_ACT_NPCCLICK_CMD
// outbound
// 8202
func ncActNpcClickCmd(ctx context.Context, np *networking.Parameters)  {
	var (
		pcne playerClicksOnNpcEvent
		mqe queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	pcne = playerClicksOnNpcEvent{
		nc:     &structs.NcActNpcClickCmd{},
		handle: session.handle,
	}

	err := structs.Unpack(np.Command.Base.Data, pcne.nc)

	if err != nil {
		log.Error(err)
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

	zm.events.send[playerClicksOnNpc] <- &pcne

}

// NC_ACT_NPCMENUOPEN_REQ
// inbound
// 8220
func ncActNpcMenuOpenReq(p *player, nc *structs.NcActNpcMenuOpenReq) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 8220,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}

// NC_ACT_NPCMENUOPEN_ACK
// 8221

// NC_MAP_TOWNPORTAL_REQ
// 6170

// NC_MAP_LINKSAME_CMD
// 6153
// move player to another map
func ncMapLinkSameCmd(p *player, nc *structs.NcMapLinkSameCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 6153,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}


// NC_MAP_LINKOTHER_CMD
// 6154
// send zone connection info

// INFO : 2020/10/01 19:40:11.294396 handlers.go:267: 2020-10-01 19:40:11.29071 +0200 CEST 9224->9312 inbound NC_BAT_TARGETINFO_CMD {"packetType":"small","length":32,"department":9,"command":"2","opCode":9218,"data":"00c642000000000000000064000000640000000000000000000000001c00","rawData":"20022400c642000000000000000064000000640000000000000000000000001c00","friendlyName":"NC_BAT_TARGETINFO_CMD"}
// INFO : 2020/10/01 19:40:11.295396 handlers.go:267: 2020-10-01 19:40:11.29071 +0200 CEST 9224->9312 inbound NC_MENU_SERVERMENU_REQ {"packetType":"small","length":210,"department":15,"command":"1","opCode":15361,"data":"446f20796f752077616e7420746f206d6f766520746f20526f756d656e206669656c643f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c64261440000213600002c0102005965730000000000000000000000000000000000000000000000000000000000014e6f000000000000000000000000000000000000000000000000000000000000","rawData":"d2013c446f20796f752077616e7420746f206d6f766520746f20526f756d656e206669656c643f000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c64261440000213600002c0102005965730000000000000000000000000000000000000000000000000000000000014e6f000000000000000000000000000000000000000000000000000000000000","friendlyName":"NC_MENU_SERVERMENU_REQ"}

// NC_MENU_SERVERMENU_REQ
// 15361
func ncMenuServerMenuReq(p *player, nc *structs.NcServerMenuReq) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 15361,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}

// 15362 NC_MENU_SERVERMENU_ACK
func ncMenuServerMenuAck(ctx context.Context, np *networking.Parameters)  {
	var (
		ppre playerPromptReplyEvent
		mqe queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	ppre = playerPromptReplyEvent{
		nc:     &structs.NcServerMenuAck{},
		s: session,
	}

	err := structs.Unpack(np.Command.Base.Data, ppre.nc)

	if err != nil {
		log.Error(err)
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

	zm.events.send[playerPromptReply] <- &ppre

}
