package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
)

// handle user, given his account
// verify account and character data
func userLoginWorldReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserLoginWorldReq{}
		if err := nc.Unpack(pc.Base.Data); err != nil {
			log.Error(err)
			// TODO: define steps for this kind of errors, either kill the connection or send error code
		} else {
			pc.NcStruct = &nc
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

		nc, err := wc.userWorldInfo(ctx)
		if err != nil {
			return
		}
		pc.NcStruct = &nc
		go pc.Send(ctx)
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
		pc.NcStruct = &nc
		go pc.Send(ctx)
		//
		//data, err := nc.Pack()
		//if err != nil {
		//	return
		//}
		//pc.Base.Data = data
		//go networking.WriteToClient(ctx, pc)
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
