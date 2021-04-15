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

// NC_MAP_LOGINCOMPLETE_CMD
// 6147
func ncMapLoginCompleteCmd(ctx context.Context, np *networking.Parameters) {
	// fetch user session
	var (
		pae playerAppearedEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
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
