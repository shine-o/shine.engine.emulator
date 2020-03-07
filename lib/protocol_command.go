package lib

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
)

type ProtocolCommandBase struct {
	packetType    string
	length        int
	department    uint16
	command       uint16
	operationCode uint16
	data          []byte
}

type PC struct {
	pcb ProtocolCommandBase
	pcc interface{} // protocol command concrete, eg: PROTO_NC_QUEST_GIVEUP_ACK
}

type NcMiscSeedAck struct{}

type NcUserClientVersionCheckReq struct{}

type NcUserClientRightversionCheckAck struct{}

type NcUserUsLoginReq struct{}

type NcUserXtrapReq struct{}

type NcUserXtrapAck struct{}

type NcUserLoginAck struct{}

type NcUserWorldStatusReq struct{}

type NcUserWorldStatusAck struct{}

// reassemble packet raw data
func (pcb *ProtocolCommandBase) RawData() []byte {
	var r []byte
	if pcb.packetType == "small" {
		r = append(r, uint8(pcb.length))
	} else {
		r = append(r, uint8(0))
		r = append(r, byte(pcb.length))
	}
	//r = append(r, byte(pcb.operationCode))
	r = append(r, pcb.data...)
	return r
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
