package world

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// NcCharLoginReq handles a petition to login to the zone where the character's location map is running
// NC_CHAR_LOGIN_REQ
func ncCharLoginReq(ctx context.Context, np *networking.Parameters) {
	var (
		nc structs.NcCharLoginReq
		cl characterLoginEvent
	)

	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		log.Error(err)
		return
	}

	cl = characterLoginEvent{
		nc:       &nc,
		np:       np,
		zoneInfo: make(chan *structs.NcCharLoginAck),
		err:      make(chan error),
	}

	worldEvents[characterLogin] <- &cl

	select {
	case nc := <-cl.zoneInfo:
		ncCharLoginAck(np, nc)
		break
	case err := <-cl.err:
		log.Error(err)
		return
	}

}

// NcCharLoginAck sends zone connection info to the client
// NC_CHAR_LOGIN_ACK
func ncCharLoginAck(np *networking.Parameters, nc *structs.NcCharLoginAck) {
	// query the zone master for connection info for the map
	// send it to the client
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4099,
		},
		NcStruct: nc,
	}
	pc.Send(np.OutboundSegments.Send)
}
