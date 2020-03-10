package protocol

import (
	"context"
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

		nc := ncMiscSeedAck{
			seed: xorOffset,
		}

		if data, err := WriteBinary(nc); err != nil {

		} else {
			pc.Base.Data = data
			go WriteToClient(ctx, pc)
		}
	}
}
