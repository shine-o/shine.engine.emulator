package data

//enum SHN_DATA_FILE_INDEX
//{
//  SHN_Abstate = 0x0,
//  SHN_ActiveSkill = 0x1,
//  SHN_CharacterTitleData = 0x2,
//  SHN_ChargedEffect = 0x3,
//  SHN_ClassName = 0x4,
//  SHN_Gather = 0x5,
//  SHN_GradeItemOption = 0x6,
//  SHN_ItemDismantle = 0x7,
//  SHN_ItemInfo = 0x8,
//  SHN_MapInfo = 0x9,
//  SHN_MiniHouse = 0xA,
//  SHN_MiniHouseFurniture = 0xB,
//  SHN_MiniHouseObjAni = 0xC,
//  SHN_MobInfo = 0xD,
//  SHN_PassiveSkill = 0xE,
//  SHN_Riding = 0xF,
//  SHN_SubAbstate = 0x10,
//  SHN_UpgradeInfo = 0x11,
//  SHN_WeaponAttrib = 0x12,
//  SHN_WeaponTitleData = 0x13,
//  SHN_MiniHouseFurnitureObjEffect = 0x14,
//  SHN_MiniHouseFurnitureEndure = 0x15,
//  SHN_DiceDividind = 0x16,
//  SHN_ActionViewInfo = 0x17,
//  SHN_MapLinkPoint = 0x18,
//  SHN_MapWayPoint = 0x19,
//  SHN_AbStateView = 0x1A,
//  SHN_ActiveSkillView = 0x1B,
//  SHN_CharacterTitleStateView = 0x1C,
//  SHN_EffectViewInfo = 0x1D,
//  SHN_ItemShopView = 0x1E,
//  SHN_ItemViewInfo = 0x1F,
//  SHN_MapViewInfo = 0x20,
//  SHN_MobViewInfo = 0x21,
//  SHN_NPCViewInfo = 0x22,
//  SHN_PassiveSkillView = 0x23,
//  SHN_ProduceView = 0x24,
//  SHN_CollectCardView = 0x25,
//  SHN_GTIView = 0x26,
//  SHN_ItemViewEquipTypeInfo = 0x27,
//  SHN_SingleData = 0x28,
//  SHN_MarketSearchInfo = 0x29,
//  SHN_ItemMoney = 0x2A,
//  SHN_PupMain = 0x2B,
//  SHN_ChatColor = 0x2C,
//  SHN_TermExtendMatch = 0x2D,
//  SHN_MinimonInfo = 0x2E,
//  SHN_MinimonAutoUseItem = 0x2F,
//  SHN_ChargedDeletableBuff = 0x30,
//  SHN_SlanderFilter = 0x31,
//  SHN_MaxCnt = 0x32,
//};

//enum CDataReader::TYPE_LIST
//{
//  TYPE_LIST_END = 0x0,
//  TYPE_LIST_BYTE = 0x1,
//  TYPE_LIST_WORD = 0x2,
//  TYPE_LIST_DWORD = 0x3,
//  TYPE_LIST_QWORD = 0x4,
//  TYPE_LIST_FLOAT = 0x5,
//  TYPE_LIST_FILENAME = 0x6,
//  TYPE_LIST_FILEAUTO = 0x7,
//  TYPE_LIST_REMARK = 0x8,
//  TYPE_LIST_STR = 0x9,
//  TYPE_LIST_STRAUTO = 0xA,
//  TYPE_LIST_INX = 0xB,
//  TYPE_LIST_INXBYTE = 0xC,
//  TYPE_LIST_INXWORD = 0xD,
//  TYPE_LIST_INXDWORD = 0xE,
//  TYPE_LIST_INXQWORD = 0xF,
//  TYPE_LIST_BYTE_BIT = 0x10,
//  TYPE_LIST_WORD_BIT = 0x11,
//  TYPE_LIST_DWORD_BIT = 0x12,
//  TYPE_LIST_QWORD_BIT = 0x13,
//  TYPE_LIST_BYTE_ARRAY = 0x14,
//  TYPE_LIST_WORD_ARRAY = 0x15,
//  TYPE_LIST_DWORD_ARRAY = 0x16,
//  TYPE_LIST_QWORD_ARRAY = 0x17,
//  TYPE_LIST_STR_ARRAY = 0x18,
//  TYPE_LIST_STRAUTO_ARRAY = 0x19,
//  TYPE_LIST_VARSTR = 0x1A,
//  TYPE_LIST_INXSTR = 0x1B,
//  TYPE_LIST_UNKNOWNED = 0x1C,
//  TYPE_LIST_TWO_INX = 0x1D,
//};
type ShineDataType uint32

const (
	TypeListEnd ShineDataType = iota
	TypeListByte
	TypeListWord
	TypeListDword
	TypeListQWord
	TypeListFloat
	TypeListFilename
	TypeListFileAuto
	TypeListRemark
	TypeListStr
	TypeListStrAuto
	TypeListInx
	TypeListInxByte
	TypeListInxWord
	TypeListInxDword
	TypeListInxQWord
	TypeListByteBit
	TypeListWordBit
	TypeListDwordBit
	TypeListQWordBit
	TypeListByteArray
	TypeListWordArray
	TypeListDWordArray
	TypeListQWordArray
	TypeListStrArray
	TypeListStrAutoArray
	TypeListVarStr
	TypeListInxStr
	TypeListUnknownEd
	TypeListTwoInx
)

//enum CDataReader::DATA_MODE
//{
//  DATA_MODE_NORMAL = 0x0,
//  DATA_MODE_ENCRYPTION = 0x1,
//};
type ShineDataMode uint32

const (
	DataModeNormal ShineDataMode = iota
	DataModeEncryption
)

//struct CDataReader::HEAD
//{
//  unsigned int nVersionKey;
//  char sVersion[20];
//  unsigned int nReserved;
//  CDataReader::DATA_MODE nDataMode;
//  unsigned int nFileSize;
//  unsigned int nDataSize;
//  unsigned int nNumOfRecord;
//  unsigned int nFieldSize;
//  unsigned int nNumOfField;
//  CDataReader::FIELD Field[];
//};
type ShineRawFile struct {
	VersionKey uint32
	Version    [20]byte
	Reserved   uint32
	DataMode   ShineDataMode
	FileSize   int32
	// column, row data
	Data []byte `struct-while:"!_eof"`
}

//struct CDataReader::FIELD
//{
//  char Name[48];
//  CDataReader::TYPE_LIST Type;
//  unsigned int Size;
//};
type ShineColumn struct {
	Name string `struct:"[48]byte"`
	Type ShineDataType
	Size uint32
}
