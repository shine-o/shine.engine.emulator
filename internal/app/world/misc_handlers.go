package world

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/crypto"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"github.com/spf13/viper"
)

func ncMiscSeedAck(ctx context.Context, np *networking.Parameters) {
	xov := ctx.Value(networking.XorOffset)
	xc := xov.(chan uint16)

	xorLimit := uint16(viper.GetInt("crypt.xorLimit"))

	xorOffset := crypto.RandomXorKey(xorLimit)
	log.Infof("[xor offset] %v", xorOffset)

	nc := structs.NcMiscSeedAck{
		Seed: xorOffset,
	}

	networking.Send(np.OutboundSegments.Send, networking.NC_MISC_SEED_ACK, &nc)

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
