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
		nc: &nc,
		np: np,
	}

	worldEvents[characterLogin] <- &cl
}