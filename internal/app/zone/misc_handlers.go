package zone

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

// for the zone service it is the server that makes use of this handler
func ncMiscHeartBeatReq(p *player) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 2052,
		},
	}
	pc.Send(p.conn.outboundData)
}

// for the zone service it is the client that makes use of this handler
func ncMiscHeartBeatAck(ctx context.Context, np *networking.Parameters) {
	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	zoneEvents[heartbeatUpdate] <- &heartbeatUpdateEvent{
		session: session,
	}
}
