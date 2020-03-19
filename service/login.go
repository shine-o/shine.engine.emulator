package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.networking"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
	"github.com/shine-o/shine.engine.structs"
	"github.com/spf13/viper"
	"reflect"
	"strings"
)

// ErrWTO world service timed out
var ErrWTO = errors.New("world timed out")

// ErrUNF user does not exist in the database
var ErrUNF = errors.New("user not found")

// ErrCC ctx.Done() signal was received
var ErrCC = errors.New("context was canceled")

// ErrBC user sent bad userName and password combination
var ErrBC = errors.New("bad credentials")

// ErrDBE database exception
var ErrDBE = errors.New("database exception")

// LoginCommand wrapper for networking command
// any information scoped to this service and its handlers can be added here
type LoginCommand struct {
	pc *networking.Command
}

// check that the client version is correct
func (lc *LoginCommand) checkClientVersion(ctx context.Context) ([]byte, error) {
	var data []byte
	select {
	case <-ctx.Done():
		return data, ErrCC
	default:
		if viper.GetBool("crypt.client_version.ignore") {
			return data, nil
		}
		if ncs, ok := lc.pc.NcStruct.(structs.NcUserClientVersionCheckReq); ok {
			vk := strings.TrimRight(string(ncs.VersionKey[:33]), "\x00") // will replace with direct binary comparison
			if vk == viper.GetString("crypt.client_version.key") {
				// xtrap info goes here, but we dont use xtrap so we don't have to send anything.
				return data, nil
			}
			return data, fmt.Errorf("client sent incorrect client version key:%v", vk)
		}
		return data, fmt.Errorf("unexpected struct type: %v", reflect.TypeOf(lc.pc.NcStruct).String())
	}
}

// check against database that the user name and password combination are correct
func (lc *LoginCommand) checkCredentials(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ErrCC
	default:
		if ncs, ok := lc.pc.NcStruct.(structs.NcUserUsLoginReq); ok {
			un := ncs.UserName[:]
			pass := ncs.Password[:]
			userName := strings.TrimRight(string(un), "\x00")
			password := strings.TrimRight(string(pass), "\x00")

			var user User
			db := database.Where("user_name = ?", userName).First(&user)

			if len(db.GetErrors()) > 0 {

				if db.RecordNotFound() {
					return ErrUNF
				}

				return fmt.Errorf("%v: [ %v ]", ErrDBE, db.GetErrors())
			}

			if user.Password == password {
				return nil
			}

			return ErrBC

		}

		return fmt.Errorf("unexpected struct type: %v", reflect.TypeOf(lc.pc.NcStruct).String())
	}
}

// check the world service is up and running
func (lc *LoginCommand) checkWorldStatus(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ErrCC
	default:
		grpcc.mu.Lock()
		conn := grpcc.services["world"]
		state := conn.GetState().String()
		grpcc.mu.Unlock()

		if state == "READY" || state == "IDLE" {
			return nil
		}
		return ErrWTO
	}
}

// request info about selected world
func (lc *LoginCommand) userSelectedServer(ctx context.Context) ([]byte, error) {
	var data []byte
	select {
	case <-ctx.Done():
		return data, ErrCC
	default:
		grpcc.mu.Lock()
		conn := grpcc.services["world"]
		c := lw.NewWorldClient(conn)
		grpcc.mu.Unlock()

		rpcCtx, _ := context.WithTimeout(context.Background(), gRPCTimeout)

		r, err := c.ConnectionInfo(rpcCtx, &lw.SelectedWorld{Num: 1})
		if err != nil {
			return data, err
		}
		return r.Info, nil
	}
}

// verify the token matches the one stored [on redis] by the world service
func (lc *LoginCommand) loginByCode(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ErrCC
	default:
		if ncs, ok := lc.pc.NcStruct.(structs.NcUserLoginWithOtpReq); ok {
			b := make([]byte, len(ncs.Otp.Name))
			copy(b, ncs.Otp.Name[:])
			if _, err := redisClient.Get(string(b)).Result(); err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("unexpected struct type: %v", reflect.TypeOf(lc.pc.NcStruct).String())
	}
}
