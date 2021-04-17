package zone

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NC_BAT_TARGETTING_REQ
// 9217
func ncBatTargetingReq(ctx context.Context, np *networking.Parameters) {
	var (
		e playerSelectsEntityEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = playerSelectsEntityEvent{
		nc:     &structs.NcBatTargetInfoReq{},
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

	zm.send[playerSelectsEntity] <- &e
}

func ncBatUntargetReq(ctx context.Context, np *networking.Parameters) {

	// remove current SelectionOrder value for player
	var (
		e playerUnselectsEntityEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = playerUnselectsEntityEvent{
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

	zm.send[playerUnselectsEntity] <- &e
}