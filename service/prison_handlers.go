package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NcPrisonGetReq client asks for how much time in prison ô_o
// NC_PRISON_GET_REQ
func ncPrisonGetReq(ctx context.Context, np * networking.Parameters) {
	go ncPrisonGetAck(np)
}

// NcPrisonGetAck sends how much time the player spends in prison
// NC_PRISON_GET_ACK
func ncPrisonGetAck(np * networking.Parameters) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 31751,
		},
		NcStruct: &structs.NcPrisonGetAck{
			Err:    3505,
			Minute: 0,
		},
	}
	pc.Send(np.OutboundSegments.Send)
}
