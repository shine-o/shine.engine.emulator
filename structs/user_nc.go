package structs

import (
	"encoding/json"
	"reflect"
)

// RE client struct:
// struct PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ
// {
//  char sVersionKey[64];
// };
type NcUserClientVersionCheckReq struct {
	VersionKey [64]byte `struct:"[64]byte"`
}

func (nc *NcUserClientVersionCheckReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserClientVersionCheckReq) PdbType() string {
	return `
	struct PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ
	{
	  char sVersionKey[64];
	};
`
}


// struct PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK
//{
//  char dummy[1];
//};
type NcUserClientWrongVersionCheckAck struct{}

func (nc *NcUserClientWrongVersionCheckAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserClientWrongVersionCheckAck) PdbType() string {
	return `
	struct PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK
	{
	  char dummy[1];
	};
`
}

// RE client struct:
// struct PROTO_NC_USER_US_LOGIN_REQ
// {
//  char sUserName[260];
//  char sPassword[36];
//  Name5 spawnapps;
// };
type NcUserUsLoginReq struct {
	UserName  [260]byte `struct:"[260]byte"`
	Password  [36]byte  `struct:"[36]byte"`
	SpawnApps Name5
}

func (nc *NcUserUsLoginReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserUsLoginReq) PdbType() string {
	return `
	struct PROTO_NC_USER_US_LOGIN_REQ
	{
	  char sUserName[260];
	  char sPassword[36];
	  Name5 spawnapps;
	};
`
}


//struct PROTO_NC_USER_LOGIN_ACK
//{
//  char numofworld;
//  PROTO_NC_USER_LOGIN_ACK::WorldInfo worldinfo[];
//};
type NcUserLoginAck struct {
	NumOfWorld byte        `struct:"byte"`
	Worlds     []WorldInfo `struct:"sizefrom=NumOfWorld"`
}

func (nc *NcUserLoginAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserLoginAck) PdbType() string {
	return `
	struct PROTO_NC_USER_LOGIN_ACK
	{
	  char numofworld;
	  PROTO_NC_USER_LOGIN_ACK::WorldInfo worldinfo[];
	};
`
}


// RE client struct:
// struct PROTO_NC_USER_LOGINFAIL_ACK
// {
//	unsigned __int16 error;
// };
type NcUserLoginFailAck struct {
	Err uint16 `struct:"uint16"`
}

func (nc *NcUserLoginFailAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserLoginFailAck) PdbType() string {
	return `
	struct PROTO_NC_USER_LOGINFAIL_ACK
	{
	  unsigned __int16 error;
	};
`
}

// RE client struct:
//struct PROTO_NC_USER_WORLDSELECT_REQ
//{
//char worldno;
//};
type NcUserWorldSelectReq struct {
	WorldNo byte `struct:"byte"`
}

func (nc *NcUserWorldSelectReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserWorldSelectReq) PdbType() string {
	return `
	struct PROTO_NC_USER_WORLDSELECT_REQ
	{
	  char worldno;
	};
`
}

//struct PROTO_NC_USER_WORLDSELECT_ACK
//{
//  char worldstatus;
//  Name4 ip;
//  unsigned __int16 port;
//  unsigned __int16 validate_new[32];
//};
type NcUserWorldSelectAck struct {
	// 1: behaviour -> cannot enter, message -> The server is under maintenance.
	// 2: behaviour -> cannot enter, message -> You cannot connect to an empty server.
	// 3: behaviour -> cannot enter, message -> The server has been reserved for a special use.
	// 4: behaviour -> cannot enter, message -> Login failed due to an unknown error.
	// 5: behaviour -> cannot enter, message -> The server is full.
	WorldStatus byte `struct:"byte"`
	Ip          Name4
	Port        uint16     `struct:"uint16"`
	ValidateNew [32]uint16 `struct:"[32]uint16"`
}

func (nc *NcUserWorldSelectAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserWorldSelectAck) PdbType() string {
	return `
	struct PROTO_NC_USER_WORLDSELECT_ACK
	{
	  char worldstatus;
	  Name4 ip;
	  unsigned __int16 port;
	  unsigned __int16 validate_new[32];
	};
`
}

// struct PROTO_NC_USER_LOGINWORLD_REQ
// {
//  Name256Byte user;
//  unsigned __int16 validate_new[32];
// };
type NcUserLoginWorldReq struct {
	User        Name256Byte
	ValidateNew [28]uint16 `struct:"[28]uint16"`
}

func (nc *NcUserLoginWorldReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserLoginWorldReq) PdbType() string {
	return `
	struct PROTO_NC_USER_LOGINWORLD_REQ
	{
	  Name256Byte user;
	  unsigned __int16 validate_new[32];
	};
`
}

//struct PROTO_NC_USER_LOGINWORLD_ACK
//{
//  unsigned __int16 worldmanager;
//  char numofavatar;
//  PROTO_AVATARINFORMATION avatar[];
//};
type NcUserLoginWorldAck struct {
	WorldManager uint16              `struct:"uint16"`
	NumOfAvatar  byte                `struct:"byte"`
	Avatars      []AvatarInformation `struct:"sizefrom=NumOfAvatar"`
}

func (nc *NcUserLoginWorldAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserLoginWorldAck) PdbType() string {
	return `
	struct PROTO_NC_USER_LOGINWORLD_ACK
	{
	  unsigned __int16 worldmanager;
	  char numofavatar;
	  PROTO_AVATARINFORMATION avatar[];
	};
`
}

// struct PROTO_NC_USER_WILL_WORLD_SELECT_ACK
// {
//	unsigned __int16 nError;
//	Name8 sOTP;
// };
type NcUserWillWorldSelectAck struct {
	Error uint16 `struct:"uint16"`
	Otp   Name8
}

func (nc *NcUserWillWorldSelectAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserWillWorldSelectAck) PdbType() string {
	return `
	struct PROTO_NC_USER_WILL_WORLD_SELECT_ACK
	{
	  unsigned __int16 nError;
	  Name8 sOTP;
	};
`
}

// struct PROTO_NC_USER_LOGIN_WITH_OTP_REQ
// {
// 	Name8 sOTP;
// };
type NcUserLoginWithOtpReq struct {
	Otp Name8
}

func (nc *NcUserLoginWithOtpReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcUserLoginWithOtpReq) PdbType() string {
	return `
	struct PROTO_NC_USER_LOGIN_WITH_OTP_REQ
	{
	  Name8 sOTP;
	};	
`
}