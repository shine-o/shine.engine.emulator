package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
)

// NcMiscGameTimeReq requests the server time
// NC_MISC_GAMETIME_REQ
func NcMiscGameTimeReq(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		go NcMiscGameTimeAck(ctx, &networking.Command{})
	}
}

// NcMiscGameTimeAck sends the current server time
// NC_MISC_GAMETIME_ACK
func NcMiscGameTimeAck(ctx context.Context, pc *networking.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		pc.Base.OperationCode = 2062
		nc, err := worldTime(ctx)
		if err != nil {
			return
		}
		pc.NcStruct = &nc
		go pc.Send(ctx)
	}
}
