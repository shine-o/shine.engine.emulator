package service

import (
	"context"
	w "github.com/shine-o/shine.engine.core/grpc/world"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NcUserClientVersionCheckReq handles the client hash verification
// NC_USER_CLIENT_VERSION_CHECK_REQ
func NcUserClientVersionCheckReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserClientVersionCheckReq{}
		//if err := nc.Unpack(pc.Base.Data); err != nil {
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			log.Error(err)
			go NcUserClientWrongVersionCheckAck(ctx)
		} else {
			if _, err := checkClientVersion(ctx, nc); err != nil { // data is irrelevant in this call
				log.Error(err)
				go NcUserClientWrongVersionCheckAck(ctx)
			} else {
				go NcUserClientRightVersionCheckAck(ctx) // will be triggered by method
			}
		}
	}
}

// NcUserClientRightVersionCheckAck acknowledges the client is correct
// NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK
func NcUserClientRightVersionCheckAck(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		pc := &networking.Command{
			Base: networking.CommandBase{
				OperationCode: 3175,
			},
		}
		go pc.Send(ctx)
	}
}

// NcUserClientWrongVersionCheckAck acknowledges the client is incorrect
// NC_USER_CLIENT_WRONGVERSION_CHECK_ACK
func NcUserClientWrongVersionCheckAck(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		pc := &networking.Command{
			Base: networking.CommandBase{
				OperationCode: 3176,
			},
		}
		go pc.Send(ctx)
	}
}

// NcUserUsLoginReq handles account login
// NC_USER_US_LOGIN_REQ
func NcUserUsLoginReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserUsLoginReq{}
		//if err := nc.Unpack(pc.Base.Data); err != nil {
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			go NcUserLoginFailAck(ctx, 69)
			return
		}
		pc.NcStruct = &nc
		if err := checkCredentials(ctx, nc); err != nil {
			log.Error(err)
			go NcUserLoginFailAck(ctx, 69)
			return
		}
		go NcUserLoginAck(ctx)
	}
}

// NcUserLoginFailAck notifies the user about an error while attempting to log in
// NC_USER_LOGINFAIL_ACK
func NcUserLoginFailAck(ctx context.Context, errCode uint16) {
	select {
	case <-ctx.Done():
		return
	default:
		pc := networking.Command{
			Base: networking.CommandBase{
				OperationCode: 3081,
			},
			NcStruct: &structs.NcUserLoginFailAck{
				Err: errCode,
			},
		}
		go pc.Send(ctx)
	}
}

// NcUserLoginAck acknowledges the login attempt and sends the server list to the client
// NC_USER_LOGIN_ACK
func NcUserLoginAck(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		log.Info("Requesting data to the World service")
		unexpectedFailure := func() {
			NcUserLoginFailAck(ctx, 70)
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

// NcUserXtrapReq has no use, but client still sends it
// NC_USER_XTRAP_REQ
func NcUserXtrapReq(ctx context.Context, pc *networking.Command) {}

// NcUserWorldStatusReq pings the world server
// NC_USER_WORLD_STATUS_REQ
func NcUserWorldStatusReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		if err := checkWorldStatus(ctx); err != nil {
			log.Error(err)
			return
		}
		go NcUserWorldStatusAck(ctx, &networking.Command{})
	}
}

// NcUserWorldStatusAck notifies the world is alive
// TODO: different world statuses...
// NC_USER_WORLD_STATUS_ACK
func NcUserWorldStatusAck(ctx context.Context, pc *networking.Command) {
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

// NcUserWorldSelectReq handles a server selected petition
// NC_USER_WORLDSELECT_REQ
func NcUserWorldSelectReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		log.Info("Requesting data to the World service")
		unexpectedFailure := func() {
			NcUserLoginFailAck(ctx, 70)
		}
		nc := structs.NcUserWorldSelectReq{}
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			go unexpectedFailure()
			return
		}
		wci, err := userSelectedServer(ctx, nc)
		if err != nil {
			go unexpectedFailure()
			return
		}
		go NcUserWorldSelectAck(ctx, wci)
	}
}

// NcUserWorldSelectAck queries the world server for connection info and sends it to the clietn
// NC_USER_WORLDSELECT_ACK
func NcUserWorldSelectAck(ctx context.Context, wd *w.WorldData) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserWorldSelectAck{
			WorldStatus: 6,
			Ip:          structs.Name4{},
			Port:        uint16(wd.Port),
		}
		copy(nc.Ip.Name[:], wd.IP)

		pc := &networking.Command{
			Base: networking.CommandBase{
				OperationCode: 3084,
			},
			NcStruct: &nc,
		}
		go pc.Send(ctx)
	}
}

// NcUserNormalLogoutCmd not implemented yet ~~
// NC_USER_NORMALLOGOUT_CMD
func NcUserNormalLogoutCmd(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
	}
}

// NcUserLoginWithOtpReq handles a petition where the client tries a false login
// NC_USER_LOGIN_WITH_OTP_REQ
func NcUserLoginWithOtpReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserLoginWithOtpReq{}
		if err := structs.Unpack(pc.Base.Data, &nc); err != nil {
			log.Info(err)
			go NcUserLoginFailAck(ctx, 70)
		} else {

			pc.NcStruct = &nc
			if err := loginByCode(ctx, nc); err != nil {
				log.Info(err)
				go NcUserLoginFailAck(ctx, 69)
			} else {
				go NcUserLoginAck(ctx)
			}
		}
	}
}
