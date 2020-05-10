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

// for the zone service it is the server that makes use of this handler
func ncMiscHeartBeatReq(p *player) {
	pc := networking.Command{
		Base:     networking.CommandBase{
			OperationCode: 2052,
		},
	}
	pc.SendDirectly(p.conn.outboundData)
}

// for the zone service it is the client that makes use of this handler
func ncMiscHeartBeatAck(ctx context.Context, np *networking.Parameters) {
	// this player's session has info about the handle
	// update it
}
