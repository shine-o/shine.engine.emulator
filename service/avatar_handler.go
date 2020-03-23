package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
	"reflect"
)

func avatarCreateReq(ctx context.Context, pc *networking.Command)   {
	select {
	case <- ctx.Done():
		return
	default:
		// basically create the character
		if nc, ok := pc.NcStruct.(structs.NcAvatarCreateReq); ok {
			log.Error(nc)
		} else {
			log.Error(reflect.TypeOf(nc).String())
			log.Error(pc.NcStruct)
		}
	}
}