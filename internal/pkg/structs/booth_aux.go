package structs

//struct STREETBOOTH_SIGNBOARD
type StreetBoothSignBoard struct {
	//Text [30]byte
	Text string `struct:"[30]byte"`
}

//struct PROTO_NC_BOOTH_ENTRY_SELL_ACK::BoothItemList
type NcBoothEntrySellAckItemList struct {
	DataSize  byte
	BoothSlot byte
	UnitCost  uint64
	//struct SHINE_ITEM_STRUCT
	ItemID   uint16
	ItemAttr []byte `struct-size:"DataSize - 11"`
}
