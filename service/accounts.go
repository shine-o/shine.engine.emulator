package service

import (
	"context"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.structs"
	"strings"
)

func authenticate(ctx context.Context, nc *structs.NcUserUsLoginReq) {
	select {
	case <-ctx.Done():
		go userLoginFailAck(ctx, &networking.Command{})
		return
	default:
		un := nc.UserName[:]
		pass := nc.Password[:]
		userName := strings.TrimRight(string(un), "\x00")
		password := strings.TrimRight(string(pass), "\x00")

		var user User
		if database.Where("user_name = ?", userName).First(&user).RecordNotFound() {
			go userLoginFailAck(ctx, &networking.Command{})
		} else {
			log.Info(user)
			if user.Password == password {
				go userLoginAck(ctx, &networking.Command{})
			} else {
				go userLoginFailAck(ctx, &networking.Command{})
			}
		}
	}
}