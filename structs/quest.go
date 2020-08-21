package structs

//struct PROTO_NC_QUEST_START_REQ
//{
//  unsigned __int16 nQuestID;
//};
type NcQuestStartReq struct {
	QuestID uint16
}

//struct PROTO_NC_QUEST_SCRIPT_CMD_ACK
//{
//  unsigned __int16 nQuestID;
//  char nQSC;
//  unsigned int nResult;
//};
type NcQuestScriptCmdAck struct {
	QuestID uint16
	QSC     byte
	Result  uint32
}

//struct PROTO_NC_QUEST_RESET_TIME_CMD
//{
//  int tResetYearQuest;
//  int tResetMonthQuest;
//  int tResetWeekQuest;
//  int tResetDayQuest;
//};
type NcQuestResetTimeCmd struct {
	ResetYear  int32
	ResetMonth int32
	ResetWeek  int32
	ResetDay   int32
}

//NC_QUEST_RESET_TIME_CLIENT_CMD
type NcQuestResetTimeClientCmd NcQuestResetTimeCmd
