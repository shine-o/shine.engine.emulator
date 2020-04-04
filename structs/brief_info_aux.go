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

//struct PROTO_NC_BRIEFINFO_REGENMOB_CMD
//{
//  unsigned __int16 handle;
//  char mode;
//  unsigned __int16 mobid;
//  SHINE_COORD_TYPE coord;
//  char flagstate;
//  PROTO_NC_BRIEFINFO_REGENMOB_CMD::<unnamed-type-flag> flag;
//  char sAnimation[32];
//  char nAnimationLevel;
//  char nKQTeamType;
//  char bRegenAni;
//};
type NcBriefInfoRegenMobCmd struct {
	Handle uint16
	Mode byte
	MobID uint16
	Coord ShineCoordType
	// 0,1
	FlagState byte
	//FlagZero RegenMobFlagZero `struct-if:"FlagState == 0"`
	FlagOne  RegenMobFlagOne  `struct-if:"FlagState == 1"`
	
	Animation [32]byte
	AnimationLevel byte
	KQTeamType byte
	RegenAni byte
}

type RegenMobFlagZero []byte

type RegenMobFlagOne []byte

func (rmfz * RegenMobFlagZero) SizeOf() int {
	return 103
}

func (rmfo * RegenMobFlagOne) SizeOf() int {
	return 12
}


//union PROTO_NC_BRIEFINFO_REGENMOB_CMD::<unnamed-type-flag>
//{
//  ABNORMAL_STATE_BIT abstatebit;   0
//  char gate2where[12];   			 1
//};
type BriefInfoRegenMobCmdFlag struct {
	AbstateBit AbnormalStateBit
}

//struct ABNORMAL_STATE_BIT
//{
//  char statebit[103];
//};
type AbnormalStateBit struct {
	StateBit [103]byte
}

//struct SHINE_COORD_TYPE
//{
//  SHINE_XY_TYPE xy;
//  char dir;
//};
type ShineCoordType struct {
	XY ShineXYType
	Direction byte
}