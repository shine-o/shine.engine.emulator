package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.structs"
	"time"
)

// handle user, given his account
// verify account and character data
func userLoginWorldReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		nc := structs.NcUserLoginWorldReq{}
		if err := networking.ReadBinary(pc.Base.Data, &nc); err != nil {
			log.Error(err)
			// TODO: define steps for this kind of errors, either kill the connection or send error code
		} else {
			pc.NcStruct = nc
			wc := WorldCommand{pc:pc}
			if err := wc.loginToWorld(ctx); err != nil {
				log.Error(err)
				return
			}
			go userLoginWorldAck(ctx, &networking.Command{})
		}
	}
}

// acknowledge request of login to the world
// send to the client world and character data
func userLoginWorldAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base.OperationCode = 3092
		wc := &WorldCommand{pc:pc}

		data, err := wc.userWorldInfo(ctx)
		if err != nil {
			return
		}
		pc.Base.Data = data
		go networking.WriteToClient(ctx, pc)
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
		otp := randStringBytesMaskImprSrcUnsafe(32)
		//otp := "a85472c3841de5fc22433560fe32a2a3"
		err := redisClient.Set(otp, otp, 20 * time.Second).Err()
		if err != nil {
			// err opcode
			return
		}
		nc := structs.NcUserWillWorldSelectAck{
			Error: 7768,
		}
		data := make([]byte, 0)

		if b, err := networking.WriteBinary(nc.Error); err == nil {
			data = append(data, b...)
		}
		data = append(data, otp...)
		pc.Base.Data = data
		go networking.WriteToClient(ctx, pc)
	}
}

//func userLoginWorldFailAck(ctx context.Context, pc *networking.Command) {
//	select {
//	case <-ctx.Done():
//		return
//	default:
//
//	}
//}
