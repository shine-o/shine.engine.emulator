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
	np.Command.Send(np.OutboundSegments.Send)

	xc <- xorOffset
}

// NcMiscGameTimeReq requests the server time
// NC_MISC_GAMETIME_REQ
func NcMiscGameTimeReq(ctx context.Context, np * networking.Parameters) {
	var ste serverTimeEvent

	ste = serverTimeEvent{
		np:  np,
		err: make(chan error),
	}

	worldEvents[serverTime] <- &ste
	select {
	case e := <- ste.err:
		log.Error(e)
	}
}

// NcMiscGameTimeAck sends the current server time
// NC_MISC_GAMETIME_ACK
func NcMiscGameTimeAck(np * networking.Parameters, nc * structs.NcMiscGameTimeAck) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 2062,
		},
		NcStruct : nc,
	}
	pc.Send(np.OutboundSegments.Send)
}
