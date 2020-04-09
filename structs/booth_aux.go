package structs

//struct CHARBRIEFINFO_CAMP
//{
//  unsigned __int16 minihouse;
//  char dummy[10];
//};
type CharBriefInfoCamp struct {
	MiniHouse uint16
	Dummy     [10]byte
}

//struct STREETBOOTH_SIGNBOARD
//{
//  char signboard[30];
//};
type StreetBoothSignBoard struct {
	Text [30]byte
}

//struct PROTO_NC_BOOTH_ENTRY_SELL_ACK::BoothItemList
//{
//  char datasize;
//  char boothslot;
//  unsigned __int64 unitcost;
//  SHINE_ITEM_STRUCT item;
//};
type NcBoothEntrySellAckItemList struct {
	DataSize byte
	BoothSlot byte
	UnitCost uint64
	//struct SHINE_ITEM_STRUCT
	//{
	//  unsigned __int16 itemid;
	//  SHINE_ITEM_ATTRIBUTE itemattr;
	//};
	Info struct {
		ItemID   uint16
		ItemAttr []byte `struct-size:"Parent.DataSize - 2"`
	}
}