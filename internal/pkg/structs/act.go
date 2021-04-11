package structs

//struct PROTO_NC_ACT_CHAT_REQ
type NcActChatReq struct {
	ItemLinkDataCount byte
	Length            byte
	Content           []byte `struct:"sizefrom=Length"`
	ItemMetadata      []byte `struct-while:"!_eof"`
}

//struct PROTO_NC_ACT_STOP_REQ
type NcActStopReq struct {
	Location ShineXYType
}

//struct PROTO_NC_ACT_MOVESPEED_CMD
type NcActMoveSpeedCmd struct {
	WalkSpeed uint16
	RunSpeed  uint16
}

//struct PROTO_NC_ACT_SOMEONEMOVEWALK_CMD
type NcActSomeoneMoveWalkCmd struct {
	Handle   uint16
	From     ShineXYType
	To       ShineXYType
	Speed    uint16
	MoveAttr NcActSomeoneMoveWalkCmdAttr
}

//struct PROTO_NC_ACT_GATHERSTART_REQ
type NcActGatherStartReq struct {
	Handle uint16
}

//struct PROTO_NC_ACT_SOMEONESHOUT_CMD
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
type NcActSomeoneStopCmd struct {
	Handle   uint16
	Location ShineXYType
}

//struct PROTO_NC_ACT_NPCCLICK_CMD
type NcActNpcClickCmd struct {
	NpcHandle uint16
}

//struct PROTO_NC_ACT_CHANGEMODE_REQ
type NcActChangeModeReq struct {
	Mode byte
}

//struct PROTO_NC_ACT_SOMEONEPRODUCE_CAST_CMD
type NcActSomeoneProduceCastCmd struct {
	Caster uint16
	Item   uint16
}

//struct PROTO_NC_ACT_SOMEONEFOLDTENT_CMD
type NcActSomeoneFoldTentCmd struct {
	Handle uint16
	Shape  CharBriefInfoNotCamped
}

//struct PROTO_NC_ACT_SOMEONECHANGEMODE_CMD
type NcActSomeoneChangeModeCmd struct {
	Handle uint16
	Mode   byte
}

//struct PROTO_NC_ACT_SOMEONEPRODUCE_MAKE_CMD
type NcActSomeoneProduceMakeCmd struct {
	Caster uint16
	Item   uint16
}

//struct PROTO_NC_ACT_SOMEEONEJUMP_CMD
type NcActSomeoneJumpCmd struct {
	Handle uint16
}

//struct PROTO_NC_ACT_WALK_REQ
type NcActWalkReq struct {
	From ShineXYType
	To   ShineXYType
}

//NC_ACT_MOVERUN_CMD
type NcActMoveRunCmd NcActWalkReq

//struct PROTO_NC_ACT_SOMEONEMOVEWALK_CMD
type NcActSomeoneWalkCmd struct {
	Handle   uint16
	From     ShineXYType
	To       ShineXYType
	Speed    uint16 // not sure
	MoveAttr uint16
}

//NC_ACT_SOMEONEMOVERUN_CMD
type NcActSomeoneMoveRunCmd NcActSomeoneWalkCmd

//struct PROTO_NC_ACT_MOVEWALK_CMD
type NcActMoveWalkCmd struct {
	From ShineXYType
	To   ShineXYType
}

// NC_ACT_NPCMENUOPEN_REQ
type NcActNpcMenuOpenReq struct {
	MobID uint16
}

// struct PROTO_NC_ACT_NPCMENUOPEN_ACK
type NcActNpcMenuOpenAck struct {
	Ack byte
}
