package structs

//struct ABSTATE_INFORMATION
//{
//  ABSTATEINDEX abstateID;
//  unsigned int restKeeptime;
//  unsigned int strength;
//};
type AbstateInformation struct {
	//enum ABSTATEINDEX
	//{
	//  STA_SEVERBONE = 0x0,
	//  STA_REDSLASH = 0x1,
	//  STA_BATTLEBLOWSTUN = 0x2,
	//  [ .... many more ]
	//  STA_MIGHTYSOULMAIN = 0x3,
	//  MAX_ABSTATEINDEX = 0x336,
	//};
	AbstateIndex uint32
	RestKeepTime uint32
	Strength     uint32
}

//struct PROTO_NC_BRIEFINFO_DROPEDITEM_CMD::<unnamed-type-attr>
//{
//  char _bf0;
//};
type NcBriefInfoDroppedItemCmdAttr struct {
	BF0 byte
}

//union PROTO_NC_BRIEFINFO_REGENMOB_CMD::<unnamed-type-flag>
//{
//  ABNORMAL_STATE_BIT abstatebit;   0
//  char gate2where[12];   			 1
//};
type BriefInfoRegenMobCmdFlag struct {
	Data [112]byte
}

//struct SHINE_COORD_TYPE
//{
//  SHINE_XY_TYPE xy;
//  char dir;
//};
type ShineCoordType struct {
	XY        ShineXYType
	Direction byte
}

//union PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD::<unnamed-type-shapedata>
//{
//  CHARBRIEFINFO_NOTCAMP notcamp;
//  CHARBRIEFINFO_CAMP camp;
//  CHARBRIEFINFO_BOOTH booth;
//  CHARBRIEFINFO_RIDE ride;
//};
type NcBriefInfoLoginCharacterCmdShapeData struct {
	Data [45]byte //
	//NotCamp CharBriefInfoNotCamped
}

//struct CHARBRIEFINFO_NOTCAMP
//{
//  PROTO_EQUIPMENT equip;
//};
type CharBriefInfoNotCamp struct {
	Equip ProtoEquipment
}

//struct CHARBRIEFINFO_CAMP
//{
//  unsigned __int16 minihouse;
//  char dummy[10];
//};
type CharBriefInfoCamp struct {
	MiniHouse uint16
	Dummy     [10]byte //
}

//struct CHARBRIEFINFO_BOOTH
//{
//  CHARBRIEFINFO_CAMP camp;
//  char issell;
//  STREETBOOTH_SIGNBOARD signboard;
//};
type CharBriefInfoBooth struct {
	Camp CharBriefInfoCamp
	IsSelling byte
	SignBoard StreetBoothSignBoard
}

//struct CHARBRIEFINFO_RIDE::RideInfo
//{
//  unsigned __int16 horse;
//};
type CharBriefInfoRideInfo struct {
	Horse uint16
}

//struct CHARBRIEFINFO_RIDE
//{
//  PROTO_EQUIPMENT equip;
//  CHARBRIEFINFO_RIDE::RideInfo rideinfo;
//};
type CharBriefInfoRide struct {
	Equip ProtoEquipment
	RideInfo CharBriefInfoRideInfo
}

//struct STOPEMOTICON_DESCRIPT
//{
//  char emoticonid;
//  unsigned __int16 emoticonframe;
//};
type StopEmoticonDescript struct {
	EmoticonID    byte
	EmoticonFrame uint16
}

//struct CHARTITLE_BRIEFINFO
//{
//  char Type;
//  char ElementNo;
//  unsigned __int16 MobID;
//};
type CharTitleBriefInfo struct {
	Type      byte
	ElementNo byte
	MobID     uint16
}

//struct ABNORMAL_STATE_BIT
//{
//  #char statebit[103]; 2017
//  #char statebit[99];  2016
//  char statebit[112]; 2020
//};
type AbstateBit struct {
	Data [111]byte
}
