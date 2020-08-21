package structs

//struct PROTO_NC_HOLY_PROMISE_LIST_CMD
//{
//  PROTO_HOLY_PROMISE_INFO UpInfo;
//  char nPart;
//  unsigned __int16 MemberCount;
//  PROTO_HOLY_PROMISE_INFO MemberInfo[];
//};
type NcHolyPromiseListCmd struct {
	UpInfo      HolyPromiseInfo
	Part        byte
	MemberCount uint16
	Members     []HolyPromiseInfo `struct:"sizefrom=MemberCount"`
}
