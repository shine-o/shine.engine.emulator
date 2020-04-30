package service

import (
	"context"
	w "github.com/shine-o/shine.engine.core/grpc/world"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NcUserClientVersionCheckReq handles the client hash verification
// NC_USER_CLIENT_VERSION_CHECK_REQ
func NcUserClientVersionCheckReq(ctx context.Context, np *networking.Parameters) {
	nc := structs.NcUserClientVersionCheckReq{}
	//if err := nc.Unpack(pc.Base.Data); err != nil {
	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		log.Error(err)
		go NcUserClientWrongVersionCheckAck(ctx)
	} else {
		if _, err := checkClientVersion(nc); err != nil { // data is irrelevant in this call
			log.Error(err)
			go NcUserClientWrongVersionCheckAck(ctx)
		} else {
			go NcUserClientRightVersionCheckAck(ctx) // will be triggered by method
		}
	}
}

// NcUserClientRightVersionCheckAck acknowledges the client is correct
// NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK
func NcUserClientRightVersionCheckAck(ctx context.Context) {
	pc := &networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3175,
		},
	}
	pc.Send(ctx)
}

// NcUserClientWrongVersionCheckAck acknowledges the client is incorrect
// NC_USER_CLIENT_WRONGVERSION_CHECK_ACK
func NcUserClientWrongVersionCheckAck(ctx context.Context) {
	pc := &networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3176,
		},
	}
	pc.Send(ctx)
}

// NcUserUsLoginReq handles account login
// NC_USER_US_LOGIN_REQ
func NcUserUsLoginReq(ctx context.Context, np *networking.Parameters) {
	nc := structs.NcUserUsLoginReq{}
	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		go NcUserLoginFailAck(ctx, 69)
		return
	}
	np.Command.NcStruct = &nc
	if err := checkCredentials(nc); err != nil {
		log.Error(err)
		go NcUserLoginFailAck(ctx, 69)
		return
	}
	go NcUserLoginAck(ctx)
}

// NcUserLoginFailAck notifies the user about an error while attempting to log in
// NC_USER_LOGINFAIL_ACK
func NcUserLoginFailAck(ctx context.Context, errCode uint16) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3081,
		},
		NcStruct: &structs.NcUserLoginFailAck{
			Err: errCode,
		},
	}
	pc.Send(ctx)
}

// NcUserLoginAck acknowledges the login attempt and sends the server list to the client
// NC_USER_LOGIN_ACK
func NcUserLoginAck(ctx context.Context) {
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
	pc.Send(ctx)
}

// NcUserXtrapReq has no use, but client still sends it
// NC_USER_XTRAP_REQ
func NcUserXtrapReq(ctx context.Context, np *networking.Parameters) {}

// NcUserWorldStatusReq pings the world server
// NC_USER_WORLD_STATUS_REQ
func NcUserWorldStatusReq(ctx context.Context, np *networking.Parameters) {
	if err := checkWorldStatus(); err != nil {
		log.Error(err)
		return
	}
	go NcUserWorldStatusAck(ctx, &networking.Command{})
}

// NcUserWorldStatusAck notifies the world is alive
// TODO: different world statuses...
// NC_USER_WORLD_STATUS_ACK
func NcUserWorldStatusAck(ctx context.Context, pc * networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base = networking.CommandBase{
			OperationCode: 3100,
		}
		pc.Send(ctx)
	}
}

// NcUserWorldSelectReq handles a server selected petition
// NC_USER_WORLDSELECT_REQ
func NcUserWorldSelectReq(ctx context.Context, np *networking.Parameters) {
	log.Info("Requesting data to the World service")
	unexpectedFailure := func() {
		NcUserLoginFailAck(ctx, 70)
	}
	nc := structs.NcUserWorldSelectReq{}
	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
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

// NcUserWorldSelectAck queries the world server for connection info and sends it to the clietn
// NC_USER_WORLDSELECT_ACK
func NcUserWorldSelectAck(ctx context.Context, wd *w.WorldData) {
	nc := structs.NcUserWorldSelectAck{
		WorldStatus: 6,
		Ip:          structs.Name4{
			Name: wd.IP,
		},
		Port:        uint16(wd.Port),
	}

	pc := &networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3084,
		},
		NcStruct: &nc,
	}
	pc.Send(ctx)
}

// NcUserNormalLogoutCmd not implemented yet ~~
// NC_USER_NORMALLOGOUT_CMD
func NcUserNormalLogoutCmd(ctx context.Context, np *networking.Parameters) {
	select {
	case <-ctx.Done():
		return
	default:
	}
}

// NcUserLoginWithOtpReq handles a petition where the client tries a false login
// NC_USER_LOGIN_WITH_OTP_REQ
func NcUserLoginWithOtpReq(ctx context.Context, np *networking.Parameters) {
	nc := structs.NcUserLoginWithOtpReq{}
	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		log.Info(err)
		go NcUserLoginFailAck(ctx, 70)
	} else {
		np.Command.NcStruct = &nc
		if err := loginByCode(nc); err != nil {
			log.Error(err)
			go NcUserLoginFailAck(ctx, 69)
		} else {
			go NcUserLoginAck(ctx)
		}
	}
}
