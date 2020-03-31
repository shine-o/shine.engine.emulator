package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
)

func avatarCreateReq(ctx context.Context, pc *networking.Command)   {
	select {
	case <- ctx.Done():
		return
	default:
		nc := structs.NcAvatarCreateReq{}
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			return
		} else {
			nc := structs.NcAvatarCreateReq{}

			err := nc.Unpack(pc.Base.Data)
			if err != nil {
				log.Error(err)
				return
			}

			err = validateCharacter(ctx, nc)
			if err != nil {
				log.Error(err)
				return
			}

			ai, err := newCharacter(ctx, nc)
			if err != nil {
				log.Error(err)
				return
			}
			go avatarCreateSuccAck(ctx, ai)
		}
	}
}

func avatarCreateSuccAck(ctx context.Context, ai structs.AvatarInformation) {
	select {
	case <-ctx.Done():
		return
	default:
		pc := &networking.Command{
			Base:     networking.CommandBase{
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
		err := deleteCharacter(ctx, nc)
		if err != nil {
			log.Error(err)
			// todo: error nc if possible
			return
		}
		go avatarEraseSuccAck(ctx, structs.NcAvatarEraseSuccAck{
			Slot:nc.Slot,
		})
	}
}

func avatarEraseSuccAck(ctx context.Context, ack structs.NcAvatarEraseSuccAck) {
	select {
	case <-ctx.Done():
		return
	default:
		pc := networking.Command{
			Base:     networking.CommandBase{
				OperationCode:		5132,
			},
		}
		pc.NcStruct = &ack
		go pc.Send(ctx)
	}
}