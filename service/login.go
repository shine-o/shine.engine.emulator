package service

import (
	"context"
	"fmt"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.structs"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

type LoginCommand struct {
	pc * networking.Command
}

// check that the client version is correct
func (lc * LoginCommand) checkClientVersion(ctx context.Context) ([]byte, error) {
	select {
	case <- ctx.Done():
		return nil, fmt.Errorf("cancel signal received")
	default:
		var data []byte
		if s, ok := lc.pc.NcStruct.(structs.NcUserClientVersionCheckReq); ok {
			vk := strings.TrimRight(string(s.VersionKey[1:33]), "\x00") // will replace with direct binary comparison
			if vk == viper.GetString("crypt.client_version") {
				// xtrap info goes here, but we dont use xtrap so we don't have to send anything.
				return data, nil
			} else {
				return data, fmt.Errorf("client sent incorrect client version key:%v", vk)
			}
		} else {
			return data, fmt.Errorf("bad NcStruct: %v", reflect.TypeOf(lc.pc.NcStruct).String())
		}
	}
}

// check against database that the user name and password combination are correct
func  (lc * LoginCommand) checkCredentials(ctx context.Context) ([]byte, error) {
	return nil, nil
}

// check the world service is up and running
func (lc * LoginCommand) checkWorldStatus(ctx context.Context) ([]byte, error) {
	return nil, nil
}

// request info about selected world
func (lc * LoginCommand) userSelectedServer(ctx context.Context) ([]byte, error) {
	return nil, nil
}

// verify the token matches the one stored [on redis] by the world service
func (lc * LoginCommand) loginByCode(ctx context.Context) ([]byte, error) {
	return nil, nil
}


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
			if user.Password == password {
				go userLoginAck(ctx, &networking.Command{})
			} else {
				go userLoginFailAck(ctx, &networking.Command{})
			}
		}
	}
}