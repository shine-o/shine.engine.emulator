package service

import (
	"context"
	"errors"
	"fmt"
	w "github.com/shine-o/shine.engine.core/grpc/world"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/spf13/viper"
	"strconv"
	"strings"
)

type world struct {
	ID   uint8
	Name string
}

type availableWorlds []world

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
func checkClientVersion(req structs.NcUserClientVersionCheckReq) ([]byte, error) {
	var data []byte
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

// check against database that the user name and password combination are correct
func checkCredentials(req structs.NcUserUsLoginReq) error {
	var storedPassword string
	err := db.Model((*User)(nil)).Column("password").Where("user_name = ?", req.UserName).Limit(1).Select(&storedPassword)

	if err != nil {
		return fmt.Errorf("%v: [ %v ]", ErrDBE, err)
	}

	if storedPassword == req.Password {
		// save login session in redis if necessary
		return nil
	}
	return ErrBC
}

// check the world service is up and running
func checkWorldStatus() error {
	if !viper.IsSet("worlds") {
		return ErrNoWorld
	}

	worlds, err := worlds()
	if err != nil {
		return err
	}
	for _, aw := range worlds {
		c, err := newRPCClient(aw.Name)
		if err != nil {
			return err
		}
		c.Close()
	}
	return nil
}


func worlds() (availableWorlds, error) {
	aw := availableWorlds{}
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
		nw := world{
			ID:   uint8(id),
			Name: v["name"],
		}
		aw = append(aw, nw)
	}
	return aw, nil
}

func serverSelectScreen(ctx context.Context) (structs.NcUserLoginAck, error) {
	nc := structs.NcUserLoginAck{}
	aw, err := worlds()
	if err != nil {
		return structs.NcUserLoginAck{}, err
	}

	for _, v := range aw {
		conn, err := newRPCClient(v.Name)
		if err != nil {
			return nc, err
		}
		c := w.NewWorldClient(conn)

		wd, err := c.GetWorldData(ctx, &w.WorldQuery{
			Name: v.Name,
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

// request info about selected world
func userSelectedServer(ctx context.Context, req structs.NcUserWorldSelectReq) (*w.WorldData, error) {
	aw, err := worlds()

	if err != nil {
		return &w.WorldData{}, err
	}

	for _, v := range aw {
		if v.ID == req.WorldNo {
			conn, err := newRPCClient(v.Name)

			if err != nil {
				return &w.WorldData{}, err
			}

			c := w.NewWorldClient(conn)

			wd, err := c.GetWorldData(ctx, &w.WorldQuery{
				ID: int32(req.WorldNo),
			})

			if err != nil {
				return &w.WorldData{}, err
			}

			return wd, nil
		}
	}
	return &w.WorldData{}, ErrNoWorld
}

// verify the token matches the one stored [on redis] by the world service
func loginByCode(req structs.NcUserLoginWithOtpReq) error {
	b := make([]byte, len(req.Otp.Name))
	copy(b, req.Otp.Name[:])
	if _, err := redisClient.Get(string(b)).Result(); err != nil {
		return err
	}
	return nil
}
