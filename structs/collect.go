package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_COLLECT_CARDREGIST_REQ
//{
//  char invenslot;
//};
type NcCollectCardRegisterReq struct {
	Slot byte
}

func (nc *NcCollectCardRegisterReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCollectCardRegisterReq) PdbType() string {
	return `
	struct PROTO_NC_COLLECT_CARDREGIST_REQ
	{
	  char invenslot;
	};
`
}
