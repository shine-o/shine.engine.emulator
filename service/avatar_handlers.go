package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

func avatarCreateReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcAvatarCreateReq{}
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			return
		}

		nc = structs.NcAvatarCreateReq{}

		err := structs.Unpack(pc.Base.Data, &nc)
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
		go avatarCreateSuccAck(ctx, ai)
	}
}

func ncAvatarCreateFailAck(ctx context.Context, errCode uint16) {
	select {
	case <-ctx.Done():
		return
	default:
		pc := networking.Command{
			Base: networking.CommandBase{
				OperationCode: 5124,
			},
		}
		pc.NcStruct = &structs.NcAvatarCreateFailAck{
			Err: errCode,
		}
		go pc.Send(ctx)
	}
}

func avatarCreateSuccAck(ctx context.Context, ai structs.AvatarInformation) {
	select {
	case <-ctx.Done():
		return
	default:
		pc := &networking.Command{
			Base: networking.CommandBase{
				OperationCode: 5126,
			},
		}
		nc := structs.NcAvatarCreateSuccAck{}
		nc.NumOfAvatar = 1 //
		nc.Avatar = ai
		pc.NcStruct = &nc
		go pc.Send(ctx)
	}
}

func avatarEraseReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcAvatarEraseReq{}
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
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
		go avatarEraseSuccAck(ctx, structs.NcAvatarEraseSuccAck{
			Slot: nc.Slot,
		})
	}
}

func avatarEraseSuccAck(ctx context.Context, ack structs.NcAvatarEraseSuccAck) {
	select {
	case <-ctx.Done():
		return
	default:
		pc := networking.Command{
			Base: networking.CommandBase{
				OperationCode: 5132,
			},
		}
		pc.NcStruct = &ack
		go pc.Send(ctx)
	}
}

func createCharErr(ctx context.Context, err error) {
	log.Error(err)
	errChar, ok := err.(*character.ErrCharacter)
	if !ok {
		return
	}
	switch errChar.Code {
	case 1:
		go ncAvatarCreateFailAck(ctx, 385)
		return
	}
}
