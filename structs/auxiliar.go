// auxiliar structs and functions for Network Commands
package structs

/* 3801 */
// union Name4
// {
//	char n4_name[16];
//	unsigned int n4_code[4];
// };
type Name4 struct {
	Name [16]byte `struct:"[16]byte"`
}

/* 3926 */
// struct __unaligned __declspec(align(2)) PROTO_NC_USER_LOGIN_ACK::WorldInfo
// {
//	char worldno;
//	Name4 worldname;
//	char worldstatus;
//};
type WorldInfo struct {
	WorldNumber byte `struct:"byte"`
	WorldName   Name4
	WorldStatus byte `struct:"byte"`
}

// union Name256Byte
// {
//	char n256_name[256];
//	unsigned __int64 n256_code[32];
// };
type Name256Byte struct {
	Name [256]byte `struct:"[256]byte"`
}

// /* 3256 */
// union Name5
// {
//  char n5_name[20];
//  unsigned int n5_code[5];
// };
type Name5 struct {
	Name [20]byte `struct:"[20]byte"`
}

// /* 2839 */
// union Name3
// {
//  char n3_name[12];
//  unsigned int n3_code[3];
// };
type Name3 struct {
	Name [12]byte `struct:"[12]byte"`
}

///* 2787 */
// union Name8
// {
//	char n8_name[32];
//	unsigned int n8_code[8];
// };
type Name8 struct {
	Name [32]byte `struct:"[32]byte"`
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
	EquHead         int16 `struct:"uint16"`
	EquMouth        int16 `struct:"uint16"`
	EquRightHand    int16 `struct:"uint16"`
	EquBody         int16 `struct:"uint16"`
	EquLeftHand     int16 `struct:"uint16"`
	EquPant         int16 `struct:"uint16"`
	EquBoot         int16 `struct:"uint16"`
	EquAccBoot      int16 `struct:"uint16"`
	EquAccPant      int16 `struct:"uint16"`
	EquAccHeadA     int16 `struct:"uint16"`
	EquMinimonR     int16 `struct:"uint16"`
	EquEye          int16 `struct:"uint16"`
	EquAccLeftHand  int16 `struct:"uint16"`
	EquAccRightHand int16 `struct:"uint16"`
	EquAccBack      int16 `struct:"uint16"`
	EquCosEff       int16 `struct:"uint16"`
	EquAccHip       int16 `struct:"uint16"`
	EquMinimon      int16 `struct:"uint16"`
	EquAccShield    int16 `struct:"uint16"`
}

// /* 2458 */
// struct SHINE_XY_TYPE
// {
//  unsigned int x;
//  unsigned int y;
// };
type ShineXYType struct {
	X uint32 `struct:"uint32"`
	Y uint32 `struct:"uint32"`
}

// /* 3808 */
//struct SHINE_DATETIME
//{
//  unsigned __int32 year : 4;
//  unsigned __int32 month : 4;
//  unsigned __int32 day : 5;
//  unsigned __int32 hour : 5;
//  unsigned __int32 min : 6;
//  unsigned __int32 sec : 6;
//};
type ShineDateTime struct {
	Year  uint32 `struct:"uint32"`
	Month uint32 `struct:"uint32"`
	Day   uint32 `struct:"uint32"`
	Hour  uint32 `struct:"uint32"`
	Min   uint32 `struct:"uint32"`
	Sec   uint32 `struct:"uint32"`
}

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
	TutorialState int  `struct:"uint32"`
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
	TsMax       = 3
)

// struct __unaligned __declspec(align(2)) PROTO_AVATARINFORMATION
// {
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
//  PROTO_TUTORIAL_INFO TutorialInfo;   nein nein x(
// };
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
