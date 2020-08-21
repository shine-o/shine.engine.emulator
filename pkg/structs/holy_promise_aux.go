package structs

//struct PROTO_HOLY_PROMISE_INFO
//{
//  Name5 PromiseCharID;
//  PROTO_HOLY_PROMISE_DATE LogonInfo;
//  char Level;
//  char Flags;
//};
type HolyPromiseInfo struct {
	PromiseCharID Name5
	LogonInfo     HolyPromiseDate
	Level         byte
	Flags         byte
}

//struct PROTO_HOLY_PROMISE_DATE
//{
//  int _bf0;
//};
type HolyPromiseDate struct {
	BF0 int32
}
