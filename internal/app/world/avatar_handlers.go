package world

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game/character"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
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
		nc:   &nc,
		np:   np,
	}

	worldEvents[createCharacter] <- &cce
}

// NcAvatarCreateSuccAck notifies the character was created and sends the character info
// NC_AVATAR_CREATESUCC_ACK
func ncAvatarCreateSuccAck(np *networking.Parameters, nc *structs.NcAvatarCreateSuccAck) {
	pc := &networking.Command{
		Base: networking.CommandBase{
			OperationCode: 5126,
		},
		NcStruct: nc,
	}
	pc.Send(np.OutboundSegments.Send)
}

// NcAvatarCreateFailAck sends error message to the client when character creation fails
// NC_AVATAR_CREATEFAIL_ACK
func ncAvatarCreateFailAck(np *networking.Parameters, errCode uint16) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 5124,
		},
	}
	pc.NcStruct = &structs.NcAvatarCreateFailAck{
		Err: errCode,
	}
	pc.Send(np.OutboundSegments.Send)
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
		nc:   &nc,
		np:   np,
		done: make(chan bool),
		err:  make(chan error),
	}

	worldEvents[deleteCharacter] <- &dce


}

// AvatarEraseSuccAck notifies the client that the character was deleted successfully
// AVATAR_ERASESUCC_ACK
func avatarEraseSuccAck(np *networking.Parameters, nc *structs.NcAvatarEraseSuccAck) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 5132,
		},
		NcStruct: nc,
	}
	pc.Send(np.OutboundSegments.Send)
}

func createCharErr(np *networking.Parameters, err error) {
	errChar, ok := err.(*character.ErrCharacter)
	if !ok {
		return
	}
	switch errChar.Code {
	case 1:
		ncAvatarCreateFailAck(np, 385)
		return
	}
}
