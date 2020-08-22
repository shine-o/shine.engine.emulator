package login

import (
	"context"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"

	wm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/world-master"
	"github.com/spf13/viper"
	"strings"
)

type login struct {
	worlds map[int]world
	events
	dynamic
}

type world struct {
	id int
	name string
	ip string
	port int
}

var loginEvents sendEvents

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


func (l *login) load() {
	indexes := []eventIndex{
		clientVersion,
		credentialsLogin,
		worldManagerStatus,
		serverList,
		serverSelect,
		tokenLogin,
	}

	l.events = events{
		send: make(sendEvents),
		recv: make(recvEvents),
	}

	for _, index := range indexes {
		c := make(chan event, 5)
		l.send[index] = c
		l.recv[index] = c
	}

	loginEvents = l.send

	err := l.availableWorlds()
	if err != nil{
		log.Fatal(err)
	}
	go l.startWorkers()
}

func (l*login) startWorkers()  {
	go l.authentication()
}

// 1: behaviour -> cannot enter, message -> The server is under maintenance.
// 2: behaviour -> cannot enter, message -> You cannot connect to an empty server.
// 3: behaviour -> cannot enter, message -> The server has been reserved for a special use.
// 4: behaviour -> cannot enter, message -> Login failed due to an unknown error.
// 5: behaviour -> cannot enter, message -> The server is full.
// 6: behaviour -> ok
func (l *login) availableWorlds() error {
	l.worlds = make(map[int]world)

	conn, err := newRPCClient("world_master")

	if err != nil {
		return err
	}

	defer conn.Close()

	c := wm.NewMasterClient(conn)

	ctx := context.Background()

	worlds, err := c.GetWorlds(ctx, &wm.Empty{})

	if err != nil {
		return err
	}

	for _, wd := range worlds.List {
		l.worlds[int(wd.ID)] = world{
			id:   int(wd.ID),
			name: wd.Name,
			ip:   wd.Conn.IP,
			port: int(wd.Conn.Port),
		}
 	}
	return  nil
}

// check that the client version is correct
func checkClientVersion(req *structs.NcUserClientVersionCheckReq) error {
	if viper.GetBool("crypt.client_version.ignore") {
		return nil
	}
	vk := strings.TrimRight(string(req.VersionKey[:33]), "\x00") // will replace with direct binary comparison
	if vk == viper.GetString("crypt.client_version.key") {
		// xtrap info goes here, but we dont use xtrap so we don't have to send anything.
		return nil
	}
	return fmt.Errorf("client sent incorrect client version key:%v", vk)
}

// check against database that the user name and password combination are correct
func checkCredentials(req *structs.NcUserUsLoginReq) error {
	var storedPassword string
	err := db.Model((*User)(nil)).Column("password").Where("user_name = ?", req.UserName).Limit(1).Select(&storedPassword)

	if err != nil {
		return fmt.Errorf("%v: [ %v ]", ErrDBE, err)
	}

	if storedPassword == req.Password {
		return nil
	}

	return ErrBC
}

// verify the token matches the one stored [on redis] by the world service
func loginByCode(req * structs.NcUserLoginWithOtpReq) error {
	b := make([]byte, len(req.Otp.Name))
	copy(b, req.Otp.Name[:])
	if _, err := redisClient.Get(string(b)).Result(); err != nil {
		return err
	}
	return nil
}

func loginSuccessful(l *login, np *networking.Parameters) {
	nc := structs.NcUserLoginAck{}
	for _, w := range l.worlds {
		nc.Worlds = append(nc.Worlds, structs.WorldInfo{
			WorldNumber: byte(w.id),
			WorldName: structs.Name4{
				Name: w.name,
			},
			WorldStatus: 6,
		})
	}
	nc.NumOfWorld = byte(len(l.worlds))
	ncUserLoginAck(np, nc)
}