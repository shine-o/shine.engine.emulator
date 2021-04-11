package structs

//struct PROTO_HOLY_PROMISE_INFO
type HolyPromiseInfo struct {
	PromiseCharID Name5
	LogonInfo     HolyPromiseDate
	Level         byte
	Flags         byte
}

//struct PROTO_HOLY_PROMISE_DATE
type HolyPromiseDate struct {
	BF0 int32
}
