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

			// check if okay
			//

			//pc.NcStruct = nc
			//wc := &WorldCommand{pc:pc}
			//anc, err := wc.createCharacter(ctx)

			//data, err := structs.Pack(&anc)
			//if err != nil {
			//	log.Error()
			//	return
			//}
			//pc.Base.Data = data
			go avatarCreateSuccAck(ctx, nc)
		}
	}
}

func avatarCreateSuccAck(ctx context.Context, req structs.NcAvatarCreateReq) {
	select {
	case <-ctx.Done():
		return
	default:
		pc := &networking.Command{
			Base:     networking.CommandBase{
				OperationCode: 5126,
			},
		}
		pc.Base.OperationCode = 5126
		nc := structs.NcAvatarCreateSuccAck{}
		nc.NumOfAvatar = 1
		avatarInfo, err := newCharacter(ctx, req)
		if err != nil {
			log.Error(err)
		}

		nc.Avatar = avatarInfo
		pc.NcStruct = &nc
		go pc.Send(ctx)
		//hexData, _ := hex.DecodeString("01ed03000046696768747265726f6f00000000000000000000010000526f754e0000000000000000002b00000085060000ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff0000f0ffffffff0000000000000000000000000000000000000000345632df0000000000000200000000")
		//pc.Base.Data = hexData
	}
}