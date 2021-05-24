package structs

// RE client struct:
// struct PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ
type NcUserClientVersionCheckReq struct {
	VersionKey [64]byte `struct:"[64]byte"`
}

// struct PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK
type NcUserClientWrongVersionCheckAck struct{}

// struct PROTO_NC_USER_US_LOGIN_REQ
type NcUserUsLoginReq struct {
	UserName  string `struct:"[260]byte"`
	Password  string `struct:"[36]byte"`
	SpawnApps Name5
}

type NewUserLoginReq struct {
	Unk1     [32]byte
	UserName string `struct:"[260]byte"`
	Password string `struct:"[37]byte"`
	Unk2     string `struct:"[20]byte"`
}

// struct PROTO_NC_USER_LOGIN_ACK
type NcUserLoginAck struct {
	NumOfWorld byte        `struct:"byte,sizeof=Worlds"`
	Worlds     []WorldInfo `struct:"sizefrom=NumOfWorld"`
}

// struct PROTO_NC_USER_LOGINFAIL_ACK
type NcUserLoginFailAck struct {
	Err uint16
}

// struct PROTO_NC_USER_WORLDSELECT_REQ
type NcUserWorldSelectReq struct {
	WorldNo byte
}

// struct PROTO_NC_USER_WORLDSELECT_ACK
type NcUserWorldSelectAck struct {
	// 1: behaviour -> cannot enter, message -> The server is under maintenance.
	// 2: behaviour -> cannot enter, message -> You cannot connect to an empty server.
	// 3: behaviour -> cannot enter, message -> The server has been reserved for a special use.
	// 4: behaviour -> cannot enter, message -> Login failed due to an unknown error.
	// 5: behaviour -> cannot enter, message -> The server is full.
	WorldStatus byte
	Ip          Name4
	Port        uint16
	ValidateNew [32]uint16 `struct:"[32]uint16"`
}

// struct PROTO_NC_USER_LOGINWORLD_REQ
type NcUserLoginWorldReq struct {
	User        Name256Byte
	ValidateNew [32]uint16 `struct:"[32]uint16"`
}

// struct PROTO_NC_USER_LOGINWORLD_ACK
type NcUserLoginWorldAck struct {
	WorldManager uint16
	NumOfAvatar  byte
	Avatars      []AvatarInformation `struct:"sizefrom=NumOfAvatar"`
}

// struct PROTO_NC_USER_WILL_WORLD_SELECT_ACK
type NcUserWillWorldSelectAck struct {
	Error uint16
	Otp   Name8
}

// struct PROTO_NC_USER_LOGIN_WITH_OTP_REQ
type NcUserLoginWithOtpReq struct {
	Otp Name8
}

// struct PROTO_NC_USER_USE_BEAUTY_SHOP_CMD
type NcUserUseBeautyShopCmd struct {
	Filler byte
}
