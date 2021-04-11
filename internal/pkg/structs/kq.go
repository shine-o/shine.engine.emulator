package structs

//struct PROTO_NC_KQ_TEAM_TYPE_CMD
type NcKqTeamTypeCmd struct {
	TeamType byte
}

//struct PROTO_NC_KQ_LIST_TIME_ACK
type NcKqListTimeAck struct {
	ServerTime   int32
	TmServerTime TM
}
