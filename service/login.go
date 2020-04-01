package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
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
		if ncs, ok := lc.pc.NcStruct.(*structs.NcUserClientVersionCheckReq); ok {
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
func checkCredentials(ctx context.Context, req structs.NcUserUsLoginReq) error {
	select {
	case <-ctx.Done():
		return ErrCC
	default:
		un := req.UserName[:]
		pass := req.Password[:]
		userName := strings.TrimRight(string(un), "\x00")
		password := strings.TrimRight(string(pass), "\x00")

		var storedPassword string
		err := db.Model((*User)(nil)).Column("password").Where("user_name = ?", userName).Limit(1).Select(&storedPassword)

		if err != nil {
			return fmt.Errorf("%v: [ %v ]", ErrDBE, err)
		}

		if storedPassword == password {
			// save login session in redis if necessary
			return nil
		}
		return ErrBC
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

func (lc *LoginCommand) serverSelectScreen(ctx context.Context) (structs.NcUserLoginAck, error) {
	grpcc.mu.Lock()
	conn := grpcc.services["world"]
	wc := lw.NewWorldClient(conn)
	grpcc.mu.Unlock()

	rpcCtx, _ := context.WithTimeout(context.Background(), gRPCTimeout)
	wi, err := wc.AvailableWorlds(rpcCtx, &lw.ClientMetadata{
		Ip: "127.0.0.01",
	})

	if err != nil {
		return structs.NcUserLoginAck{}, err
	}

	nc := structs.NcUserLoginAck{}
	for _, w := range wi.Worlds {
		ws := structs.WorldInfo{
			WorldNumber: byte(w.WorldNumber),
			// 1: behaviour -> cannot enter, message -> The server is under maintenance.
			// 2: behaviour -> cannot enter, message -> You cannot connect to an empty server.
			// 3: behaviour -> cannot enter, message -> The server has been reserved for a special use.
			// 4: behaviour -> cannot enter, message -> Login failed due to an unknown error.
			// 5: behaviour -> cannot enter, message -> The server is full.
			// 6: behaviour -> ok
			WorldStatus: byte(w.WorldStatus),
		}
		copy(ws.WorldName.Name[:], w.WorldName)
		nc.Worlds = append(nc.Worlds, ws)
	}
	nc.NumOfWorld = byte(len(nc.Worlds))
	return nc, nil
}

// request info about selected world
func userSelectedServer(ctx context.Context) (*lw.WorldConnectionInfo, error) {
	select {
	case <-ctx.Done():
		return &lw.WorldConnectionInfo{}, ErrCC
	default:
		grpcc.mu.Lock()
		conn := grpcc.services["world"]
		c := lw.NewWorldClient(conn)
		grpcc.mu.Unlock()

		rpcCtx, _ := context.WithTimeout(context.Background(), gRPCTimeout)

		wci, err := c.ConnectionInfo(rpcCtx, &lw.SelectedWorld{Num: 1})
		if err != nil {
			return &lw.WorldConnectionInfo{}, err
		}
		return wci, nil
	}
}

// verify the token matches the one stored [on redis] by the world service
func (lc *LoginCommand) loginByCode(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ErrCC
	default:
		if ncs, ok := lc.pc.NcStruct.(*structs.NcUserLoginWithOtpReq); ok {
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
