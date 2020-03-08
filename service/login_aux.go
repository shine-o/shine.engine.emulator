// auxiliar structs and functions for Network Commands
package service

import protocol "shine.engine.packet-protocol"

//import protocol "github.com/shine-o/shine.engine.protocol"

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
// seems like a utility map struct for names, maybe related with NIF files
// union Name4
// {
//	char n4_name[16];
//	unsigned int n4_code[4];
// };
type ComplexName struct {
	Name     [16]byte
	NameCode [4]uint16
}

/* 3256 */
// union Name5
// {
//	char n5_name[20];
//	unsigned int n5_code[5];
// };
type ComplexName1 struct {
	Name     [20]byte
	NameCode [5]uint32
}

func logOutboundPacket(pc *protocol.Command) {
	log.Infof("Outbound packet %v", pc.Base.String())
}
