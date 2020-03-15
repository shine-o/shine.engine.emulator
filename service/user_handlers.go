package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.structs"
)

// handle user, given his account
// verify account and character data
func userLoginWorldReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := &structs.NcUserLoginWorldReq{}
		if err := networking.ReadBinary(pc.Base.Data, nc); err != nil {
			//if err := readBinary(pc.Base.Data, nc); err != nil {
			log.Error(err)
			// TODO: define steps for this kind of errors, either kill the connection or send error code
		} else {
			//go authenticate(ctx, nc)
			// query character data for the user with given id
			// [ future game logic: anything that should be checked before allowing the account to connect to the world ]
			go userLoginWorldAck(ctx, &networking.Command{})
		}
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
		go networking.WriteToClient(ctx, pc)
	}
}

// acknowledge request of login to the world
// send to the client world and character data
func userLoginWorldAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:

		// for now hardcoded... for testing purposes

		pc.Base.OperationCode = 3092

		//nc := &structs.NcUserLoginWorldAck{
		//	WorldManager: 1,
		//	NumOfAvatar:  0,
		//	Avatars:      nil,
		//}
		pc.Base.Data = []byte{1, 0}

		go networking.WriteToClient(ctx, pc)

	}
}
