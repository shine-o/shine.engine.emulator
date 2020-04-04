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
