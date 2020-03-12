package networking

import (
	"context"
	gs "shine.engine.game_structs"
)

func miscSeedAck(ctx context.Context, pc *Command) {
	select {
	case <-ctx.Done():
		return
	default:
		xov := ctx.Value("xorOffset")
		xo := xov.(*uint16)

		xorOffset := RandomXorKey()
		log.Infof("XorKey: %v", xorOffset)

		*xo = xorOffset

		nc := gs.NcMiscSeedAck{
			Seed: xorOffset,
		}

		if data, err := WriteBinary(nc); err != nil {

		} else {
			pc.Base.Data = data
			go WriteToClient(ctx, pc)
		}
	}
}
