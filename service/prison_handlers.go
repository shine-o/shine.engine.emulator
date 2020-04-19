package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

//NC_PRISON_GET_REQ
func prisonGetReq(ctx context.Context, pc * networking.Command) {
	go prisonGetAck(ctx)
}


//NC_PRISON_GET_ACK
func prisonGetAck(ctx context.Context) {
	pc := networking.Command{
		Base:     networking.CommandBase{
			OperationCode:31751,
		},
		NcStruct: &structs.NcPrisonGetAck{
			Err:    3505,
			Minute: 0,
		},
	}
	go pc.Send(ctx)
}