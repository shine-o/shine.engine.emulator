package structs



// /* 3809 */
// struct __unaligned __declspec(align(2)) CHAR_ID_CHANGE_DATA
// {
//  char bNeedChangeID;
//  char bInit;
//  unsigned int nRowNo;
// };
type CharIdChangeData struct {
	NeedChangeId byte   `struct:"byte"`
	Init         byte   `struct:"byte"`
	RowNo        uint32 `struct:"uint32"`
}

///* 3810 */
// struct __unaligned __declspec(align(1)) PROTO_TUTORIAL_INFO
// {
//	TUTORIAL_STATE nTutorialState;
//	char nTutorialStep;
// };
type ProtoTutorialInfo struct {
	TutorialState int  `struct:"int32"`
	TutorialStep  byte `struct:"byte"`
}

///* 1387 */
//enum TUTORIAL_STATE
//{
//TS_PROGRESS = 0x0,
//TS_DONE = 0x1,
//TS_SKIP = 0x2,
//TS_EXCEPTION = 0x3,
//TS_MAX = 0x4,
//};
const (
	TsProgress  = 0
	TsDone      = 1
	TsSkip      = 2
	TsException = 3
	TsMax       = 4
)

// struct PROTO_AVATARINFORMATION
//{
//  unsigned int chrregnum;
//  Name5 name;
//  unsigned __int16 level;
//  char slot;
//  Name3 loginmap;
//  PROTO_AVATAR_DELETE_INFO delinfo;
//  PROTO_AVATAR_SHAPE_INFO shape;
//  PROTO_EQUIPMENT equip;
//  unsigned int nKQHandle;
//  Name3 sKQMapName;
//  SHINE_XY_TYPE nKQCoord;
//  SHINE_DATETIME dKQDate;
//  CHAR_ID_CHANGE_DATA CharIDChangeData;
//  PROTO_TUTORIAL_INFO TutorialInfo;
//};
type AvatarInformation struct {
	ChrRegNum        uint32 `struct:"uint32"`
	Name             Name5
	Level            uint16 `struct:"uint16"`
	Slot             byte   `struct:"byte"`
	LoginMap         Name3
	DelInfo          ProtoAvatarDeleteInfo
	Shape            ProtoAvatarShapeInfo
	Equip            ProtoEquipment
	KqHandle         uint32 `struct:"uint32"`
	KqMapName        Name3
	KqCoord          ShineXYType
	KqDate           ShineDateTime
	CharIdChangeData CharIdChangeData
	TutorialInfo     ProtoTutorialInfo
}

// /* 3807 */
// struct PROTO_AVATAR_DELETE_INFO
// {
//  char year;
//  char month;
//  char day;
//  char hour;
//  char min;
// };
type ProtoAvatarDeleteInfo struct {
	Year  byte `struct:"byte"`
	Month byte `struct:"byte"`
	Day   byte `struct:"byte"`
	Hour  byte `struct:"byte"`
	Min   byte `struct:"byte"`
}

// struct PROTO_AVATAR_SHAPE_INFO
//{
//char _bf0;
//char hairtype;
//char haircolor;
//char faceshape;
//};
type ProtoAvatarShapeInfo struct {
	BF        byte `struct:"byte"`
	HairType  byte `struct:"byte"`
	HairColor byte `struct:"byte"`
	FaceShape byte `struct:"byte"`
}

// /* 2317 */
// struct __unaligned __declspec(align(1)) PROTO_EQUIPMENT
// {
//  unsigned __int16 Equ_Head;
//  unsigned __int16 Equ_Mouth;
//  unsigned __int16 Equ_RightHand;
//  unsigned __int16 Equ_Body;
//  unsigned __int16 Equ_LeftHand;
//  unsigned __int16 Equ_Pant;
//  unsigned __int16 Equ_Boot;
//  unsigned __int16 Equ_AccBoot;
//  unsigned __int16 Equ_AccPant;
//  unsigned __int16 Equ_AccBody;
//  unsigned __int16 Equ_AccHeadA;
//  unsigned __int16 Equ_MiniMon_R;
//  unsigned __int16 Equ_Eye;
//  unsigned __int16 Equ_AccLeftHand;
//  unsigned __int16 Equ_AccRightHand;
//  unsigned __int16 Equ_AccBack;
//  unsigned __int16 Equ_CosEff;
//  unsigned __int16 Equ_AccHip;
//  unsigned __int16 Equ_Minimon;
//  unsigned __int16 Equ_AccShield;
//  $050E0EECA9116B4E3A3935292D917DD5 upgrade; ???
// };
type ProtoEquipment struct {
	EquHead         int16 `struct:"int16"`
	EquMouth        int16 `struct:"int16"`
	EquRightHand    int16 `struct:"int16"`
	EquBody         int16 `struct:"int16"`
	EquLeftHand     int16 `struct:"int16"`
	EquPant         int16 `struct:"int16"`
	EquBoot         int16 `struct:"int16"`
	EquAccBoot      int16 `struct:"int16"`
	EquAccPant      int16 `struct:"int16"`
	EquAccHeadA     int16 `struct:"int16"`
	EquMinimonR     int16 `struct:"int16"`
	EquEye          int16 `struct:"int16"`
	EquAccLeftHand  int16 `struct:"int16"`
	EquAccRightHand int16 `struct:"int16"`
	EquAccBack      int16 `struct:"int16"`
	EquCosEff       int16 `struct:"int16"`
	EquAccHip       int16 `struct:"int16"`
	EquMinimon      int16 `struct:"int16"`
	EquAccShield    int16 `struct:"int16"`
}
