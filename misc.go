package networking

// RE client struct:
// struct PROTO_NC_MISC_SEED_ACK
// {
//	unsigned __int16 seed;
// };
// xorKey offset used by client to encrypt data
// same offset is used on the server side to decrypt data sent by the client
type ncMiscSeedAck struct {
	seed uint16
}
