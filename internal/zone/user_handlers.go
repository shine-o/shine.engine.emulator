package zone

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
)

//3096
// NC_USER_NORMALLOGOUT_CMD
func ncUserNormalLogoutCmd(ctx context.Context, np *networking.Parameters) {
	var (
		plce playerLogoutConcludeEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	plce = playerLogoutConcludeEvent{
		sessionID: session.id,
		err:       make(chan error),
	}

	zoneEvents[playerLogoutConclude] <- &plce

	select {
	case e := <-plce.err:
		log.Error(e)
	}
}
