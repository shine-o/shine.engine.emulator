package structs

// RE client struct:
// struct PROTO_NC_MISC_SEED_ACK
// {
//	unsigned __int16 seed;
// };
// xorKey offset used by client to encrypt data
// same offset is used on the server side to decrypt data sent by the client
type NcMiscSeedAck struct {
	Seed uint16 `struct:"uint16"`
}

// struct PROTO_NC_MISC_GAMETIME_ACK
// {
//	char hour;
//	char minute;
//	char second;
// };
type NcMiscGameTimeAck struct {
	Hour   byte `struct:"byte"`
	Minute byte `struct:"byte"`
	Second byte `struct:"byte"`
}
