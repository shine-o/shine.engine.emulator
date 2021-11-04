package zone

import (
	"context"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NC_ACT_MOVEWALK_CMD
func ncActMoveWalkCmd(ctx context.Context, np *networking.Parameters) {
	var e playerWalksEvent

	session, ok := np.Session.(*session)

	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneNoSessionAvailable,
		})
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
			Code: errors.ZoneMapNotFound,
			Details: errors.Details{
				"session": session,
			},
		})
		return
	}
	zm.events.send[playerWalks] <- &e
}

// NC_ACT_MOVERUN_CMD
func ncActMoveRunCmd(ctx context.Context, np *networking.Parameters) {
	var e playerRunsEvent

	session, ok := np.Session.(*session)

	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneNoSessionAvailable,
		})
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
			Code: errors.ZoneMapNotFound,
			Details: errors.Details{
				"session": session,
			},
		})
		return
	}

	zm.events.send[playerRuns] <- &e
}

// NC_ACT_EMOTICON_CMD
func ncActEmoticonCmd(ctx context.Context, np *networking.Parameters) {
	var e playerEmotedEvent

	session, ok := np.Session.(*session)

	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneNoSessionAvailable,
		})
		return
	}

	e = playerEmotedEvent{
		handle: session.handle,
	}

	zm, ok := maps.list[session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneMapNotFound,
			Details: errors.Details{
				"session": session,
			},
		})
		return
	}

	zm.events.send[playerEmoted] <- &e
}

// NC_ACT_JUMP_CMD
func ncActJumpCmd(ctx context.Context, np *networking.Parameters) {
	var e playerJumpedEvent

	session, ok := np.Session.(*session)

	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneNoSessionAvailable,
		})
		return
	}

	e = playerJumpedEvent{
		handle: session.handle,
	}

	zm, ok := maps.list[session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneMapNotFound,
			Details: errors.Details{
				"session": session,
			},
		})
		return
	}

	zm.events.send[playerJumped] <- &e
}

// NC_ACT_STOP_REQ
func ncActStopReq(ctx context.Context, np *networking.Parameters) {
	var e playerStoppedEvent

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
			Code: errors.ZoneMapNotFound,
			Details: errors.Details{
				"session": session,
			},
		})
		return
	}

	zm.events.send[playerStopped] <- &e
}

// NC_ACT_NPCCLICK_CMD
func ncActNpcClickCmd(ctx context.Context, np *networking.Parameters) {
	var e playerClicksOnNpcEvent

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
			Code: errors.ZoneMapNotFound,
			Details: errors.Details{
				"session": session,
			},
		})
		return
	}

	zm.events.send[playerClicksOnNpc] <- &e
}

func ncActNpcMenuOpenAck(ctx context.Context, np *networking.Parameters) {
	// INFO : 2021/04/26 10:00:13.035168 handlers.go:272: 2021-04-26 10:00:13.021596 +0200 CEST 6233->9120 outbound NC_ACT_NPCMENUOPEN_ACK {"packetType":"small","length":3,"department":8,"command":"1D","opCode":8221,"data":"01","rawData":"031d2001","friendlyName":""}
}

// NC_ACT_NPCMENUOPEN_ACK
// 8221

// NC_MAP_TOWNPORTAL_REQ
// 6170

// NC_MAP_LINKOTHER_CMD
// send zone connection info

// NC_MENU_SERVERMENU_ACK
func ncMenuServerMenuAck(ctx context.Context, np *networking.Parameters) {
	var e playerPromptReplyEvent

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
			Code: errors.ZoneMapNotFound,
			Details: errors.Details{
				"session": session,
			},
		})
		return
	}

	zm.events.send[playerPromptReply] <- &e
}
