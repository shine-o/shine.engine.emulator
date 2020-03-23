package service

import (
	"context"
	"encoding/binary"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
	"gopkg.in/restruct.v1"
)

// handle user, given his account
// verify account and character data
func userLoginWorldReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserLoginWorldReq{}
		if err := restruct.Unpack(pc.Base.Data, binary.LittleEndian, &nc); err != nil {
			log.Error(err)
			// TODO: define steps for this kind of errors, either kill the connection or send error code
		} else {
			pc.NcStruct = nc
			wc := WorldCommand{pc: pc}
			if err := wc.loginToWorld(ctx); err != nil {
				log.Error(err)
				return
			}
			go userLoginWorldAck(ctx, &networking.Command{})
		}
	}
}

// acknowledge request of login to the world
// send to the client world and character data
func userLoginWorldAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base.OperationCode = 3092
		wc := &WorldCommand{pc: pc}

		data, err := wc.userWorldInfo(ctx)
		if err != nil {
			return
		}
		pc.Base.Data = data
		go networking.WriteToClient(ctx, pc)
	}
}

func userWillWorldSelectReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		go userWillWorldSelectAck(ctx, &networking.Command{})
	}
}

func userWillWorldSelectAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base.OperationCode = 3124
		wc := &WorldCommand{pc: pc}

		nc, err :=  wc.returnToServerSelect(ctx)
		if err != nil {
			return
		}

		data, err := structs.Pack(&nc)
		if err != nil {
			return
		}

		pc.Base.Data = data
		go networking.WriteToClient(ctx, pc)
	}
}

//func userLoginWorldFailAck(ctx context.Context, pc *networking.Command) {
//	select {
//	case <-ctx.Done():
//		return
//	default:
//
//	}
//}
