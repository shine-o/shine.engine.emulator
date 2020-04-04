package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_CMD
//{
//  unsigned __int16 handle;
//  ABSTATE_INFORMATION info;
//};
type NcBriefInfoAbstateChangeCmd struct {
	Handle uint16
	Info AbstateInformation
}

func (nc *NcBriefInfoAbstateChangeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBriefInfoAbstateChangeCmd) PdbType() string {
	return `
	struct PROTO_NC_BRIEFINFO_ABSTATE_CHANGE_CMD
	{
	  unsigned __int16 handle;
	  ABSTATE_INFORMATION info;
	};
`
}
