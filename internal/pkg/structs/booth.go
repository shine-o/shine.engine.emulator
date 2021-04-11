package structs

//struct PROTO_NC_BOOTH_ENTRY_REQ
type NcBoothEntryReq struct {
	Booth uint16
}

//struct PROTO_NC_BOOTH_SOMEONEOPEN_CMD
type NcBoothSomeoneOpenCmd struct {
	Handle    uint16
	Tent      CharBriefInfoCamp
	IsSelling byte
	Sign      StreetBoothSignBoard
}

//struct PROTO_NC_BOOTH_REFRESH_REQ
type NcBoothRefreshReq struct {
	Booth uint16
}

//struct PROTO_NC_BOOTH_ENTRY_SELL_ACK
type NcBoothEntrySellAck struct {
	Err         uint16
	BoothHandle uint16
	NumOfItems  byte
	Items       []NcBoothEntrySellAckItemList `struct:"sizefrom=NumOfItems"`
}

//struct PROTO_NC_BOOTH_SEARCH_BOOTH_CLOSED_CMD
type NcBoothSearchBoothClosedCmd struct {
	ClosedBoothOwnerHandle uint16
}
