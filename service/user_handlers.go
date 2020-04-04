package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
)

func userClientVersionCheckReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserClientVersionCheckReq{}
		//if err := nc.Unpack(pc.Base.Data); err != nil {
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			log.Error(err)
			go userClientWrongVersionAck(ctx, &networking.Command{})
		} else {

			if _, err := checkClientVersion(ctx, nc); err != nil { // data is irrelevant in this call
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
		pc.Base = networking.CommandBase{
			OperationCode: 3175,
		}
		go pc.Send(ctx)
	}
}

func userClientWrongVersionAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base = networking.CommandBase{
			OperationCode: 3176,
		}
		go pc.Send(ctx)
	}
}

func userUsLoginReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserUsLoginReq{}
		//if err := nc.Unpack(pc.Base.Data); err != nil {
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			go userLoginFailAck(ctx, &networking.Command{})
			return
		}
		pc.NcStruct = &nc
		if err := checkCredentials(ctx, nc); err != nil {
			log.Error(err)
			go userLoginFailAck(ctx, &networking.Command{})
			return
		}
		go userLoginAck(ctx)
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
		nc := structs.NcUserLoginFailAck{
			Err: uint16(69),
		}
		pc.NcStruct = &nc
		go pc.Send(ctx)
	}
}

func userLoginAck(ctx context.Context) {
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

		pc := &networking.Command{
			Base: networking.CommandBase{
				OperationCode: 3082,
			},
		}

		nc, err := serverSelectScreen(ctx)

		if err != nil {
			go unexpectedFailure()
			log.Error(err)
			return
		}
		pc.NcStruct = &nc

		go pc.Send(ctx)
	}
}

func userXtrapReq(ctx context.Context, pc *networking.Command) {}

func userWorldStatusReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		if err := checkWorldStatus(ctx); err != nil {
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
		go pc.Send(ctx)
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
		nc := structs.NcUserWorldSelectReq{}
		//if err := nc.Unpack(pc.Base.Data); err != nil {
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			go unexpectedFailure()
			return
		}
		wci, err := userSelectedServer(ctx)
		if err != nil {
			go unexpectedFailure()
			return
		}
		go userWorldSelectAck(ctx, wci)
	}
}

func userWorldSelectAck(ctx context.Context, wci *lw.WorldConnectionInfo) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserWorldSelectAck{
			WorldStatus: 6,
			Ip:          structs.Name4{},
			Port:        uint16(wci.Port),
		}
		copy(nc.Ip.Name[:], wci.IP)

		pc := &networking.Command{
			Base: networking.CommandBase{
				OperationCode: 3084,
			},
			NcStruct: &nc,
		}
		go pc.Send(ctx)
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
		//if err := restruct.Unpack(pc.Base.Data, binary.LittleEndian, &nc); err != nil {
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			log.Info(err)
			go userLoginFailAck(ctx, &networking.Command{})
		} else {

			pc.NcStruct = &nc
			if err := loginByCode(ctx, nc); err != nil {
				log.Info(err)
				go userLoginFailAck(ctx, &networking.Command{})
			} else {
				go userLoginAck(ctx)
			}
		}
	}
}
