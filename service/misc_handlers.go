package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
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
		wc := &WorldCommand{pc:pc}
		if data, err := wc.worldTime(ctx); err != nil {
			return
		} else {
			pc.Base.Data = data
			go networking.WriteToClient(ctx, pc)
		}
	}
}
