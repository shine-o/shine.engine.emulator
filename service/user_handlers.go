package service

import (
	"bytes"
	"context"
	"encoding/binary"
	networking "github.com/shine-o/shine.engine.networking"
	gs "shine.engine.game_structs"
)

// handle user, given his account
// verify account and character data
func userLoginWorldReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := &gs.NcUserLoginWorldReq{}
		//if err := networking.ReadBinary(pc.Base.Data, nc); err != nil {
		if err := readBinary(pc.Base.Data, nc); err != nil {
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

func readBinary(data []byte, nc interface{}) error {
	// data is only 83. it reads all 83 bytes into the struct
	// then tries again, since there's still 16 bytes to read
	// but of course, there is nothing to read from.
	// so. if struct is
	//fmt.Println(binary.Size(loginAck))
	//structSize := binary.Size(loginAck)
	//buffer := make([]byte, structSize)
	//actualData := b[3:]
	//copy(buffer, actualData)
	structSize := binary.Size(nc)
	buffer := make([]byte, structSize)
	copy(buffer, data)

	buf := bytes.NewBuffer(buffer)
	if err := binary.Read(buf, binary.LittleEndian, nc); err != nil {
		log.Error(err)
		return err
	}
	return nil
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

		//nc := &gs.NcUserLoginWorldAck{
		//	WorldManager: 1,
		//	NumOfAvatar:  0,
		//	Avatars:      nil,
		//}
		pc.Base.Data = []byte{1, 0}

		go networking.WriteToClient(ctx, pc)

	}
}
