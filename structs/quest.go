package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_QUEST_START_REQ
//{
//  unsigned __int16 nQuestID;
//};
type NcQuestStartReq struct {
	QuestID uint16
}

func (nc *NcQuestStartReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcQuestStartReq) PdbType() string {
	return `
	struct PROTO_NC_COLLECT_CARDREGIST_REQ
	{
	  char invenslot;
	};
`
}

//struct PROTO_NC_QUEST_SCRIPT_CMD_ACK
//{
//  unsigned __int16 nQuestID;
//  char nQSC;
//  unsigned int nResult;
//};
type NcQuestScriptCmdAck struct {
	QuestID uint16
	QSC byte
	Result uint32
}

func (nc *NcQuestScriptCmdAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcQuestScriptCmdAck) PdbType() string {
	return `
	struct PROTO_NC_QUEST_SCRIPT_CMD_ACK
	{
	  unsigned __int16 nQuestID;
	  char nQSC;
	  unsigned int nResult;
	};
`
}