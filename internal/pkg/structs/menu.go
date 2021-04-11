package structs

//struct PROTO_NC_MENU_SERVERMENU_REQ
type NcServerMenuReq struct {
	Title       string `struct:"[128]byte"`
	Priority    byte
	NpcHandle   uint16
	NpcPosition ShineXYType
	LimitRange  uint16
	MenuNumber  byte
	Menu        []ServerMenu `struct:"sizefrom=MenuNumber"`
}

//struct PROTO_NC_MENU_SERVERMENU_ACK
type NcServerMenuAck struct {
	Reply byte
}
