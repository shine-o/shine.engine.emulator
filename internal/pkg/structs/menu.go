package structs

//struct PROTO_NC_MENU_SERVERMENU_REQ
//{
//  char title[128];
//  char priority;
//  unsigned __int16 npcHandle;
//  SHINE_XY_TYPE npcPosition;
//  unsigned __int16 limitRange;
//  char menunum;
//  SERVERMENU menu[];
//};
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
//{
//  char reply;
//};
type NcServerMenuAck struct {
	Reply byte
}
