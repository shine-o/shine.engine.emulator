package login

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// NcUserClientVersionCheckReq handles the client hash verification
// NC_USER_CLIENT_VERSION_CHECK_REQ
func ncUserClientVersionCheckReq(ctx context.Context, np *networking.Parameters) {
	var (
		nc  structs.NcUserClientVersionCheckReq
		cve clientVersionEvent
	)

	err := structs.Unpack(np.Command.Base.Data, &nc)
	if err != nil {
		log.Error(err)
		go ncUserClientWrongVersionCheckAck(np)
		return
	}

	cve = clientVersionEvent{
		nc: &nc,
		np: np,
	}

	loginEvents[clientVersion] <- &cve
}

// NcUserClientRightVersionCheckAck acknowledges the client is correct
// NC_USER_CLIENT_RIGHTVERSION_CHECK_ACK
func ncUserClientRightVersionCheckAck(np *networking.Parameters) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3175,
		},
	}
	pc.Send(np.OutboundSegments.Send)
}

// NcUserClientWrongVersionCheckAck acknowledges the client is incorrect
// NC_USER_CLIENT_WRONGVERSION_CHECK_ACK
func ncUserClientWrongVersionCheckAck(np *networking.Parameters) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3176,
		},
	}
	pc.Send(np.OutboundSegments.Send)
}

// NcUserUsLoginReq handles account login
// NC_USER_US_LOGIN_REQ
func ncUserUsLoginReq(ctx context.Context, np *networking.Parameters) {

	var cle credentialsLoginEvent

	nc := structs.NcUserUsLoginReq{}

	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		ncUserLoginFailAck(np, 69)
		return
	}

	cle = credentialsLoginEvent{
		nc: &nc,
		np: np,
	}

	loginEvents[credentialsLogin] <- &cle
}

// NcUserLoginFailAck notifies the user about an error while attempting to log in
// NC_USER_LOGINFAIL_ACK
func ncUserLoginFailAck(np *networking.Parameters, errCode uint16) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3081,
		},
		NcStruct: &structs.NcUserLoginFailAck{
			Err: errCode,
		},
	}
	pc.Send(np.OutboundSegments.Send)
}

// NcUserLoginAck acknowledges the login attempt and sends the server list to the client
// NC_USER_LOGIN_ACK
func ncUserLoginAck(np *networking.Parameters, nc structs.NcUserLoginAck) {
	pc := &networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3082,
		},
		NcStruct: &nc,
	}
	pc.Send(np.OutboundSegments.Send)
}

// NcUserXtrapReq has no use, but client still sends it
// NC_USER_XTRAP_REQ
func ncUserXtrapReq(ctx context.Context, np *networking.Parameters) {}

// NcUserWorldStatusReq pings the world server
// NC_USER_WORLD_STATUS_REQ
func ncUserWorldStatusReq(ctx context.Context, np *networking.Parameters) {
	loginEvents[worldManagerStatus] <- &worldManagerStatusEvent{
		np: np,
	}
}

// NcUserWorldStatusAck notifies the world is alive
// TODO: different world statuses...
// NC_USER_WORLD_STATUS_ACK
func ncUserWorldStatusAck(np *networking.Parameters) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3100,
		},
	}
	pc.Send(np.OutboundSegments.Send)
}

// NcUserWorldSelectReq handles a server selected petition
// NC_USER_WORLDSELECT_REQ
func ncUserWorldSelectReq(ctx context.Context, np *networking.Parameters) {
	var sse serverSelectEvent

	unexpectedFailure := func() {
		ncUserLoginFailAck(np, 70)
	}

	nc := structs.NcUserWorldSelectReq{}
	err := structs.Unpack(np.Command.Base.Data, &nc)

	if err != nil {
		log.Error(err)
		unexpectedFailure()
		return
	}

	sse = serverSelectEvent{
		nc: &nc,
		np: np,
	}

	loginEvents[serverSelect] <- &sse
}

// NcUserWorldSelectAck queries the world server for connection info and sends it to the client
// NC_USER_WORLDSELECT_ACK
func ncUserWorldSelectAck(np *networking.Parameters, ack structs.NcUserWorldSelectAck) {
	pc := &networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3084,
		},
		NcStruct: &ack,
	}
	pc.Send(np.OutboundSegments.Send)
}

// NcUserNormalLogoutCmd not implemented yet ~~
// NC_USER_NORMALLOGOUT_CMD
func ncUserNormalLogoutCmd(ctx context.Context, np *networking.Parameters) {

}

// NcUserLoginWithOtpReq handles a petition where the client tries a false login
// NC_USER_LOGIN_WITH_OTP_REQ
func ncUserLoginWithOtpReq(ctx context.Context, np *networking.Parameters) {
	var tl tokenLoginEvent

	nc := structs.NcUserLoginWithOtpReq{}
	err := structs.Unpack(np.Command.Base.Data, &nc)
	if err != nil {
		log.Info(err)
		ncUserLoginFailAck(np, 70)
		return
	}

	tl = tokenLoginEvent{
		nc: &nc,
		np: np,
	}

	loginEvents[tokenLogin] <- &tl

}
