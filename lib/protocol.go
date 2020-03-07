package lib

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type ProtocolCommand struct {
	pcb ProtocolCommandBase
	pcc interface{} // protocol command concrete, eg: PROTO_NC_QUEST_GIVEUP_ACK
}

type ProtocolCommandBase struct {
	packetType    string
	length        int
	department    uint16
	command       uint16
	operationCode uint16
	data          []byte
}

// reassemble packet raw data
func (pcb *ProtocolCommandBase) RawData() []byte {
	var header []byte
	var data []byte

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, pcb.operationCode); err != nil {
		log.Fatalf("failed writing operation code to buffer %v", err)
	}

	data = append(data, buf.Bytes()...)
	data = append(data, pcb.data...)

	if len(data) > 255 { // means big packet
		header = append(header, byte(0))
		lenBuf := new(bytes.Buffer)
		if err := binary.Write(lenBuf, binary.LittleEndian, uint16(buf.Len())); err != nil {
			log.Fatalf("failed writing length for big packet to buffer %v", err)
		}
		header = append(header, buf.Bytes()...)
	} else {
		header = append(header, byte(len(data)))
	}

	return append(header, data...)
}

func (pcb *ProtocolCommandBase) String() string {
	type exportedPcb struct {
		PacketType    string `json:"packetType"`
		Length        int    `json:"length"`
		Department    uint16 `json:"department"`
		Command       string `json:"command"`
		OperationCode uint16 `json:"opCode"`
		Data          string `json:"data"`
		RawData       string `json:"rawData"`
		FriendlyName  string `json:"friendlyName"`
	}

	ePcb := exportedPcb{
		PacketType:    pcb.packetType,
		Length:        pcb.length,
		Department:    pcb.department,
		Command:       fmt.Sprintf("%X", pcb.command),
		OperationCode: pcb.operationCode,
		Data:          hex.EncodeToString(pcb.data),
		RawData:       hex.EncodeToString(pcb.RawData()),
	}

	if rawJson, err := json.Marshal(&ePcb); err != nil {
		log.Error(err)
		return ""
	} else {
		return string(rawJson)
	}
}

// RE client struct:
// struct PROTO_NC_MISC_SEED_ACK
// {
//	unsigned __int16 seed;
// };
//
// xorKey offset used by client to encrypt data
// same offset is used on the server side to decrypt data sent by the client
type ncMiscSeedAck struct {
	seed uint16
}

// RE client struct:
// struct PROTO_NC_USER_CLIENT_VERSION_CHECK_REQ
// {
//  char sVersionKey[64];
// };
type ncUserClientVersionCheckReq struct {
	VersionKey [64]byte
}

// RE client struct:
//struct __cppobj PROTO_NC_USER_CLIENT_WRONGVERSION_CHECK_ACK
//{
//};
type ncUserClientRightversionCheckAck struct{}

// RE client struct:
// struct PROTO_NC_USER_US_LOGIN_REQ
// {
//  char sUserName[260];
//  char sPassword[36];
//  Name5 spawnapps;
// };
type ncUserUsLoginReq struct{}

// RE client struct:
// struct PROTO_NC_USER_XTRAP_REQ
// {
//  char XTrapClientKeyLength;
//  char XTrapClientKey[];
// };
type ncUserXtrapReq struct{}

// RE client struct:
// struct PROTO_NC_USER_XTRAP_ACK
// {
//  char bSuccess;
// };
type ncUserXtrapAck struct{}

// RE client struct:
// struct __unaligned __declspec(align(1)) PROTO_NC_USER_LOGIN_ACK
// {
//  char numofworld;
//  PROTO_NC_USER_LOGIN_ACK::WorldInfo worldinfo[];
// };
type ncUserLoginAck struct {
	NumOfWorld byte
	Worlds     [1]WorldInfo
}

// RE client struct:
// struct __cppobj PROTO_NC_USER_WORLD_STATUS_REQ
// {
// };
type ncUserWorldStatusReq struct{}

// OPERATION CODE ONLY 3100
type ncUserWorldStatusAck struct{}
