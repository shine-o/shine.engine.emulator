package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_ACT_CHAT_REQ
//{
//	char itemLinkDataCount;
//	char len;
//	char content[];
//};

type NcActChatReq struct {
	ItemLinkDataCount byte   `struct:"byte"`
	Length            byte   `struct:"byte"`
	Content           []byte `struct:"sizefrom=Length"`
}

func (nc *NcActChatReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcActChatReq) PdbType() string {
	return `
	struct PROTO_NC_ACT_CHAT_REQ
	{
	  char itemLinkDataCount;
	  char len;
	  char content[];
	};
`
}
