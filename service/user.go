package service

import (
	"context"
	protocol "github.com/shine-o/shine.engine.protocol"
	"strings"
)

// RE client struct:
// struct PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ
// {
//  char sVersionKey[64];
// };
type ncUserClientVersionCheckReq struct {
	VersionKey [64]byte
}

// RE client struct:
// struct __cppobj PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK
// {
// };
type ncUserClientRightversionCheckAck struct{}

// RE client struct:
// struct PROTO_NC_USER_US_LOGIN_REQ
// {
//  char sUserName[260];
//  char sPassword[36];
//  Name5 spawnapps;
// };
type ncUserUsLoginReq struct {
	UserName  [260]byte
	Password  [36]byte
	SpawnApps [0]ComplexName1
}

func (nc *ncUserUsLoginReq) authenticate(ctx context.Context) {
	un := nc.UserName[:]
	pass := nc.Password[:]
	userName := strings.TrimRight(string(un), "\x00")
	password := strings.TrimRight(string(pass), "\x00")
	if userName == "admin" && password == "21232f297a57a5a743894a0e4a801fc3" { // temporary :)
		go userLoginAck(ctx, &protocol.Command{})
	} else {
		go userLoginFailAck(ctx, &protocol.Command{})
	}
}

// RE client struct:
// struct PROTO_NC_USER_XTRAP_REQ
// {
//  char XTrapClientKeyLength;
//  char XTrapClientKey[];
// };
type ncUserXtrapReq struct{}

// RE client struct:
// struct PROTO_NC_USER_XTRAP_ACK
// {
//  char bSuccess;
// };
type ncUserXtrapAck struct{}

// RE client struct:
// struct __unaligned __declspec(align(1)) PROTO_NC_USER_LOGIN_ACK
// {
//  char numofworld;
//  PROTO_NC_USER_LOGIN_ACK::WorldInfo worldinfo[];
// };
type ncUserLoginAck struct {
	NumOfWorld byte
	Worlds     [1]WorldInfo
}

func (nc *ncUserLoginAck) setServerInfo(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		// in future communicate with World Service in order to get correct info about IDs of the servers, which will have a linked IP Address
		var worlds [1]WorldInfo
		w1 := WorldInfo{
			WorldNumber: 0,
			WorldName:   ComplexName{},
			WorldStatus: 1,
		}
		copy(w1.WorldName.Name[:], "INITIO")
		copy(w1.WorldName.NameCode[:], []uint16{262, 16720, 17735, 76})
		worlds[0] = w1

		nc.NumOfWorld = byte(1)
		nc.Worlds = worlds
	}
}

// RE client struct:
// struct PROTO_NC_USER_LOGINFAIL_ACK
// {
//	unsigned __int16 error;
// };
type ncUserLoginFailAck struct {
	Err uint16
}

// RE client struct:
// struct __cppobj PROTO_NC_USER_WORLD_STATUS_REQ
// {
// };
type ncUserWorldStatusReq struct{}

// OPERATION CODE ONLY 3100
type ncUserWorldStatusAck struct{}

// RE client struct:
//struct PROTO_NC_USER_WORLDSELECT_REQ
//{
//char worldno;
//};
type ncUserWorldSelectReq struct {
	WorldNo byte
}

// RE client struct:
// struct __unaligned __declspec(align(1)) PROTO_NC_USER_WORLDSELECT_ACK
// {
//	char worldstatus;
//	Name4 ip;
//	unsigned __int16 port;
//	unsigned __int16 validate_new[32];
//};
type ncUserWorldSelectAck struct {
	WorldStatus byte
	Ip          ComplexName
	Port        uint16
	ValidateNew [32]uint16
}
