package world

import (
	"context"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NcPrisonGetReq client asks for how much time in prison ô_o
// NC_PRISON_GET_REQ
func ncPrisonGetReq(ctx context.Context, np *networking.Parameters) {
	nc := &structs.NcPrisonGetAck{
		Err:    3505,
		Minute: 0,
	}
	networking.Send(np.OutboundSegments.Send, networking.NC_PRISON_GET_ACK, &nc)
}
