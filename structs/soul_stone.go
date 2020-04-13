package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD
//{
//  unsigned __int16 player;
//};
type NcSoulStoneSpSomeoneUseCmd struct {
	Player uint16
}

func (nc *NcSoulStoneSpSomeoneUseCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcSoulStoneSpSomeoneUseCmd) PdbType() string {
	return `
	struct PROTO_NC_SOULSTONE_SP_SOMEONEUSE_CMD
	{
	  unsigned __int16 player;
	};
`
}

//struct PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD
//{
//  unsigned __int16 player;
//};
type NcSoulStoneHpSomeoneUseCmd struct {
	Player uint16
}

func (nc *NcSoulStoneHpSomeoneUseCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcSoulStoneHpSomeoneUseCmd) PdbType() string {
	return `
	struct PROTO_NC_SOULSTONE_HP_SOMEONEUSE_CMD
	{
	  unsigned __int16 player;
	};
`
}