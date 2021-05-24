package world

import (
	"context"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NcAvatarCreateReq handles character creation request
// NC_AVATAR_CREATE_REQ
func ncAvatarCreateReq(ctx context.Context, np *networking.Parameters) {
	var e createCharacterEvent

	e = createCharacterEvent{
		nc: &structs.NcAvatarCreateReq{},
		np: np,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)
	if err != nil {
		log.Error(err)
		return
	}

	worldEvents[createCharacter] <- &e
}

// NcAvatarCreateFailAck sends error message to the client when character creation fails
// NC_AVATAR_CREATEFAIL_ACK
func ncAvatarCreateFailAck(np *networking.Parameters, errCode uint16) {
	nc := &structs.NcAvatarCreateFailAck{
		Err: errCode,
	}
	networking.Send(np.OutboundSegments.Send, networking.NC_AVATAR_CREATEFAIL_ACK, nc)
}

// NcAvatarEraseReq handles a petition to delete a character
// NC_AVATAR_ERASE_REQ
func ncAvatarEraseReq(ctx context.Context, np *networking.Parameters) {
	var e deleteCharacterEvent

	e = deleteCharacterEvent{
		nc: &structs.NcAvatarEraseReq{},
		np: np,
	}

	if err := structs.Unpack(np.Command.Base.Data, e.nc); err != nil {
		// todo: error nc if possible
		log.Error(err)
		return
	}

	worldEvents[deleteCharacter] <- &e
}
