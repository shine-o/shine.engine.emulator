package login

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"

	"github.com/spf13/viper"
)

func ncMiscSeedAck(ctx context.Context, np *networking.Parameters) {
	xov := ctx.Value(networking.XorOffset)
	xc := xov.(chan uint16)

	xorLimit := uint16(viper.GetInt("crypt.xorLimit"))

	xorOffset := networking.RandomXorKey(xorLimit)
	log.Infof("[xor offset] %v", xorOffset)

	nc := structs.NcMiscSeedAck{
		Seed: xorOffset,
	}

	networking.Send(np.OutboundSegments.Send, networking.NC_MISC_SEED_ACK, nc)

	xc <- xorOffset
}
