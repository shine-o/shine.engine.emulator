package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.networking/structs"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

// ErrWTO world service timed out
var ErrWTO = errors.New("world timed out")

// ErrCC ctx.Done() signal was received
var ErrCC = errors.New("context was canceled")

// ErrNoWorld no world is available
var ErrNoWorld = errors.New("no world is available")

// ErrBC user sent bad userName and password combination
var ErrBC = errors.New("bad credentials")

// ErrDBE database exception
var ErrDBE = errors.New("database exception")

// check that the client version is correct
func checkClientVersion(ctx context.Context, req structs.NcUserClientVersionCheckReq) ([]byte, error) {
	var data []byte
	select {
	case <-ctx.Done():
		return data, ErrCC
	default:
		if viper.GetBool("crypt.client_version.ignore") {
			return data, nil
		}
		vk := strings.TrimRight(string(req.VersionKey[:33]), "\x00") // will replace with direct binary comparison
		if vk == viper.GetString("crypt.client_version.key") {
			// xtrap info goes here, but we dont use xtrap so we don't have to send anything.
			return data, nil
		}
		return data, fmt.Errorf("client sent incorrect client version key:%v", vk)
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
func checkWorldStatus(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ErrCC
	default:
		if !viper.IsSet("worlds") {
			return ErrNoWorld
		}

		worlds := viper.GetStringSlice("worlds")

		for _, w := range worlds {
			clientKey := fmt.Sprintf("gRPC.clients.%v", w)
			_, err := newRpcConn(clientKey)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

type World struct {
	ID uint8
	Name string
}

type AvailableWorlds []World

func avalWorlds() (AvailableWorlds, error)  {
	aw := AvailableWorlds{}
	if !viper.IsSet("worlds") {
		return aw, ErrNoWorld
	}
	
	worlds := make([]map[string]string, 0)
	var m map[string]string
	worldsI := viper.Get("worlds")
	worldsS := worldsI.([]interface{})
	for _, s := range worldsS {
		serviceMap := s.(map[interface{}]interface{})
		m = make(map[string]string)
		for k, v := range serviceMap {
			m[k.(string)] = v.(string)
		}
		worlds = append(worlds, m)
	}
	
	for _, v := range worlds {
		id, err := strconv.Atoi(v["id"])
		if err != nil {
			log.Error(err)
			continue
		}
		w := World{
			ID:   uint8(id),
			Name: v["name"],
		}
		aw = append(aw, w)
	}
	return aw, nil
}

func serverSelectScreen(ctx context.Context) (structs.NcUserLoginAck, error) {
	select {
	case <-ctx.Done():
		return structs.NcUserLoginAck{}, nil
	default:
		nc := structs.NcUserLoginAck{}

		aw, err := avalWorlds()
		if err != nil {
			return structs.NcUserLoginAck{}, err
		}

		for _, w := range aw {
			conn, err := newRpcConn(w.Name)
			if err != nil {
				return nc, err
			}
			c := lw.NewWorldClient(conn)

			wd, err := c.GetWorldData(ctx, &lw.WorldQuery{
				Name: w.Name,
			})

			if err != nil {
				log.Error(err)
				continue
			}
			wi := structs.WorldInfo{
				WorldNumber: byte(wd.WorldNumber),
				// 1: behaviour -> cannot enter, message -> The server is under maintenance.
				// 2: behaviour -> cannot enter, message -> You cannot connect to an empty server.
				// 3: behaviour -> cannot enter, message -> The server has been reserved for a special use.
				// 4: behaviour -> cannot enter, message -> Login failed due to an unknown error.
				// 5: behaviour -> cannot enter, message -> The server is full.
				// 6: behaviour -> ok
				WorldStatus: byte(wd.WorldStatus),
			}
			nc.Worlds = append(nc.Worlds, wi)
		}
		return nc, nil
	}
}

// request info about selected world
func userSelectedServer(ctx context.Context, req structs.NcUserWorldSelectReq) (*lw.WorldData, error) {
	select {
	case <-ctx.Done():
		return &lw.WorldData{}, ErrCC
	default:

		aw, err := avalWorlds()

		if err != nil {
			return &lw.WorldData{}, err

		}
		for _, w := range aw {
			if w.ID == req.WorldNo {
				conn, err := newRpcConn(w.Name)

				if err != nil {
					return &lw.WorldData{}, err
				}

				c := lw.NewWorldClient(conn)

				wd, err := c.GetWorldData(ctx, &lw.WorldQuery{
					Id: int32(req.WorldNo),
				})

				if err != nil {
					return &lw.WorldData{}, err
				}

				return wd, nil
			}
		}
		return &lw.WorldData{}, ErrNoWorld
	}
}

// verify the token matches the one stored [on redis] by the world service
func loginByCode(ctx context.Context, req structs.NcUserLoginWithOtpReq) error {
	select {
	case <-ctx.Done():
		return ErrCC
	default:
		b := make([]byte, len(req.Otp.Name))
		copy(b, req.Otp.Name[:])
		if _, err := redisClient.Get(string(b)).Result(); err != nil {
			return err
		}
		return nil
	}
}
