package manager

import (
	"context"
	"errors"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.networking/structs"
	lw "github.com/shine-o/shine.engine.protocol-buffers/login-world"
	"strconv"
	"strings"
	"time"
)

var errCC = errors.New("context was canceled")

var errAccountInfo = errors.New("failed to fetch account info")

func worldTime(ctx context.Context) (structs.NcMiscGameTimeAck, error) {
	select {
	case <-ctx.Done():
		return structs.NcMiscGameTimeAck{}, errCC
	default:
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
}

// user wants to log to given world
// check if world is okay
// take user name, persist to redis
func loginToWorld(ctx context.Context, req structs.NcUserLoginWorldReq) error {
	select {
	case <-ctx.Done():
		return errCC
	default:
		wsi := ctx.Value(networking.ShineSession)
		ws := wsi.(*session)
		userName := strings.TrimRight(string(req.User.Name[:]), "\x00")
		ws.UserName = userName

		// fetch user id using user name, login serve will check if it has a session for that user
		grpcc.mu.Lock()
		conn := grpcc.services["login"]
		c := lw.NewLoginClient(conn)
		grpcc.mu.Unlock()

		rpcCtx, _ := context.WithTimeout(context.Background(), gRPCTimeout)

		ui, err := c.AccountInfo(rpcCtx, &lw.User{
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
}

func userWorldInfo(ctx context.Context) (structs.NcUserLoginWorldAck, error) {
	select {
	case <-ctx.Done():
		return structs.NcUserLoginWorldAck{}, errCC
	default:
		wsi := ctx.Value(networking.ShineSession)
		ws := wsi.(*session)
		worldID, err := strconv.Atoi(ws.WorldID)

		if err != nil {
			return structs.NcUserLoginWorldAck{}, err
		}

		var avatars []structs.AvatarInformation
		var chars []Character

		err = worldDB.Model(&chars).
			Relation("Appearance").
			//Where("user_id = ?", ws.UserID).
			Relation("Attributes").
			Relation("Location").
			Relation("Inventory").
			Relation("EquippedItems").
			Where("user_id = ?", ws.UserID).
			Select()

		if err != nil {
			return structs.NcUserLoginWorldAck{}, err
		}

		if len(chars) > 0 {
			for _, c := range chars {
				avatars = append(avatars, c.ncRepresentation())
			}
		}

		nc := structs.NcUserLoginWorldAck{
			WorldManager: uint16(worldID),
			NumOfAvatar:  byte(len(chars)),
			Avatars:      avatars,
		}

		return nc, nil
	}
}

// user clicked previous
// generate a otp token and store it in redis
// login manager will use the token to authenticate the user and send him to server select
func returnToServerSelect(ctx context.Context) (structs.NcUserWillWorldSelectAck, error) {
	select {
	case <-ctx.Done():
		return structs.NcUserWillWorldSelectAck{}, errCC
	default:
		otp := randStringBytesMaskImprSrcUnsafe(32)
		err := redisClient.Set(otp, otp, 20*time.Second).Err()
		if err != nil {
			return structs.NcUserWillWorldSelectAck{}, err
		}
		nc := structs.NcUserWillWorldSelectAck{
			Error: 7768,
			Otp:   structs.Name8{},
		}
		copy(nc.Otp.Name[:], otp)

		return nc, nil
	}
}
