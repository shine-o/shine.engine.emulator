package structs

//union ITEM_INVEN
type ItemInventory struct {
	Inventory uint16
}

//struct SHINE_ITEM_VAR_STRUCT
type ShineItemVar struct {
	ItemID   uint16
	ItemAttr []byte
}

//struct PROTO_ITEMPACKET_INFORM
type ItemPacketInfo struct {
	DataSize byte
	Location ItemInventory
	ItemID   uint16
	ItemAttr []byte `struct-size:"DataSize - 4"`
}
