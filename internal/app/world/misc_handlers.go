package world

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

	np.Command.NcStruct = &nc
	np.Command.Send(np.OutboundSegments.Send)

	xc <- xorOffset
}

// NcMiscGameTimeReq requests the server time
// NC_MISC_GAMETIME_REQ
func ncMiscGameTimeReq(ctx context.Context, np *networking.Parameters) {
	ste := serverTimeEvent{
		np: np,
	}

	worldEvents[serverTime] <- &ste
}

// NcMiscGameTimeAck sends the current server time
// NC_MISC_GAMETIME_ACK
func NcMiscGameTimeAck(np *networking.Parameters, nc *structs.NcMiscGameTimeAck) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 2062,
		},
		NcStruct: nc,
	}
	pc.Send(np.OutboundSegments.Send)
}
