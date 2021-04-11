package structs

//union ITEM_INVEN
type ItemInventory struct {
	Inventory uint16
}

//struct PROTO_CHARGED_ITEM_INFO
type ChargedItemInfo struct {
	ItemOrderNo      uint32
	ItemCode         uint32
	ItemAmount       uint32
	ItemRegisterDate ShineDateTime
}

//struct SHINE_ITEM_VAR_STRUCT
type ShineItemVar struct {
	ItemID   uint16
	ItemAttr [19]byte // why 19?
}

//struct PROTO_ITEMPACKET_INFORM
type ItemPacketInfo struct {
	DataSize byte
	Location ItemInventory
	//struct SHINE_ITEM_STRUCT
	ItemID   uint16
	ItemAttr []byte `struct-size:"DataSize - 4"`
}
