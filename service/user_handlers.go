package service

import (
	"context"
	networking "github.com/shine-o/shine.engine.networking"
)

func userClientVersionCheckReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:

		nc := &ncUserClientVersionCheckReq{}

		if err := networking.ReadBinary(pc.Base.Data, nc); err != nil {
			// TODO: define steps for this kind of errors, either kill the connection or send error code

		} else {
			// method for checking version
			go userClientVersionCheckAck(ctx, &networking.Command{}) // will be triggered by method
		}
	}
}

func userClientVersionCheckAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		base := networking.CommandBase{}
		base.OperationCode = 3175
		pc.Base = base
		go networking.WriteToClient(ctx, pc)
	}
}

func userUsLoginReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := &ncUserUsLoginReq{}
		if err := networking.ReadBinary(pc.Base.Data, nc); err != nil {
			// TODO: define steps for this kind of errors, either kill the connection or send error code
		} else {
			go nc.authenticate(ctx)
		}
	}
}

func userLoginFailAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base = networking.CommandBase{
			OperationCode: 3081,
		}

		// 090c 4500
		nc := &ncUserLoginFailAck{
			Err: uint16(69),
		}

		if data, err := networking.WriteBinary(nc); err != nil {

		} else {
			pc.Base.Data = data
			go networking.WriteToClient(ctx, pc)
		}
	}
}

func userLoginAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base = networking.CommandBase{
			OperationCode: 3082,
		}
		nc := &ncUserLoginAck{}
		nc.setServerInfo(ctx)

		if data, err := networking.WriteBinary(nc); err == nil {
			pc.Base.Data = data
			go networking.WriteToClient(ctx, pc)
		}
	}
}

func userXtrapReq(ctx context.Context, pc *networking.Command) {}

func userXtrapAck(ctx context.Context, pc *networking.Command) {}

func userWorldStatusReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		// ping World service for status :)
		go userWorldStatusAck(ctx, &networking.Command{})
	}
}

func userWorldStatusAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base = networking.CommandBase{
			OperationCode: 3100,
		}
		go networking.WriteToClient(ctx, pc)
	}
}

func userWorldSelectAck(ctx context.Context, pc *networking.Command) {}

func userWorldSelectReq(ctx context.Context, pc *networking.Command) {}

func userNormalLogoutCmd(ctx context.Context, pc *networking.Command) {}
