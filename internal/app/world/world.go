package world

import (
	"context"
	"errors"
	"github.com/go-pg/pg/v10"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/login"
	wm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/world-master"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"github.com/spf13/viper"
	"time"
)

type world struct {
	db *pg.DB
	events
	dynamic
}

var worldEvents sendEvents

func (w *world) load() {
	// register to world-master
	err := registerWorld()

	if err != nil {
		log.Fatal(err)
	}

	w.events = events{
		send: make(sendEvents),
		recv: make(recvEvents),
	}

	events := []eventIndex{
		serverSelect, serverSelectToken, serverTime,
		createCharacter, deleteCharacter,
		characterLogin, characterSettings,
		updateShortcuts, updateGameSettings, updateKeymap, characterSelect,
	}

	for _, index := range events {
		c := make(chan event, 5)
		w.send[index] = c
		w.recv[index] = c
	}

	worldEvents = w.send
}

func (w *world) run() {
	go w.session()
	go w.characterCRUD()
	go w.characterSession()
}

func registerWorld() error {
	id := viper.GetInt32("world.id")
	name := viper.GetString("world.name")
	ip := viper.GetString("world.external_ip")
	port := viper.GetInt32("world.port")

	conn, err := newRPCClient("world_master")

	if err != nil {
		return err
	}
	c := wm.NewMasterClient(conn)
	rpcCtx, cancel := context.WithTimeout(context.Background(), gRPCTimeout)
	defer cancel()

	wr, err := c.RegisterWorld(rpcCtx, &wm.WorldDetails{
		ID:   id,
		Name: name,
		Conn: &wm.ConnectionInfo{
			IP:   ip,
			Port: port,
		},
	})

	if err != nil {
		return err
	}

	if !wr.Success {
		return errors.New("failed to register against the world master")
	}

	return nil
}

var errCC = errors.New("context was canceled")

var errAccountInfo = errors.New("failed to fetch account info")

// user wants to log to given service
// check if service is okay
// take user name, persist to redis
func verifyUser(ws *session, req *structs.NcUserLoginWorldReq) error {
	conn, err := newRPCClient("login")
	if err != nil {
		return err
	}

	defer conn.Close()

	c := login.NewLoginClient(conn)

	rpcCtx, cancel := context.WithTimeout(context.Background(), gRPCTimeout)
	defer cancel()

	ui, err := c.AccountInfo(rpcCtx, &login.User{
		UserName: req.User.Name,
	})

	if err != nil {
		return errAccountInfo
	}

	ws.UserID = ui.UserID
	ws.UserName = req.User.Name

	if err := persistSession(ws); err != nil {
		ws.UserID = 0
		ws.UserName = ""
		return err
	}
	return nil
}

func userCharacters(db *pg.DB, ws *session) (structs.NcUserLoginWorldAck, error) {
	worldID := viper.GetInt("service.id")

	var avatars []structs.AvatarInformation
	var chars []game.Character

	err := db.Model(&chars).
		Relation("Appearance").
		Relation("Attributes").
		Relation("Location").
		Relation("EquippedItems").
		Where("user_id = ?", ws.UserID).
		Select()

	if err != nil {
		return structs.NcUserLoginWorldAck{}, err
	}

	if len(chars) > 0 {
		for _, c := range chars {
			avatars = append(avatars, c.NcRepresentation())
		}
	}

	nc := structs.NcUserLoginWorldAck{
		WorldManager: uint16(worldID),
		NumOfAvatar:  byte(len(chars)),
		Avatars:      avatars,
	}

	return nc, nil
}

// user clicked previous
// generate a otp token and store it in redis
// login service will use the token to authenticate the user and send him to server select
func returnToServerSelect() (structs.NcUserWillWorldSelectAck, error) {
	otp := randStringBytesMaskImprSrcUnsafe(32)
	err := redisClient.Set(otp, otp, 20*time.Second).Err()
	if err != nil {
		return structs.NcUserWillWorldSelectAck{}, err
	}
	nc := structs.NcUserWillWorldSelectAck{
		Error: 7768,
		Otp: structs.Name8{
			Name: otp,
		},
	}
	return nc, nil
}
