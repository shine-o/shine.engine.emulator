package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
	"github.com/shine-o/shine.engine.structs"
	"time"
)

const gRPCTimeout = time.Second * 2

func userClientVersionCheckReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserClientVersionCheckReq{}
		if err := networking.ReadBinary(pc.Base.Data, &nc); err != nil {
			log.Error(err)
			go userClientWrongVersionAck(ctx, &networking.Command{})
		} else {
			pc.NcStruct = nc
			lc := &LoginCommand{pc: pc}

			if _, err := lc.checkClientVersion(ctx); err != nil { // data is irrelevant in this call
				log.Error(err)
				go userClientWrongVersionAck(ctx, &networking.Command{})
			} else {
				go userClientVersionCheckAck(ctx, &networking.Command{}) // will be triggered by method
			}
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

func userClientWrongVersionAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		base := networking.CommandBase{}
		base.OperationCode = 3176
		pc.Base = base
		go networking.WriteToClient(ctx, pc)
	}
}

func userUsLoginReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserUsLoginReq{}
		if err := networking.ReadBinary(pc.Base.Data, &nc); err != nil {
			go userLoginFailAck(ctx, &networking.Command{})
		} else {
			pc.NcStruct = nc
			lc := &LoginCommand{pc: pc}
			if err := lc.checkCredentials(ctx); err != nil {
				log.Error(err)
				go userLoginFailAck(ctx, &networking.Command{})
			} else {
				go userLoginAck(ctx, &networking.Command{})
			}
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
		nc := &structs.NcUserLoginFailAck{
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
		log.Info("Requesting data to the World service")
		unexpectedFailure := func() {
			userLoginFailAck(ctx, &networking.Command{
				NcStruct: &structs.NcUserLoginFailAck{
					Err: 70,
				},
			})
		}

		pc.Base = networking.CommandBase{
			OperationCode: 3082,
		}

		grpcc.mu.Lock()
		conn := grpcc.services["world"]
		c := lw.NewWorldClient(conn)
		grpcc.mu.Unlock()

		rpcCtx, _ := context.WithTimeout(context.Background(), gRPCTimeout)
		if r, err := c.AvailableWorlds(rpcCtx, &lw.ClientMetadata{
			Ip: "127.0.0.01",
		}); err != nil {
			log.Error(err)
			go unexpectedFailure()
		} else {
			pc.Base.Data = r.Info
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
		lc := &LoginCommand{pc: pc}
		if err := lc.checkWorldStatus(ctx); err != nil {
			return
		}
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

func userWorldSelectReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		log.Info("Requesting data to the World service")
		unexpectedFailure := func() {
			userLoginFailAck(ctx, &networking.Command{
				NcStruct: &structs.NcUserLoginFailAck{
					Err: 70,
				},
			})
		}
		nc := &structs.NcUserWorldSelectReq{}
		if err := networking.ReadBinary(pc.Base.Data, nc); err != nil {
			go unexpectedFailure()
		} else {

			lc := &LoginCommand{pc: pc}
			data, err := lc.userSelectedServer(ctx)
			if err != nil {
				return
			}
			go userWorldSelectAck(ctx, &networking.Command{
				Base: networking.CommandBase{
					Data: data,
				},
			})
		}
	}
}

func userWorldSelectAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base.OperationCode = 3084
		go networking.WriteToClient(ctx, pc)
	}
}

func userNormalLogoutCmd(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
	}
}

func userLoginWithOtpReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserLoginWithOtpReq{}
		if err := networking.ReadBinary(pc.Base.Data, &nc); err != nil {
			log.Info(err)
			go userLoginFailAck(ctx, &networking.Command{})
		} else {

			pc.NcStruct = nc
			lc := &LoginCommand{pc: pc}

			if err := lc.loginByCode(ctx); err != nil {
				log.Info(err)
				go userLoginFailAck(ctx, &networking.Command{})
			} else {
				go userLoginAck(ctx, &networking.Command{})
			}
		}
	}
}
