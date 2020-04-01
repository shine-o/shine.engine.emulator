package structs

import (
	"encoding/json"
	"reflect"
)

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

func (nc * NcMiscSeedAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc * NcMiscSeedAck) PdbAnalog() string {
	return `
		struct PROTO_NC_MISC_SEED_ACK
		{
		  unsigned __int16 seed;
		};
`
}

func (nc * NcMiscSeedAck) Pack() ([]byte, error) {
	return Pack(nc)
}

func (nc * NcMiscSeedAck) Unpack(data []byte) error {
	return Unpack(data, nc)
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

func (nc * NcMiscGameTimeAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc * NcMiscGameTimeAck) PdbAnalog() string {
	return `
	struct PROTO_NC_MISC_GAMETIME_ACK
	{
	  char hour;
	  char minute;
	  char second;
	};
`
}

func (nc * NcMiscGameTimeAck) Pack() ([]byte, error) {
	return Pack(nc)
}

func (nc * NcMiscGameTimeAck) Unpack(data []byte) error {
	return Unpack(data, nc)
}