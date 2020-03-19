package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.structs"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	log = logger.Init("test logger", true, false, ioutil.Discard)
	log.Info("test logger")
	if path, err := filepath.Abs("../defaults"); err != nil {
		log.Fatal(err)
	} else {
		viper.AddConfigPath(path)
		viper.SetConfigType("yaml")

		viper.SetConfigName(".world.circleci")
		// for running tests locally, use this:
		//viper.SetConfigName(".world.test")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}

		viper.SetDefault("serve.port", 9010)
		viper.SetDefault("crypt.xorKey", "0759694a941194858c8805cba09ecd583a365b1a6a16febddf9402f82196c8e99ef7bfbdcfcdb27a009f4022fc11f90c2e12fba7740a7d78401e2ca02d06cba8b97eefde49ea4e13161680f43dc29ad486d7942417f4d665bd3fdbe4e10f50f6ec7a9a0c273d2466d322689c9a520be0f9a50b25da80490dfd3e77d156a8b7f40f9be80f5247f56f832022db0f0bb14385c1cba40b0219dff08becdb6c6d66ad45be89147e2f8910b89360d860def6fe6e9bca06c1759533cfc0b2e0cca5ce12f6e5b5b426c5b2184f2a5d261b654df545c98414dc7c124b189cc724e73c64ffd63a2cee8c8149396cb7dcbd94e232f7dd0afc020164ec4c940ab156f5c9a934de0f3827bc81300f7b3825fee83e29ba5543bf6b9f1f8a4952187f8af888245c4fe1a830878e501f2fd10cb4fd0abcdc1285e252ee4a5838abffc63db960640ab450d54089179ad585cfec0d7e817fe3c3040122ec27ccfa3e21a654c8de00b6df279ff625340785bfa7a5a5e0830c3d5d2040af60a36456f305c41c7d3798c3e85a6e5885a49a6b6af4a37b619b09401e604b32d951a4fef95d4e4afb4ad47c330233d59dce5baa5a7cd8f805fa1f2b8c725750ae6c1989ca01fcfc299b61126863654626c45b50aa2bbeef9a790223752c2013fdd95a7623f10bb5b859f99f7ae606e9a53ab450bf165898b39a6e36ee8deb")
		viper.SetDefault("crypt.xorLimit", 350)
		viper.SetDefault("crypt.client_version", "ebe951290be5b45f1fb075f5505a90b6")
		viper.SetDefault("protocol.nc-data", "defaults/protocol-commands.yml")
	}

	initRedis()
	go selfRPC(ctx)
	os.Exit(m.Run())
}

func TestWorldTime(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// sniffer output -> {"packetType":"small","length":2,"department":2,"command":"D","opCode":2061,"data":"","rawData":"020d08","friendlyName":"NC_MISC_GAMETIME_REQ"}
	if data, err := hex.DecodeString("020d08"); err != nil {
		t.Error(err)
	} else {
		if pc, err := networking.DecodePacket("small", 2, data[1:]); err != nil { // as previous bytes are length info for this packet
			t.Error(err)
		} else {
			nc := structs.NcMiscGameTimeAck{}
			if err := networking.ReadBinary(pc.Base.Data, &nc); err != nil {
				t.Error(err)
			} else {
				pc.NcStruct = nc
				wc := WorldCommand{
					pc: &pc,
				}
				if data, err := wc.worldTime(ctx); err != nil {
					t.Error(err)
				} else {
					// expected network command
					enc := structs.NcMiscGameTimeAck{}
					if err := networking.ReadBinary(data, &enc); err != nil {
						t.Error(err)
					}
				}
			}
		}
	}
}

func TestLoginToWorld(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	s := &session{
		Id:       "bcd1fde6-f9d0-451d-a4b6-4992bd6207e1",
		WorldId:  "1",
		UserName: "admin",
	}
	ctx  = context.WithValue(ctx, "session", s)

	defer cancel()
	// sniffer output -> {"packetType":"big","length":322,"department":3,"command":"F","opCode":3087,"data":"61646d696e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000050410000474500004c00000001e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868","rawData":"0042010f0c61646d696e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000050410000474500004c00000001e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868","friendlyName":"NC_USER_LOGINWORLD_REQ"}
	if data, err := hex.DecodeString("0042010f0c61646d696e000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000050410000474500004c00000001e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868e80868"); err != nil {
		t.Error(err)
	} else {
		if pc, err := networking.DecodePacket("big", 322, data[3:]); err != nil { // as previous bytes are length info for this packet
			t.Error(err)
		} else {
			nc := structs.NcUserLoginWorldReq{}
			if err := networking.ReadBinary(pc.Base.Data, &nc); err != nil {
				t.Error(err)
			} else {
				pc.NcStruct = nc
				wc := WorldCommand{
					pc: &pc,
				}
				if err := wc.loginToWorld(ctx); err != nil {
					t.Error(err)
				} else {
					// check redis for key world-isya-admin
					if r, err := redisClient.Get("admin-world").Result(); err != nil {
						t.Error(err)
					} else {
						s := session{}
						if err := json.Unmarshal([]byte(r), &s); err != nil {
							t.Error(err)
						} else {
							if s.UserName != "admin" {
								t.Errorf("session user: %v is not the xpected one %v", s.UserName, "admin")
							}
						}
					}
				}
			}
		}
	}
}

func TestUserWorldInfo(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// session key is required in context
	s := &session{
		Id:       "bcd1fde6-f9d0-451d-a4b6-4992bd6207e1",
		WorldId:  "1",
		UserName: "admin",
	}

	ctx  = context.WithValue(ctx, "session", s)
	defer cancel()

	pc := &networking.Command{
		Base:     networking.CommandBase{
			OperationCode: 3092,
	}}

	wc := &WorldCommand{pc:pc}

	if data, err := wc.userWorldInfo(ctx); err != nil {
		t.Error(err)
	} else {
		var (
			worldId uint16
			numOfCharacters byte
		)

		buf := bytes.NewBuffer(data)

		if err := binary.Read(buf, binary.LittleEndian, &worldId); err != nil {
			t.Error(err)
		} else {
			if worldId != uint16(1) {
				t.Errorf("result nc.WorldManager: %v is diferent that the expected nc.WorldManager: %v", worldId, 1)
			}
		}

		if err := binary.Read(buf, binary.LittleEndian, &numOfCharacters); err != nil {
			t.Error(err)
		} else {
			if numOfCharacters != byte(0) {
				t.Errorf("result nc.NumOfAvatar: %v is diferent that the expected nc.NumOfAvatar: %v", numOfCharacters, 0)
			}
		}
	}
}