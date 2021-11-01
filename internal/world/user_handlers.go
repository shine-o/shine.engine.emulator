package world

import (
	"context"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NcUserLoginWorldReq handles the player's attempt to login to the server
// handle user, given his account
// verify account and character data
// NC_USER_LOGINWORLD_REQ
func ncUserLoginWorldReq(ctx context.Context, np *networking.Parameters) {
	var e serverSelectEvent

	e = serverSelectEvent{
		nc: &structs.NcUserLoginWorldReq{},
		np: np,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)
	if err != nil {
		log.Error(err)
		return
	}

	worldEvents[serverSelect] <- &e
}

// NcUserWillWorldSelectReq handles a petition to return to server select
// NC_USER_WILL_WORLD_SELECT_REQ
func ncUserWillWorldSelectReq(ctx context.Context, np *networking.Parameters) {
	worldEvents[serverSelectToken] <- &serverSelectTokenEvent{
		np: np,
	}
}

// NC_USER_AVATAR_LIST_REQ
// 3103
func ncUserAvatarListReq(ctx context.Context, np *networking.Parameters) {
	var e characterSelectEvent

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = characterSelectEvent{
		np:      np,
		session: session,
	}

	worldEvents[characterSelect] <- &e
}
