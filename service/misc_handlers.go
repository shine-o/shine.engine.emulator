package service

import (
	"context"
	protocol "github.com/shine-o/shine.engine.protocol"
)

func miscSeedAck(ctx context.Context, pc *protocol.Command) {
	select {
	case <-ctx.Done():
		return
	default:
		xov := ctx.Value("xorOffset")
		xo := xov.(*uint16)

		xorOffset := protocol.RandomXorKey()
		log.Infof("XorKey: %v", xorOffset)

		*xo = xorOffset

		nc := ncMiscSeedAck{
			seed: xorOffset,
		}

		if data, err := writeBinary(nc); err != nil {

		} else {
			pc.Base.SetData(data)
			go writeToClient(ctx, pc)
		}
	}
}
