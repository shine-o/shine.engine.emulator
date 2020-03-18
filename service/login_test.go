package service

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.structs"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
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

		viper.SetConfigName(".login.circleci")
		// for running tests locally, use this:
		//viper.SetConfigName(".login.test")

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}

		viper.SetDefault("serve.port", 9010)
		viper.SetDefault("crypt.xorKey", "0759694a941194858c8805cba09ecd583a365b1a6a16febddf9402f82196c8e99ef7bfbdcfcdb27a009f4022fc11f90c2e12fba7740a7d78401e2ca02d06cba8b97eefde49ea4e13161680f43dc29ad486d7942417f4d665bd3fdbe4e10f50f6ec7a9a0c273d2466d322689c9a520be0f9a50b25da80490dfd3e77d156a8b7f40f9be80f5247f56f832022db0f0bb14385c1cba40b0219dff08becdb6c6d66ad45be89147e2f8910b89360d860def6fe6e9bca06c1759533cfc0b2e0cca5ce12f6e5b5b426c5b2184f2a5d261b654df545c98414dc7c124b189cc724e73c64ffd63a2cee8c8149396cb7dcbd94e232f7dd0afc020164ec4c940ab156f5c9a934de0f3827bc81300f7b3825fee83e29ba5543bf6b9f1f8a4952187f8af888245c4fe1a830878e501f2fd10cb4fd0abcdc1285e252ee4a5838abffc63db960640ab450d54089179ad585cfec0d7e817fe3c3040122ec27ccfa3e21a654c8de00b6df279ff625340785bfa7a5a5e0830c3d5d2040af60a36456f305c41c7d3798c3e85a6e5885a49a6b6af4a37b619b09401e604b32d951a4fef95d4e4afb4ad47c330233d59dce5baa5a7cd8f805fa1f2b8c725750ae6c1989ca01fcfc299b61126863654626c45b50aa2bbeef9a790223752c2013fdd95a7623f10bb5b859f99f7ae606e9a53ab450bf165898b39a6e36ee8deb")
		viper.SetDefault("crypt.xorLimit", 350)
		viper.SetDefault("crypt.client_version", "ebe951290be5b45f1fb075f5505a90b6")
		viper.SetDefault("protocol.nc-data", "defaults/protocol-commands.yml")

		requiredParams := []string{
			"database.postgres.host",
			"database.postgres.port",
			"database.postgres.db_user",
			"database.postgres.db_password",
		}
		for _, rp := range requiredParams {
			if !viper.IsSet(rp) {
				log.Fatalf("missing required parameter %v", rp)
			}
		}
	}
	initDatabase()
	initRedis()
	//gRpcClients(ctx)
	os.Exit(m.Run())
}

func TestCheckClientVersion(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// sniffer output -> {"packetType":"small","length":66,"department":3,"command":"65","opCode":3173,"data":"656265393531323930626535623435663166623037356635353035613930623600ad5f76d8f61e1630ad5f7678fc1900c5b88d00fc17dd020118dd0200000000","rawData":"42650c656265393531323930626535623435663166623037356635353035613930623600ad5f76d8f61e1630ad5f7678fc1900c5b88d00fc17dd020118dd0200000000","friendlyName":"NC_USER_CLIENT_VERSION_CHECK_REQ"}

	if data, err := hex.DecodeString("42650c656265393531323930626535623435663166623037356635353035613930623600ad5f76d8f61e1630ad5f7678fc1900c5b88d00fc17dd020118dd0200000000"); err != nil {
		t.Error(err)
	} else {

		if pc, err := networking.DecodePacket("small", 66, data[1:]); err != nil { // from index 1, as previous bytes are length info for this packet
			t.Error(err)
		} else {
			nc := structs.NcUserClientVersionCheckReq{}
			if err := networking.ReadBinary(pc.Base.Data, &nc); err != nil {
				t.Error(err)
			} else {
				pc.NcStruct = nc
				lc := LoginCommand{
					pc: &pc,
				}
				if _, err := lc.checkClientVersion(ctx); err != nil {
					t.Error(err)
				}
			}
		}
	}
}

