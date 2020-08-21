package structs

//struct PROTO_NC_BOOTH_ENTRY_REQ
//{
//  unsigned __int16 booth;
//};
type NcBoothEntryReq struct {
	Booth uint16
}

//struct PROTO_NC_BOOTH_SOMEONEOPEN_CMD
//{
//  unsigned __int16 handle;
//  CHARBRIEFINFO_CAMP tent;
//  char issell;
//  STREETBOOTH_SIGNBOARD signboard;
//};
type NcBoothSomeoneOpenCmd struct {
	Handle    uint16
	Tent      CharBriefInfoCamp
	IsSelling byte
	Sign      StreetBoothSignBoard
}

//struct PROTO_NC_BOOTH_REFRESH_REQ
//{
//  unsigned __int16 booth;
//};
type NcBoothRefreshReq struct {
	Booth uint16
}

//struct PROTO_NC_BOOTH_ENTRY_SELL_ACK
//{
//  unsigned __int16 err;
//  unsigned __int16 boothhandle;
//  char numofitem;
//  PROTO_NC_BOOTH_ENTRY_SELL_ACK::BoothItemList items[];
//};
type NcBoothEntrySellAck struct {
	Err         uint16
	BoothHandle uint16
	NumOfItems  byte
	Items       []NcBoothEntrySellAckItemList `struct:"sizefrom=NumOfItems"`
}

//struct PROTO_NC_BOOTH_SEARCH_BOOTH_CLOSED_CMD
//{
//  unsigned __int16 nClosedBoothOwnerHandle;
//};
type NcBoothSearchBoothClosedCmd struct {
	ClosedBoothOwnerHandle uint16
}
