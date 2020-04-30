package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NcAvatarCreateReq handles character creation request
// NC_AVATAR_CREATE_REQ
func NcAvatarCreateReq(ctx context.Context, np * networking.Parameters) {
	nc := structs.NcAvatarCreateReq{}
	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		return
	}

	nc = structs.NcAvatarCreateReq{}

	err := structs.Unpack(np.Command.Base.Data, &nc)
	if err != nil {
		log.Error(err)
		return
	}
	wsi := ctx.Value(networking.ShineSession)
	ws := wsi.(*session)

	err = character.Validate(db, ws.UserID, nc)
	if err != nil {
		go createCharErr(ctx, err)
		return
	}

	ai, err := character.New(db, ws.UserID, nc)
	if err != nil {
		log.Error(err)
		return
	}
	go NcAvatarCreateSuccAck(ctx, ai)
}

// NcAvatarCreateFailAck sends error message to the client when character creation fails
// NC_AVATAR_CREATEFAIL_ACK
func NcAvatarCreateFailAck(ctx context.Context, errCode uint16) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 5124,
		},
	}
	pc.NcStruct = &structs.NcAvatarCreateFailAck{
		Err: errCode,
	}
	pc.Send(ctx)
}

// NcAvatarCreateSuccAck notifies the character was created and sends the character info
// NC_AVATAR_CREATESUCC_ACK
func NcAvatarCreateSuccAck(ctx context.Context, ai structs.AvatarInformation) {
	pc := &networking.Command{
		Base: networking.CommandBase{
			OperationCode: 5126,
		},
	}
	nc := structs.NcAvatarCreateSuccAck{}
	nc.NumOfAvatar = 1 //
	nc.Avatar = ai
	pc.NcStruct = &nc
	pc.Send(ctx)
}

// NcAvatarEraseReq handles a petition to delete a character
// NC_AVATAR_ERASE_REQ
func NcAvatarEraseReq(ctx context.Context, np * networking.Parameters) {
	nc := structs.NcAvatarEraseReq{}
	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		// todo: error nc if possible
		log.Error(err)
		return
	}
	wsi := ctx.Value(networking.ShineSession)
	ws := wsi.(*session)

	err := character.Delete(db, ws.UserID, nc)
	if err != nil {
		log.Error(err)
		// todo: error nc if possible
		return
	}
	go AvatarEraseSuccAck(ctx, structs.NcAvatarEraseSuccAck{
		Slot: nc.Slot,
	})
}

// AvatarEraseSuccAck notifies the client that the character was deleted successfully
// AVATAR_ERASESUCC_ACK
func AvatarEraseSuccAck(ctx context.Context, ack structs.NcAvatarEraseSuccAck) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 5132,
		},
	}
	pc.NcStruct = &ack
	pc.Send(ctx)
}

func createCharErr(ctx context.Context, err error) {
	log.Error(err)
	errChar, ok := err.(*character.ErrCharacter)
	if !ok {
		return
	}
	switch errChar.Code {
	case 1:
		go NcAvatarCreateFailAck(ctx, 385)
		return
	}
}
