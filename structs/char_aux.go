package structs

//struct PROTO_SKILLREADBLOCKCLIENT
//{
//	unsigned __int16 skillid;
//	unsigned int cooltime;
//	PROTO_SKILLREADBLOCKCLIENT::<unnamed-type-empow> empow;
//	unsigned int mastery;
//};
type SkillReadBlockClient struct {
	SkillID  uint16
	CoolTime uint32
	Empower  SkillReadBlockClientEmpower
	Mastery  uint32
}

//struct PROTO_SKILLREADBLOCKCLIENT::<unnamed-type-empow>
//{
//  _BYTE gap0[1];
//  char _bf1;
//};
type SkillReadBlockClientEmpower struct {
	Gap0 byte
	BF1  byte
}

//struct PARTMARK
//{
//	char _bf0;
//};
type PartMark struct {
	BF0 byte
}

//struct PROTO_NC_CHAR_CLIENT_ITEM_CMD::<unnamed-type-flag>
//{
//	char _bf0;
//};
type ProtoNcCharClientItemCmdFlag struct {
	BF0 byte
}

//struct PROTO_ITEMPACKET_INFORM
//{
//	char datasize;
//	ITEM_INVEN location;
//	SHINE_ITEM_STRUCT info;
//};
type ProtoItemPacketInformation struct {
	DataSize byte
	// can't be done like this, since data size also covers Location and Info and there's no way to use sizefrom with operators -+ :(
	// at the handler level, i would have to read the fields manually.
	ItemData []byte `struct:"sizefrom=DataSize"`
}

//struct CT_INFO
//{
//  char Type;
//  char _bf1;
//};
type CharTitleInfo struct {
	Type byte
	BF1  byte
}
