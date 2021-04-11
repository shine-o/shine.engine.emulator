package structs

//struct ABSTATE_INFORMATION
type AbstateInformation struct {
	//enum ABSTATEINDEX
	AbstateIndex uint32
	RestKeepTime uint32
	Strength     uint32
}

//struct PROTO_NC_BRIEFINFO_DROPEDITEM_CMD::<unnamed-type-attr>
type NcBriefInfoDroppedItemCmdAttr struct {
	BF0 byte
}

//union PROTO_NC_BRIEFINFO_REGENMOB_CMD::<unnamed-type-flag>
type BriefInfoRegenMobCmdFlag struct {
	Data string `struct:"[112]byte"`
}

//struct SHINE_COORD_TYPE
type ShineCoordType struct {
	XY        ShineXYType
	Direction byte
}

//union PROTO_NC_BRIEFINFO_LOGINCHARACTER_CMD::<unnamed-type-shapedata>
type NcBriefInfoLoginCharacterCmdShapeData struct {
	Data [45]byte //
	//NotCamp CharBriefInfoNotCamped
}

//struct CHARBRIEFINFO_NOTCAMP
type CharBriefInfoNotCamp struct {
	Equip ProtoEquipment
}

//struct CHARBRIEFINFO_CAMP
type CharBriefInfoCamp struct {
	MiniHouse uint16
	Dummy     [10]byte //
}

//struct CHARBRIEFINFO_BOOTH
type CharBriefInfoBooth struct {
	Camp      CharBriefInfoCamp
	IsSelling byte
	SignBoard StreetBoothSignBoard
}

//struct CHARBRIEFINFO_RIDE::RideInfo
type CharBriefInfoRideInfo struct {
	Horse uint16
}

//struct CHARBRIEFINFO_RIDE
type CharBriefInfoRide struct {
	Equip    ProtoEquipment
	RideInfo CharBriefInfoRideInfo
}

//struct STOPEMOTICON_DESCRIPT
type StopEmoticonDescript struct {
	EmoticonID    byte
	EmoticonFrame uint16
}

//struct CHARTITLE_BRIEFINFO
type CharTitleBriefInfo struct {
	Type      byte
	ElementNo byte
	MobID     uint16
}

//struct ABNORMAL_STATE_BIT
type AbstateBit struct {
	Data [111]byte
}
