package world

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// NcUserLoginWorldReq handles the player's attempt to login to the server
// handle user, given his account
// verify account and character data
// NC_USER_LOGINWORLD_REQ
func ncUserLoginWorldReq(ctx context.Context, np *networking.Parameters) {
	var sse serverSelectEvent
	nc := structs.NcUserLoginWorldReq{}
	err := structs.Unpack(np.Command.Base.Data, &nc)
	if err != nil {
		log.Error(err)
		return
	}

	sse = serverSelectEvent{
		nc: &nc,
		np: np,
	}

	worldEvents[serverSelect] <- &sse
}

// NcUserWillWorldSelectReq handles a petition to return to server select
// NC_USER_WILL_WORLD_SELECT_REQ
func ncUserWillWorldSelectReq(ctx context.Context, np *networking.Parameters) {
	sste := serverSelectTokenEvent{
		np: np,
	}
	worldEvents[serverSelectToken] <- &sste
}

// NC_USER_AVATAR_LIST_REQ
// 3103
func ncUserAvatarListReq(ctx context.Context, np *networking.Parameters) {
	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	cs := characterSelectEvent{
		np:      np,
		session: session,
	}

	worldEvents[characterSelect] <- &cs
}
