package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/spf13/viper"
)

func NcMiscSeedAck(ctx context.Context, np *networking.Parameters) {
	xov := ctx.Value(networking.XorOffset)
	xc := xov.(chan uint16)

	xorLimit := uint16(viper.GetInt("crypt.xorLimit"))

	xorOffset := networking.RandomXorKey(xorLimit)
	log.Infof("[xor offset] %v", xorOffset)

	nc := structs.NcMiscSeedAck{
		Seed: xorOffset,
	}

	np.Command.NcStruct = &nc
	np.Command.Send(ctx)

	xc <- xorOffset
}


// NcMiscGameTimeReq requests the server time
// NC_MISC_GAMETIME_REQ
func NcMiscGameTimeReq(ctx context.Context, np * networking.Parameters) {
	go NcMiscGameTimeAck(ctx, &networking.Command{})
}

// NcMiscGameTimeAck sends the current server time
// NC_MISC_GAMETIME_ACK
func NcMiscGameTimeAck(ctx context.Context, pc *networking.Command) {
	pc.Base.OperationCode = 2062
	nc, err := worldTime(ctx)
	if err != nil {
		return
	}
	pc.NcStruct = &nc
	pc.Send(ctx)
}
