// auxiliar structs
package structs

// struct PROTO_NC_USER_LOGIN_ACK::WorldInfo
type WorldInfo struct {
	WorldNumber byte
	WorldName   Name4
	WorldStatus byte
}

// struct SHINE_XY_TYPE
type ShineXYType struct {
	X uint32
	Y uint32
}

// struct SHINE_DATETIME
type ShineDateTime struct {
	BF0 int32
}

//union NETCOMMAND
type NetCommand struct {
	Protocol uint16
}
