package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// handle user, given his account
// verify account and character data
func userLoginWorldReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserLoginWorldReq{}
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			log.Error(err)
			// TODO: define steps for this kind of errors, either kill the connection or send error code
		} else {
			pc.NcStruct = &nc
			if err := loginToWorld(ctx, nc); err != nil {
				log.Error(err)
				return
			}
			go userLoginWorldAck(ctx, &networking.Command{})
		}
	}
}

// acknowledge request of login to the service
// send to the client service and character data
func userLoginWorldAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base.OperationCode = 3092

		nc, err := userWorldInfo(ctx)
		if err != nil {
			log.Error(err)
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
		go userWillWorldSelectAck(ctx)
	}
}

func userWillWorldSelectAck(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		pc := &networking.Command{
			Base: networking.CommandBase{
				OperationCode: 3124,
			},
			NcStruct: nil,
		}
		nc, err := returnToServerSelect(ctx)
		if err != nil {
			return
		}
		pc.NcStruct = &nc
		go pc.Send(ctx)
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
