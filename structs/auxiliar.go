// auxiliar structs and functions for Network Commands
package structs

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

// struct SHINE_XY_TYPE
// {
//  unsigned int x;
//  unsigned int y;
// };
type ShineXYType struct {
	X uint32 `struct:"uint32"`
	Y uint32 `struct:"uint32"`
}

// struct SHINE_DATETIME
// {
//  int _bf0;
// };
type ShineDateTime struct {
	BF0 int32 `struct:"int32"`
}
