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

// NC_ACT_NPCCLICK_CMD
// outbound
// 8202
func ncActNpcClickCmd(ctx context.Context, np *networking.Parameters) {
	var (
		pcne playerClicksOnNpcEvent
		mqe  queryMapEvent
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

// NC_ACT_NPCMENUOPEN_ACK
// 8221

// NC_MAP_TOWNPORTAL_REQ
// 6170

// NC_MAP_LINKOTHER_CMD
// 6154
// send zone connection info

// 15362 NC_MENU_SERVERMENU_ACK
func ncMenuServerMenuAck(ctx context.Context, np *networking.Parameters) {
	var (
		ppre playerPromptReplyEvent
		mqe  queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	ppre = playerPromptReplyEvent{
		nc: &structs.NcServerMenuAck{},
		s:  session,
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
