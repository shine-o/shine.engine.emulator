package lib

// auxiliar structs for Network Commands

/* 3926 */
// struct __unaligned __declspec(align(2)) PROTO_NC_USER_LOGIN_ACK::WorldInfo
// {
//	char worldno;
//	Name4 worldname;
//	char worldstatus;
//};
type WorldInfo struct {
	WorldNumber byte
	WorldName   ComplexName
	WorldStatus byte
}

/* 3801 */
// seems like a utility map struct for names
// union Name4
// {
//	char n4_name[16];
//	unsigned int n4_code[4];
// };
type ComplexName struct {
	Name     [16]byte
	NameCode [4]uint16
}
