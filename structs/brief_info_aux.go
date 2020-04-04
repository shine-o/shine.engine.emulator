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
	Strength uint32
}