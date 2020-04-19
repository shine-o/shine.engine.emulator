package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NcPrisonGetReq client asks for how much time in prison ô_o
// NC_PRISON_GET_REQ
func NcPrisonGetReq(ctx context.Context, pc * networking.Command) {
	go NcPrisonGetAck(ctx)
}

// NcPrisonGetAck sends how much time the player spends in prison
// NC_PRISON_GET_ACK
func NcPrisonGetAck(ctx context.Context) {
	pc := networking.Command{
		Base:     networking.CommandBase{
			OperationCode: 31751,
		},
		NcStruct: &structs.NcPrisonGetAck{
			Err:    3505,
			Minute: 0,
		},
	}
	go pc.Send(ctx)
}