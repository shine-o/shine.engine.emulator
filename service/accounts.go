package service

import (
	"context"
	networking "github.com/shine-o/shine.engine.networking"
	gs "shine.engine.game_structs"
	"strings"
)

func authenticate(ctx context.Context, nc *gs.NcUserUsLoginReq) {
	select {
	case <-ctx.Done():
		go userLoginFailAck(ctx, &networking.Command{})
		return
	default:
		un := nc.UserName[:]
		pass := nc.Password[:]
		userName := strings.TrimRight(string(un), "\x00")
		password := strings.TrimRight(string(pass), "\x00")
		if userName == "admin" && password == "21232f297a57a5a743894a0e4a801fc3" { // temporary :)
			go userLoginAck(ctx, &networking.Command{})
		} else {
			go userLoginFailAck(ctx, &networking.Command{})
		}
	}
}
