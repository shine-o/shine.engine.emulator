package structs

//union ITEM_INVEN
//{
//  unsigned __int16 Inven;
//};
type ItemInventory struct {
	Inventory uint16
}

//struct PROTO_CHARGED_ITEM_INFO
//{
//  unsigned int ItemOrderNo;
//  unsigned int ItemCode;
//  unsigned int ItemAmount;
//  ShineDateTime ItemRegisterDate;
//};
type ChargedItemInfo struct {
	ItemOrderNo      uint32
	ItemCode         uint32
	ItemAmount       uint32
	ItemRegisterDate ShineDateTime
}

//struct SHINE_ITEM_VAR_STRUCT
//{
//  unsigned __int16 itemid;
//  char itemattr[];
//};
type ShineItemVar struct {
	ItemID   uint16
	ItemAttr [19]byte // why 19?
}

//struct PROTO_ITEMPACKET_INFORM
//{
//  char datasize;
//  ITEM_INVEN location;
//  SHINE_ITEM_STRUCT info;
//};
type ItemPacketInfo struct {
	DataSize byte
	Location ItemInventory
	//struct SHINE_ITEM_STRUCT
	//{
	//  unsigned __int16 itemid;
	//  SHINE_ITEM_ATTRIBUTE itemattr;
	//};
	ItemID   uint16
	ItemAttr []byte `struct-size:"DataSize - 4"`
	//ShineItem struct {
	//
	//}
}
