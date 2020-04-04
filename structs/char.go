package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_CHAR_CLIENT_SKILL_CMD
//{
//	char restempow;
//	PARTMARK PartMark;
//	unsigned __int16 nMaxNum;
//	PROTO_NC_CHAR_SKILLCLIENT_CMD skill;
//};
type NcCharClientSkillCmd struct {
	RestEmpower byte
	PartMark    PartMark
	MaxNum      uint16
	Skills      NcCharSkillClientCmd
}

func (nc *NcCharClientSkillCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharClientSkillCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_CLIENT_SKILL_CMD
	{
	  char restempow;
	  PARTMARK PartMark;
	  unsigned __int16 nMaxNum;
	  PROTO_NC_CHAR_SKILLCLIENT_CMD skill;
	};
`
}


//struct PROTO_NC_CHAR_SKILLCLIENT_CMD
//{
//	unsigned int chrregnum;
//	unsigned __int16 number;
//	PROTO_SKILLREADBLOCKCLIENT skill[];
//};
type NcCharSkillClientCmd struct {
	ChrRegNum uint32
	Number    uint16
	Skills    []SkillReadBlockClient `struct:"sizefrom=Number"`
}

func (nc *NcCharSkillClientCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharSkillClientCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_SKILLCLIENT_CMD
	{
	  unsigned int chrregnum;
	  unsigned __int16 number;
	  PROTO_SKILLREADBLOCKCLIENT skill[];
	};
`
}


//struct PROTO_NC_CHAR_CLIENT_ITEM_CMD
//{
//	char numofitem;
//	char box;
//	PROTO_NC_CHAR_CLIENT_ITEM_CMD::<unnamed-type-flag> flag;
//	PROTO_ITEMPACKET_INFORM ItemArray[];
//};
type NcCharClientItemCmd struct {
	NumOfItem byte `struct:"byte"`
	Box       byte `struct:"byte"`
	Flag      ProtoNcCharClientItemCmdFlag
	Items     []ProtoItemPacketInformation `struct:"sizefrom=NumOfItem"`
}

func (nc *NcCharClientItemCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcCharClientItemCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_CLIENT_ITEM_CMD
	{
	  char numofitem;
	  char box;
	  PROTO_NC_CHAR_CLIENT_ITEM_CMD::<unnamed-type-flag> flag;
	  PROTO_ITEMPACKET_INFORM ItemArray[];
	};
`
}


//struct PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD
//{
//  char CurrentTitle;
//  char CurrentTitleElement;
//  unsigned __int16 CurrentTitleMobID;
//  unsigned __int16 NumOfTitle;
//  CT_INFO TitleArray[];
//};
type NcClientCharTitleCmd struct {
	CurrentTitle        byte
	CurrentTitleElement byte
	CurrentTitleMobID   uint16
	NumOfTitle          uint16
	Titles              []CharTitleInfo `struct:"sizefrom=NumOfTitle"`
}

func (nc *NcClientCharTitleCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcClientCharTitleCmd) PdbType() string {
	return `
	struct PROTO_NC_CHAR_CLIENT_CHARTITLE_CMD
	{
	  char CurrentTitle;
	  char CurrentTitleElement;
	  unsigned __int16 CurrentTitleMobID;
	  unsigned __int16 NumOfTitle;
	  CT_INFO TitleArray[];
	};
`
}

