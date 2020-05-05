package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/spf13/viper"
)

func ncMiscSeedAck(ctx context.Context, np *networking.Parameters) {
	select {
	case <-ctx.Done():
		return
	default:

		xov := ctx.Value(networking.XorOffset)
		xo := xov.(*uint16)

		xorLimit := uint16(viper.GetInt("crypt.xorLimit"))

		xorOffset := networking.RandomXorKey(xorLimit)
		log.Infof("[xor offset] %v", xorOffset)

		*xo = xorOffset

		nc := structs.NcMiscSeedAck{
			Seed: xorOffset,
		}
		np.Command.NcStruct = &nc
		np.Command.Send(ctx)
	}
}
