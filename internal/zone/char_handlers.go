package zone

import (
	"context"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NC_MAP_LOGIN_REQ
func ncMapLoginReq(ctx context.Context, np *networking.Parameters) {
	var (
		nc structs.NcMapLoginReq
		e  playerMapLoginEvent
	)

	err := structs.Unpack(np.Command.Base.Data, &nc)
	if err != nil {
		log.Error(err)
		return
	}

	e = playerMapLoginEvent{
		nc: nc,
		np: np,
	}

	zoneEvents[playerMapLogin] <- &e
}

// NC_MAP_LOGINCOMPLETE_CMD
func ncMapLoginCompleteCmd(ctx context.Context, np *networking.Parameters) {
	// fetch user session
	var (
		e playerAppearedEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	zm, ok := maps.list[session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}

	e = playerAppearedEvent{
		handle: session.handle,
	}

	zm.events.send[playerAppeared] <- &e
}

// NC_CHAR_LOGOUTCANCEL_CMD
func ncCharLogoutCancelCmd(ctx context.Context, np *networking.Parameters) {
	var e playerLogoutCancelEvent

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = playerLogoutCancelEvent{
		sessionID: session.id,
	}

	zoneEvents[playerLogoutCancel] <- &e
}

// NC_CHAR_LOGOUTREADY_CMD
func ncCharLogoutReadyCmd(ctx context.Context, np *networking.Parameters) {
	var plse playerLogoutStartEvent

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
