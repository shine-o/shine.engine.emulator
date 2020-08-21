package structs

//struct PROTO_NC_ACT_SOMEONEMOVEWALK_CMD::<unnamed-type-moveattr>
//{
//  __int16 _bf0;
//};
type NcActSomeoneMoveWalkCmdAttr struct {
	BF0 int16
}

//union PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-speaker>
//{
//  char charID[20];
//  unsigned __int16 mobID;
//};
type NcActSomeoneShoutCmdSpeaker struct {
	Data [20]byte
}

//struct PROTO_NC_ACT_SOMEONESHOUT_CMD::<unnamed-type-flag>
//{
//  char _bf0;
//};
type NcActSomeoneShoutCmdFlag struct {
	BF0 byte
}

//struct CHARBRIEFINFO_NOTCAMP
//{
//  PROTO_EQUIPMENT equip;
//};
type CharBriefInfoNotCamped struct {
	Equip ProtoEquipment
}
