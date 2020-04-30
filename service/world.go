package service

import (
	"context"
	"errors"
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/grpc/login"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/spf13/viper"
	"strings"
	"time"
)

var errCC = errors.New("context was canceled")

var errAccountInfo = errors.New("failed to fetch account info")

func worldTime(ctx context.Context) (structs.NcMiscGameTimeAck, error) {
	var (
		t                    time.Time
		hour, minute, second byte
	)

	t = time.Now()
	hour = byte(t.Hour())
	minute = byte(t.Minute())
	second = byte(t.Second())

	return structs.NcMiscGameTimeAck{
		Hour:   hour,
		Minute: minute,
		Second: second,
	}, nil
}

// user wants to log to given service
// check if service is okay
// take user name, persist to redis
func loginToWorld(ctx context.Context, req structs.NcUserLoginWorldReq) error {
	wsi := ctx.Value(networking.ShineSession)
	ws := wsi.(*session)
	userName := strings.TrimRight(string(req.User.Name[:]), "\x00")
	ws.UserName = userName

	conn, err := newRPCClient("login")
	if err != nil {
		return err
	}
	defer conn.Close()

	c := login.NewLoginClient(conn)

	rpcCtx, _ := context.WithTimeout(context.Background(), gRPCTimeout)

	ui, err := c.AccountInfo(rpcCtx, &login.User{
		UserName: userName,
	})

	if err != nil {
		return errAccountInfo
	}

	ws.UserID = ui.UserID

	if err := persistSession(ws); err != nil {
		return err
	}
	return nil
}

func userWorldInfo(ctx context.Context) (structs.NcUserLoginWorldAck, error) {
	wsi := ctx.Value(networking.ShineSession)
	ws := wsi.(*session)
	worldID := viper.GetInt("service.id")

	var avatars []structs.AvatarInformation
	var chars []character.Character

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
	err := redisClient.Set(otp, otp, 20 * time.Second).Err()
	if err != nil {
		return structs.NcUserWillWorldSelectAck{}, err
	}
	nc := structs.NcUserWillWorldSelectAck{
		Error: 7768,
		Otp:   structs.Name8{
			Name: otp,
		},
	}
	return nc, nil
}
