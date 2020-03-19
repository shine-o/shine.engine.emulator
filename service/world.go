package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.structs"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var errCC = errors.New("context was canceled")

var errHF = errors.New("hardcoded feature")

// WorldCommand wrapper for networking command
// any information scoped to this service and its handlers can be added here
type WorldCommand struct {
	pc *networking.Command
}

func (wc *WorldCommand) worldTime(ctx context.Context) ([]byte, error) {
	var data []byte
	select {
	case <-ctx.Done():
		return data, errCC
	default:
		var (
			t                    time.Time
			hour, minute, second byte
		)

		t = time.Now()
		hour = byte(t.Hour())
		minute = byte(t.Minute())
		second = byte(t.Second())

		nc := &structs.NcMiscGameTimeAck{
			Hour:   hour,
			Minute: minute,
			Second: second,
		}

		data, err := networking.WriteBinary(nc)
		if err != nil {
			return data, err
		}
		return data, nil
	}
}

// user wants to log to given world
// check if world is okay
// take user name, persist to redis
func (wc *WorldCommand) loginToWorld(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return errCC
	default:
		if ncs, ok := wc.pc.NcStruct.(structs.NcUserLoginWorldReq); ok {
			wsi := ctx.Value(networking.ShineSession)
			ws := wsi.(*session)
			userName := strings.TrimRight(string(ncs.User.Name[:]), "\x00")
			ws.UserName = userName
			if err := persistSession(ws); err != nil {
				return err
			}
			return nil
		}
		return fmt.Errorf("unexpected struct type: %v", reflect.TypeOf(wc.pc.NcStruct).String())
	}
}

func (wc *WorldCommand) userWorldInfo(ctx context.Context) ([]byte, error) {
	var data []byte
	select {
	case <-ctx.Done():
		return data, errCC
	default:
		// world id is in the session
		// user name is in the session
		wsi := ctx.Value(networking.ShineSession)
		ws := wsi.(*session)

		if ws.UserName == "admin" { // no database for now, so I hardcode the avatar info
			worldID, err := strconv.Atoi(ws.WorldID)
			if err != nil {
				return data, err
			}
			nc := structs.NcUserLoginWorldAck{
				WorldManager: uint16(worldID),
				NumOfAvatar:  0,
			}

			b, err := networking.WriteBinary(nc.WorldManager)
			if err != nil {
				return data, err
			}
			data = append(data, b...)
			data = append(data, nc.NumOfAvatar)
			return data, nil
		}
		return data, errHF
	}
}

// user clicked previous
// generate a otp token and store it in redis
// login service will use the token to authenticate the user and send him to server select
func (wc *WorldCommand) returnToServerSelect() ([]byte, error) {
	var data []byte
	return data, nil
}
