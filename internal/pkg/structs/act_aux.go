package structs

//struct PROTO_NC_ACT_SOMEONEMOVEWALK_CMD::<unnamed-type-moveattr>
type NcActSomeoneMoveWalkCmdAttr struct {
	BF0 int16
}

//union PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-speaker>
type NcActSomeoneShoutCmdSpeaker struct {
	Data [20]byte
}

//struct PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-flag>
type NcActSomeoneShoutCmdFlag struct {
	BF0 byte
}

//struct CHARBRIEFINFO_NOTCAMP
type CharBriefInfoNotCamped struct {
	Equip ProtoEquipment
}
