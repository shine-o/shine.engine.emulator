package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_HOLY_PROMISE_LIST_CMD
//{
//  PROTO_HOLY_PROMISE_INFO UpInfo;
//  char nPart;
//  unsigned __int16 MemberCount;
//  PROTO_HOLY_PROMISE_INFO MemberInfo[];
//};
type NcHolyPromiseListCmd struct {
	UpInfo      HolyPromiseInfo
	Part        byte
	MemberCount uint16
	Members     []HolyPromiseInfo `struct:"sizefrom=MemberCount"`
}

func (nc *NcHolyPromiseListCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcHolyPromiseListCmd) PdbType() string {
	return `
	struct PROTO_NC_HOLY_PROMISE_LIST_CMD
	{
	  PROTO_HOLY_PROMISE_INFO UpInfo;
	  char nPart;
	  unsigned __int16 MemberCount;
	  PROTO_HOLY_PROMISE_INFO MemberInfo[];
	};
`
}
