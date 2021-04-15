package zone

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NC_ACT_MOVEWALK_CMD
// 8215
// walk
func ncActMoveWalkCmd(ctx context.Context, np *networking.Parameters) {
	var (
		e playerWalksEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = playerWalksEvent{
		nc:     &structs.NcActMoveRunCmd{},
		handle: session.handle,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)
	if err != nil {
		log.Error(err)
		return
	}

	zm, ok := maps.list[session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code:    errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}
	zm.send[playerWalks] <- &e
}

// NC_ACT_MOVERUN_CMD
// 8217
func ncActMoveRunCmd(ctx context.Context, np *networking.Parameters) {
	var (
		e playerRunsEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = playerRunsEvent{
		nc:     &structs.NcActMoveRunCmd{},
		handle: session.handle,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)
	if err != nil {
		log.Error(err)
		return
	}

	zm, ok := maps.list[session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code:    errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}

	zm.send[playerRuns] <- &e
}

// NC_ACT_JUMP_CMD
// 8228
// jump
func ncActJumpCmd(ctx context.Context, np *networking.Parameters) {
	var (
		e playerJumpedEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = playerJumpedEvent{
		handle: session.handle,
	}

	zm, ok := maps.list[session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code:    errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}

	zm.send[playerJumped] <- &e
}

// NC_ACT_STOP_REQ
// 8210
// stop walk/run, a.k.a last known position of entity
func ncActStopReq(ctx context.Context, np *networking.Parameters) {
	var (
		e playerStoppedEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = playerStoppedEvent{
		nc:     &structs.NcActStopReq{},
		handle: session.handle,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)
	if err != nil {
		log.Error(err)
		return
	}

	zm, ok := maps.list[session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code:    errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}

	zm.send[playerStopped] <- &e
}

// NC_ACT_NPCCLICK_CMD
// outbound
// 8202
func ncActNpcClickCmd(ctx context.Context, np *networking.Parameters) {
	var (
		e playerClicksOnNpcEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = playerClicksOnNpcEvent{
		nc:     &structs.NcActNpcClickCmd{},
		handle: session.handle,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)

	if err != nil {
		log.Error(err)
		return
	}

	zm, ok := maps.list[session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code:    errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}

	zm.events.send[playerClicksOnNpc] <- &e

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
		e playerPromptReplyEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = playerPromptReplyEvent{
		nc: &structs.NcServerMenuAck{},
		s:  session,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)

	if err != nil {
		log.Error(err)
		return
	}

	zm, ok := maps.list[session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code:    errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}

	zm.events.send[playerPromptReply] <- &e

}
