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

var CC = errors.New("context was canceled")

var HF = errors.New("hardcoded feature")


type WorldCommand struct {
	pc *networking.Command
}

func (wc * WorldCommand) worldTime(ctx context.Context) ([]byte, error) {
	var data []byte
	select {
	case <- ctx.Done():
		return data, CC
	default:
		var (
			t time.Time
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

		if data, err := networking.WriteBinary(nc); err != nil {
			return data, err
		} else {
			return data, nil
		}
	}
}

// user wants to log to given world
// check if world is okay
// take user name, persist to redis
func (wc * WorldCommand) loginToWorld(ctx context.Context) error {
	select {
	case <- ctx.Done():
		return CC
	default:
		if ncs, ok := wc.pc.NcStruct.(structs.NcUserLoginWorldReq); ok {
			wsi := ctx.Value("session")
			ws := wsi.(*session)
			userName := strings.TrimRight(string(ncs.User.Name[:]), "\x00")
			ws.UserName = userName
			if err := persistSession(ws); err != nil {
				return err
			} else {
				return nil
			}
		} else {
			return fmt.Errorf("unexpected struct type: %v", reflect.TypeOf(wc.pc.NcStruct).String())
		}
	}
}

func (wc * WorldCommand) userWorldInfo(ctx context.Context) ([]byte, error)  {
	var data []byte
	select {
	case <-ctx.Done():
		return data, CC
	default:
		// world id is in the session
		// user name is in the session
		wsi := ctx.Value("session")
		ws := wsi.(*session)

		if ws.UserName == "admin" { // no database for now, so I hardcode the avatar info
			if worldId, err := strconv.Atoi(ws.WorldId); err != nil {
				return data, err
			} else {

				nc := structs.NcUserLoginWorldAck{
					WorldManager: uint16(worldId),
					NumOfAvatar:  0,
				}

				if b, err := networking.WriteBinary(nc.WorldManager); err != nil {
					return data, err
				} else {
					data = append(data, b...)
					data = append(data, nc.NumOfAvatar)
				}
			}
			return data, nil
		} else {
			return data, HF
		}
	}
}

// user clicked previous
// generate a otp token and store it in redis
// login service will use the token to authenticate the user and send him to server select
func (wc * WorldCommand) returnToServerSelect() ([]byte, error) {
	var data []byte
	return data, nil
}