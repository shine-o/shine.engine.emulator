package net

import (
	"context"
	"github.com/shine-o/shine.engine.core/structs"
)

func miscSeedAck(ctx context.Context, pc *Command) {
	select {
	case <-ctx.Done():
		return
	default:
		xov := ctx.Value(XorOffset)
		xo := xov.(*uint16)

		xorOffset := RandomXorKey()
		log.Infof("[xor offset] %v", xorOffset)

		*xo = xorOffset

		nc := structs.NcMiscSeedAck{
			Seed: xorOffset,
		}

		if data, err := structs.Pack(&nc); err != nil {

		} else {
			pc.Base.Data = data
			go WriteToClient(ctx, pc)
		}
	}
}
