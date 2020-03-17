package service

import (
	"context"
	"encoding/hex"
	"github.com/shine-o/shine.engine.networking"
	"github.com/shine-o/shine.engine.structs"
	"testing"
)

func TestCheckClientVersion(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// sniffer output
	// {"packetType":"small","length":66,"department":3,"command":"65","opCode":3173,"data":"656265393531323930626535623435663166623037356635353035613930623600ad5f76d8f61e1630ad5f7678fc1900c5b88d00fc17dd020118dd0200000000","rawData":"42650c656265393531323930626535623435663166623037356635353035613930623600ad5f76d8f61e1630ad5f7678fc1900c5b88d00fc17dd020118dd0200000000","friendlyName":"NC_USER_CLIENT_VERSION_CHECK_REQ"}
	if data, err := hex.DecodeString("42650c656265393531323930626535623435663166623037356635353035613930623600ad5f76d8f61e1630ad5f7678fc1900c5b88d00fc17dd020118dd0200000000"); err != nil {
		t.Error(err)
	} else {

		if pc, err := networking.DecodePacket("small", 66, data); err != nil {
			t.Error(err)
		} else {
			nc := &structs.NcUserClientVersionCheckReq{}
			if err := networking.ReadBinary(data, nc); err != nil {
				t.Error(err)
			} else {
				pc.NcStruct = nc
				lc := LoginCommand{
					pc: &pc,
				}
				if data, err := lc.checkClientVersion(ctx); err != nil {
					t.Error(err)
				} else {
					// assert data is what it should be
					if len(data) <= 0 {
						t.Errorf("bad client version")
					}
				}
			}
		}
	}
}

func TestCheckCredentials(t *testing.T) {}

func TestCheckWorldStatus(t *testing.T) {}

func TestUserSelectedServer(t *testing.T) {}

func TestLoginByCode(t *testing.T) {}