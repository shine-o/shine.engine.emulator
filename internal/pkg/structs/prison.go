package structs

//struct PROTO_NC_PRISON_GET_ACK
//{
//  unsigned __int16 err;
//  unsigned __int16 nMinute;
//  char sReason[16];
//  char sRemark[64];
//};
type NcPrisonGetAck struct {
	Err    uint16
	Minute uint16

	// seems 2016 files do not send this data
	//Reason [16]byte
	//Remark [64]byte
}
