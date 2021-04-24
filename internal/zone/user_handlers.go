package zone

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
)

// NC_USER_NORMALLOGOUT_CMD
func ncUserNormalLogoutCmd(ctx context.Context, np *networking.Parameters) {
	var (
		e playerLogoutConcludeEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = playerLogoutConcludeEvent{
		sessionID: session.id,
		err:       make(chan error),
	}

	zoneEvents[playerLogoutConclude] <- &e

	select {
	case e := <-e.err:
		log.Error(e)
	}
}