// requires database connection
func TestCheckCredentials(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// sniffer output -> {"packetType":"big","length":318,"department":3,"command":"5A","opCode":3162,"data":"61646d696e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003231323332663239376135376135613734333839346130653461383031666333000000004f726967696e616c000000000000000000000000","rawData":"005a0c5a0c61646d696e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003231323332663239376135376135613734333839346130653461383031666333000000004f726967696e616c000000000000000000000000","friendlyName":"NC_USER_US_LOGIN_REQ"}
	if data, err := hex.DecodeString("003e015a0c61646d696e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003231323332663239376135376135613734333839346130653461383031666333000000004f726967696e616c000000000000000000000000"); err != nil {
		t.Error(err)
	} else {
		if pc, err := networking.DecodePacket("big", 318, data[3:]); err != nil { // from index 3, as previous bytes are length info for this packet
			t.Error(err)
		} else {
			nc := structs.NcUserUsLoginReq{}
			if err := networking.ReadBinary(pc.Base.Data, &nc); err != nil {
				t.Error(err)
			} else {
				pc.NcStruct = nc
				lc := LoginCommand{
					pc: &pc,
				}
				if err := lc.checkCredentials(ctx); err != nil {
					t.Error(err)
				}
			}
		}
	}
}

// depends on external service and manages no data
// for completeness sake, run this test when developing
//func TestCheckWorldStatus(t *testing.T) {
//	// timeout
//	ctx := context.Background()
//	ctx, cancel := context.WithCancel(ctx)
//	defer cancel()
//	lc := LoginCommand{}
//	if err := lc.checkWorldStatus(ctx); err != nil {
//		t.Error(err)
//	}
//}

// depends on external service and data is handled on the server side
// for completeness sake, run this test when developing
//func TestUserSelectedServer(t *testing.T) {
//	ctx := context.Background()
//	ctx, cancel := context.WithCancel(ctx)
//	defer cancel()
//	// sniffer output -> {"packetType":"small","length":3,"department":3,"command":"B","opCode":3083,"data":"00","rawData":"030b0c00","friendlyName":"NC_USER_WORLDSELECT_REQ"}
//	if data, err := hex.DecodeString("030b0c00"); err != nil {
//		t.Error(err)
//	} else {
//		if pc, err := networking.DecodePacket("small", 3, data[1:]); err != nil { // from index 1, as previous bytes are length info for this packet
//			t.Error(err)
//		} else {
//			nc := structs.NcUserUsLoginReq{}
//			if err := networking.ReadBinary(pc.Base.Data, &nc); err != nil {
//				t.Error(err)
//			} else {
//				pc.NcStruct = nc
//				lc := LoginCommand{
//					pc: &pc,
//				}
//				if data, err := lc.userSelectedServer(ctx); err != nil {
//					t.Error(err)
//				} else {
//					rnc := structs.NcUserWorldSelectAck{}
//					if err := networking.ReadBinary(data, &rnc); err != nil {
//						t.Error(err)
//					} else {
//						if rnc.WorldStatus != byte(6) {
//							t.Errorf("unexpected world status %v", rnc.WorldStatus)
//						}
//					}
//				}
//			}
//		}
//	}
//}

//
func TestLoginByCode(t *testing.T) {
	// setup dummy otp token in redis store
	otp := "tEGMohMSNboCclYHGIXUOGHKZTKcjLfr"
	if err := redisClient.Set(otp, otp, 20*time.Second).Err(); err != nil {
		t.Error(err)
	} else {
		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()
		// sniffer output -> {"packetType":"small","length":34,"department":3,"command":"37","opCode":3127,"data":"7445474d6f684d534e626f43636c5948474958554f47484b5a544b636a4c6672","rawData":"22370c7445474d6f684d534e626f43636c5948474958554f47484b5a544b636a4c6672","friendlyName":"NC_USER_LOGIN_WITH_OTP_REQ"}
		if data, err := hex.DecodeString("22370c7445474d6f684d534e626f43636c5948474958554f47484b5a544b636a4c6672"); err != nil {
			t.Error(err)
		} else {
			if pc, err := networking.DecodePacket("small", 34, data[1:]); err != nil { // from index 1, as previous bytes are length info for this packet
				t.Error(err)
			} else {
				nc := structs.NcUserLoginWithOtpReq{}
				if err := networking.ReadBinary(pc.Base.Data, &nc); err != nil {
					t.Error(err)
				} else {
					pc.NcStruct = nc
					lc := LoginCommand{
						pc: &pc,
					}
					if err := lc.loginByCode(ctx); err != nil {
						t.Error(err)
					}
				}
			}
		}
	}
}
