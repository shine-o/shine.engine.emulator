package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NcUserLoginWorldReq handles the player's attempt to login to the server
// handle user, given his account
// verify account and character data
// NC_USER_LOGINWORLD_REQ
func NcUserLoginWorldReq(ctx context.Context, np * networking.Parameters) {
	nc := structs.NcUserLoginWorldReq{}
	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		log.Error(err)
		// TODO: define steps for this kind of errors, either kill the connection or send error code
	} else {
		np.Command.NcStruct = &nc
		if err := loginToWorld(ctx, nc); err != nil {
			log.Error(err)
			return
		}
		go NcUserLoginWorldAck(ctx, &networking.Command{})
	}
}

// NcUserLoginWorldAck acknowledges the player's attempt to login, returning player character information
// acknowledge request of login to the service
// send to the client service and character data
// NC_USER_LOGINWORLD_ACK
func NcUserLoginWorldAck(ctx context.Context, c * networking.Command) {
	c.Base.OperationCode = 3092
	nc, err := userWorldInfo(ctx)
	if err != nil {
		log.Error(err)
		return
	}
	c.NcStruct = &nc
	c.Send(ctx)
}

// NcUserWillWorldSelectReq handles a petition to return to server select
// NC_USER_WILL_WORLD_SELECT_REQ
func NcUserWillWorldSelectReq(ctx context.Context, np * networking.Parameters) {
	go NcUserWillWorldSelectAck(ctx)
}

// NcUserWillWorldSelectAck acknowledges a petition to return to server select
// NcUserWillWorldSelectAck
func NcUserWillWorldSelectAck(ctx context.Context) {
	pc := &networking.Command{
		Base: networking.CommandBase{
			OperationCode: 3124,
		},
	}
	nc, err := returnToServerSelect()
	if err != nil {
		return
	}
	pc.NcStruct = &nc
	pc.Send(ctx)
}

//NC_USER_NORMALLOGOUT_CMD
func NcUserNormalLogoutCmd(ctx context.Context, np * networking.Parameters) {
	np.NetVars.CloseConnection <- true
}

//func userLoginWorldFailAck(ctx context.Context, pc *networking.Command) {
//	select {
//	case <-ctx.Done():
//		return
//	default:
//
//	}
//}
