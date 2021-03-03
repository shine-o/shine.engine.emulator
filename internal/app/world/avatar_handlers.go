package world

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// NcAvatarCreateReq handles character creation request
// NC_AVATAR_CREATE_REQ
func ncAvatarCreateReq(ctx context.Context, np *networking.Parameters) {
	var (
		nc  structs.NcAvatarCreateReq
		cce createCharacterEvent
	)

	err := structs.Unpack(np.Command.Base.Data, &nc)
	if err != nil {
		log.Error(err)
		return
	}

	cce = createCharacterEvent{
		nc: &nc,
		np: np,
	}

	worldEvents[createCharacter] <- &cce
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
	var (
		nc  structs.NcAvatarEraseReq
		dce deleteCharacterEvent
	)

	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		// todo: error nc if possible
		log.Error(err)
		return
	}

	dce = deleteCharacterEvent{
		nc: &nc,
		np: np,
	}

	worldEvents[deleteCharacter] <- &dce

}

func createCharErr(np *networking.Parameters, err error) {
	errChar, ok := err.(*persistence.ErrCharacter)
	if !ok {
		return
	}
	switch errChar.Code {
	case 1:
		ncAvatarCreateFailAck(np, 385)
		return
	}
}
