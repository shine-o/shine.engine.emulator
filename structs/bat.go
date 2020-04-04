package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_BAT_ABSTATERESET_CMD
//{
//  unsigned __int16 handle;
//  ABSTATEINDEX abstate;
//};
type NcBatAbstateResetCmd struct {
	Handle       uint16
	AbstateIndex uint32
}

func (nc *NcBatAbstateResetCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatAbstateResetCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_ABSTATERESET_CMD
	{
	  unsigned __int16 handle;
	  ABSTATEINDEX abstate;
	};
`
}

//struct PROTO_NC_BAT_SPCHANGE_CMD
//{
//  unsigned int sp;
//};
type NcBatSpChangeCmd struct {
	SP uint32
}

func (nc *NcBatSpChangeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBatSpChangeCmd) PdbType() string {
	return `
	struct PROTO_NC_BAT_SPCHANGE_CMD
	{
	  unsigned int sp;
	};
`
}