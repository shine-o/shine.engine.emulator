package service

import (
	"context"
	networking "github.com/shine-o/shine.engine.networking"
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
		go userLoginAck(ctx, &networking.Command{})
	} else {
		go userLoginFailAck(ctx, &networking.Command{})
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
	Ip          Name4
	Port        uint16
	ValidateNew [32]uint16
}
