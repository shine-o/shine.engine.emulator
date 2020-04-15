package structs

//struct PROTO_NC_KQ_TEAM_TYPE_CMD
//{
//  char nTeamType;
//};
type NcKqTeamTypeCmd struct {
	TeamType byte
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