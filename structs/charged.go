package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_CHARGED_BOOTHSLOTSIZE_CMD
//{
//  char boothsize;
//};
type NcChargedBoothSlotSizeCmd struct {
	BoothSize byte
}

func (nc *NcChargedBoothSlotSizeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcChargedBoothSlotSizeCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_OPTION_GET_WINDOWPOS_ACK
	{
	  char bSuccess;
	  PROTO_NC_CHAR_OPTION_WINDOWPOS Data;
	};
`
}
