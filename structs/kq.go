package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_KQ_TEAM_TYPE_CMD
//{
//  char nTeamType;
//};
type NcKqTeamTypeCmd struct {
	TeamType byte
}

func (nc *NcKqTeamTypeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcKqTeamTypeCmd) PdbType() string {
	return `
	struct PROTO_NC_KQ_TEAM_TYPE_CMD
	{
	  char nTeamType;
	};
`
}

//struct PROTO_NC_KQ_LIST_TIME_ACK
//{
//  int ServerTime;
//  tm tm_ServerTime;
//};
type NcKqListTimeAck struct {
	ServerTime   int32
	TmServerTime TM
}

func (nc *NcKqListTimeAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc NcKqListTimeAck) PdbType() string {
	return `
	struct PROTO_NC_KQ_LIST_TIME_ACK
	{
	  int ServerTime;
	  tm tm_ServerTime;
	};
`
}
