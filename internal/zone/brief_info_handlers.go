package zone

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NC_BRIEFINFO_INFORM_CMD
func ncBriefInfoInformCmd(ctx context.Context, np *networking.Parameters) {
	// trigger handleInfo
	// if targetHandle is within range of affectedHandle
	//		send NC_BRIEFINFO_LOGINCHARACTER_CMD of the targetHandle to the affectedHandle
	var (
		e unknownHandleEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = unknownHandleEvent{
		handle: session.handle,
		nc:     &structs.NcBriefInfoInformCmd{},
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
			Details: errors.ErrDetails{
				"session": session,
			},
		})
		return
	}

	zm.send[unknownHandle] <- &e
}
