package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
	"shine.engine.structs"
	"time"
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

		var (
			t                    time.Time
			hour, minute, second byte
		)

		t = time.Now()
		hour = byte(t.Hour())
		minute = byte(t.Minute())
		second = byte(t.Second())

		nc := &structs.NcMiscGameTimeAck{
			Hour:   hour,
			Minute: minute,
			Second: second,
		}
		if data, err := networking.WriteBinary(nc); err != nil {

		} else {
			pc.Base.Data = data
			go networking.WriteToClient(ctx, pc)
		}
	}
}
