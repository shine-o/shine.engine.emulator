package structs

//struct CHAR_ID_CHANGE_DATA
type CharIdChangeData struct {
	NeedChangeId byte
	Init         byte
	RowNo        uint32
}

//struct PROTO_TUTORIAL_INFO
type ProtoTutorialInfo struct {
	TutorialState int32
	TutorialStep  byte
}

//enum TUTORIAL_STATE
const (
	TsProgress  = 0
	TsDone      = 1
	TsSkip      = 2
	TsException = 3
	TsMax       = 4
)

//struct PROTO_AVATARINFORMATION
type AvatarInformation struct {
	ChrRegNum        uint32
	Name             Name5
	Level            uint16
	Slot             byte
	LoginMap         Name3
	DelInfo          ProtoAvatarDeleteInfo
	Shape            ProtoAvatarShapeInfo
	Equip            ProtoEquipment
	KqHandle         uint32
	KqMapName        Name3
	KqCoord          ShineXYType
	KqDate           ShineDateTime
	CharIdChangeData CharIdChangeData
	TutorialInfo     ProtoTutorialInfo
}

// struct PROTO_AVATAR_DELETE_INFO
type ProtoAvatarDeleteInfo struct {
	Year  byte
	Month byte
	Day   byte
	Hour  byte
	Min   byte
}

// struct PROTO_AVATAR_SHAPE_INFO
type ProtoAvatarShapeInfo struct {
	BF        byte
	HairType  byte
	HairColor byte
	FaceShape byte
}

//struct PROTO_EQUIPMENT::<unnamed-type-upgrade>
//  __int8 lefthand : 4;
//  __int8 righthand : 4;
//  __int8 body : 4;
//  __int8 leg : 4;
//  __int8 shoes : 4;
type EquipmentUpgrade struct {
	// first byte is the upgrade level
	Gap [2]uint8 `struct:"[2]uint8"`
	//Gap uint16
	BF2 byte
}

//struct PROTO_EQUIPMENT
type ProtoEquipment struct {
	EquHead         uint16
	EquMouth        uint16
	EquRightHand    uint16
	EquBody         uint16
	EquLeftHand     uint16
	EquPant         uint16
	EquBoot         uint16
	EquAccBoot      uint16
	EquAccPant      uint16
	EquAccBody      uint16
	EquAccHeadA     uint16
	EquMinimonR     uint16
	EquEye          uint16
	EquAccLeftHand  uint16
	EquAccRightHand uint16
	EquAccBack      uint16
	EquCosEff       uint16
	EquAccHip       uint16
	EquMinimon      uint16
	EquAccShield    uint16
	Upgrade         EquipmentUpgrade
}
