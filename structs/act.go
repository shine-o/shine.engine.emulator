package structs

//struct PROTO_NC_ACT_CHAT_REQ
//{
//	char itemLinkDataCount;
//	char len;
//	char content[];
//};
type NcActChatReq struct {
	ItemLinkDataCount byte
	Length            byte
	Content           []byte `struct:"sizefrom=Length"`
	ItemMetadata []byte `struct-while:"!_eof"`
}

//struct PROTO_NC_ACT_STOP_REQ
//{
//  SHINE_XY_TYPE loc;
//};
type NcActStopReq struct {
	Location ShineXYType
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

//struct PROTO_NC_ACT_GATHERSTART_REQ
//{
//  unsigned __int16 objhandle;
//};
type NcActGatherStartReq struct {
	Handle uint16
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
	Count   byte
	Speaker NcActSomeoneShoutCmdSpeaker
	Flag    NcActSomeoneShoutCmdFlag
	Len     byte
	Content []byte `struct:"sizefrom=Len"`
	// data for viewing the item listed in the chat
	ItemMetadata []byte `struct-while:"!_eof"`
}

//struct PROTO_NC_ACT_SOMEONESTOP_CMD
//{
//  unsigned __int16 handle;
//  SHINE_XY_TYPE loc;
//};
type NcActSomeoneStopCmd struct {
	Handle   uint16
	Location ShineXYType
}

//struct PROTO_NC_ACT_NPCCLICK_CMD
//{
//  unsigned __int16 npchandle;
//};
type NcActNpcClickCmd struct {
	NpcHandle uint16
}

//struct PROTO_NC_ACT_CHANGEMODE_REQ
//{
//  char mode;
//};
type NcActChangeModeReq struct {
	Mode byte
}

//struct PROTO_NC_ACT_SOMEONEPRODUCE_CAST_CMD
//{
//  unsigned __int16 caster;
//  unsigned __int16 item;
//};
type NcActSomeoneProduceCastCmd struct {
	Caster uint16
	Item   uint16
}

//struct PROTO_NC_ACT_SOMEONEFOLDTENT_CMD
//{
//  unsigned __int16 handle;
//  CHARBRIEFINFO_NOTCAMP shape;
//};
type NcActSomeoneFoldTentCmd struct {
	Handle uint16
	Shape  CharBriefInfoNotCamped
}

//struct PROTO_NC_ACT_SOMEONECHANGEMODE_CMD
//{
//  unsigned __int16 handle;
//  char mode;
//};
type NcActSomeoneChangeModeCmd struct {
	Handle uint16
	Mode   byte
}

//struct PROTO_NC_ACT_SOMEONEPRODUCE_MAKE_CMD
//{
//  unsigned __int16 caster;
//  unsigned __int16 item;
//};
type NcActSomeoneProduceMakeCmd struct {
	Caster uint16
	Item   uint16
}

//struct PROTO_NC_ACT_SOMEEONEJUMP_CMD
//{
//  unsigned __int16 handle;
//};
type NcActSomeoneJumpCmd struct {
	Handle uint16
}

//struct PROTO_NC_ACT_WALK_REQ
//{
//  SHINE_XY_TYPE from;
//  SHINE_XY_TYPE to;
//};
type NcActWalkReq struct {
	From ShineXYType
	To   ShineXYType
}

//NC_ACT_MOVERUN_CMD
type NcActMoveRunCmd NcActWalkReq

//struct PROTO_NC_ACT_SOMEONEMOVEWALK_CMD
//{
//  unsigned __int16 handle;
//  SHINE_XY_TYPE from;
//  SHINE_XY_TYPE to;
//  unsigned __int16 speed;
//  PROTO_NC_ACT_SOMEONEMOVEWALK_CMD::<unnamed-type-moveattr> moveattr;
//};
type NcActSomeoneWalkCmd struct {
	Handle uint16
	From   ShineXYType
	To     ShineXYType
	Speed uint16 // not sure
	MoveAttr uint16
}

//NC_ACT_SOMEONEMOVERUN_CMD
type NcActSomeoneMoveRunCmd NcActSomeoneWalkCmd
