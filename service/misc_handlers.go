package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
)

func miscGametimeReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		go miscGameTimeAck(ctx, &networking.Command{})
	}
}

func miscGameTimeAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base.OperationCode = 2062
		wc := &WorldCommand{pc: pc}
		nc, err := wc.worldTime(ctx)
		if err != nil {
			return
		}
		data, err := structs.Pack(nc)
		if err != nil {
			return
		}
		pc.Base.Data = data
		go networking.WriteToClient(ctx, pc)
	}
}