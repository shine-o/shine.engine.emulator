package structs

import (
	"encoding/json"
	"reflect"
)

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
	Title       [128]byte
	Priority    byte
	NpcHandle   uint16
	NpcPosition ShineXYType
	LimitRange  uint16
	MenuNumber  byte
	Menu        []ServerMenu `struct:"sizefrom=MenuNumber"`
}

func (nc *NcServerMenuReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcServerMenuReq) PdbType() string {
	return `
	struct PROTO_NC_MENU_SERVERMENU_REQ
	{
	  char title[128];
	  char priority;
	  unsigned __int16 npcHandle;
	  SHINE_XY_TYPE npcPosition;
	  unsigned __int16 limitRange;
	  char menunum;
	  SERVERMENU menu[];
	};
`
}
