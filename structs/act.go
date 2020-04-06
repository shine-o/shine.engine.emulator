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

//struct PROTO_NC_ACT_STOP_REQ
//{
//  SHINE_XY_TYPE loc;
//};
type NcActStopReq struct {
	Location ShineXYType
}

func (nc *NcActStopReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcActStopReq) PdbType() string {
	return `
	struct PROTO_NC_ACT_STOP_REQ
	{
	  SHINE_XY_TYPE loc;
	};
`
}

//struct PROTO_NC_ACT_MOVESPEED_CMD
//{
//  unsigned __int16 walkspeed;
//  unsigned __int16 runspeed;
//};
type NcActMoveSpeedCmd struct {
	WalkSpeed uint16
	RunSpeed  uint16
}

func (nc *NcActMoveSpeedCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcActMoveSpeedCmd) PdbType() string {
	return `
	struct PROTO_NC_ACT_MOVESPEED_CMD
	{
	  unsigned __int16 walkspeed;
	  unsigned __int16 runspeed;
	};
`
}

//struct PROTO_NC_ACT_SOMEONEMOVEWALK_CMD
//{
//  unsigned __int16 handle;
//  SHINE_XY_TYPE from;
//  SHINE_XY_TYPE to;
//  unsigned __int16 speed;
//  PROTO_NC_ACT_SOMEONEMOVEWALK_CMD::<unnamed-type-moveattr> moveattr;
//};
type NcActSomeoneMoveWalkCmd struct {
	Handle   uint16
	From     ShineXYType
	To       ShineXYType
	Speed    uint16
	MoveAttr NcActSomeoneMoveWalkCmdAttr
}

func (nc *NcActSomeoneMoveWalkCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcActSomeoneMoveWalkCmd) PdbType() string {
	return `
	struct PROTO_NC_ACT_SOMEONEMOVEWALK_CMD
	{
	  unsigned __int16 handle;
	  SHINE_XY_TYPE from;
	  SHINE_XY_TYPE to;
	  unsigned __int16 speed;
	  PROTO_NC_ACT_SOMEONEMOVEWALK_CMD::<unnamed-type-moveattr> moveattr;
	};
`
}

//struct PROTO_NC_ACT_GATHERSTART_REQ
//{
//  unsigned __int16 objhandle;
//};
type NcActGatherStartReq struct {
	Handle uint16
}

func (nc *NcActGatherStartReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcActGatherStartReq) PdbType() string {
	return `
	struct PROTO_NC_ACT_GATHERSTART_REQ
	{
	  unsigned __int16 objhandle;
	};
`
}


//struct PROTO_NC_ACT_SOMEONESHOUT_CMD
//{
//  char itemLinkDataCount;
//  PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-speaker> speaker;
//  PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-flag> flag;
//  char len;
//  char content[];
//};
type NcActSomeoneShoutCmd struct {
	Count byte
	Speaker NcActSomeoneShoutCmdSpeaker
	Flag NcActSomeoneShoutCmdFlag
	Len byte
	Content []byte `struct:"sizefrom=Len"`
}

func (nc *NcActSomeoneShoutCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcActSomeoneShoutCmd) PdbType() string {
	return `
	struct PROTO_NC_ACT_SOMEONESHOUT_CMD
	{
	  char itemLinkDataCount;
	  PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-speaker> speaker;
	  PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-flag> flag;
	  char len;
	  char content[];
	};
`
}

//struct PROTO_NC_ACT_SOMEONESTOP_CMD
//{
//  unsigned __int16 handle;
//  SHINE_XY_TYPE loc;
//};
type NcActSomeoneStopCmd struct {
	Handle uint16
	Location ShineXYType
}

func (nc *NcActSomeoneStopCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcActSomeoneStopCmd) PdbType() string {
	return `
	struct PROTO_NC_ACT_SOMEONESTOP_CMD
	{
	  unsigned __int16 handle;
	  SHINE_XY_TYPE loc;
	};
`
}