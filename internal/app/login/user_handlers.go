package login

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
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
		networking.Send(np.OutboundSegments.Send, networking.NC_USER_CLIENT_WRONGVERSION_CHECK_ACK, nil)
		return
	}

	cve = clientVersionEvent{
		nc: &nc,
		np: np,
	}

	loginEvents[clientVersion] <- &cve
}

// NcUserUsLoginReq handles account login
// NC_USER_US_LOGIN_REQ
func ncUserUsLoginReq(ctx context.Context, np *networking.Parameters) {

	var cle credentialsLoginEvent

	//nc := structs.NcUserUsLoginReq{}
	nc := structs.NewUserLoginReq{}

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
