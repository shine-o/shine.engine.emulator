package structs

import (
	"encoding/json"
	"reflect"
)

// struct PROTO_NC_AVATAR_CREATE_REQ
// {
//	char slotnum;
//	Name5 name;
//	PROTO_AVATAR_SHAPE_INFO char_shape;
// };
type NcAvatarCreateReq struct {
	SlotNum byte `struct:"byte"`
	Name    Name5
	Shape   ProtoAvatarShapeInfo
}

func (nc *NcAvatarCreateReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcAvatarCreateReq) PdbType() string {
	return `
		struct PROTO_NC_AVATAR_CREATE_REQ
		{
		  char slotnum;
		  Name5 name;
		  PROTO_AVATAR_SHAPE_INFO char_shape;
		};
`
}

// struct PROTO_NC_AVATAR_CREATESUCC_ACK
//{
//  char numofavatar;
//  PROTO_AVATARINFORMATION avatar;
//};
type NcAvatarCreateSuccAck struct {
	NumOfAvatar byte `struct:"byte"`
	Avatar      AvatarInformation
}

func (nc *NcAvatarCreateSuccAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcAvatarCreateSuccAck) PdbType() string {
	return `
		struct PROTO_NC_AVATAR_CREATESUCC_ACK
		{
		  char numofavatar;
		  PROTO_AVATARINFORMATION avatar;
		};
`
}

//struct PROTO_NC_AVATAR_CREATEFAIL_ACK
//{
//  unsigned __int16 err;
//};
type NcAvatarCreateFailAck struct {
	Err uint16 `struct:"uint16"`
}

func (nc *NcAvatarCreateFailAck) PdbType() string {
	return `
	struct PROTO_NC_AVATAR_CREATEFAIL_ACK
	{
	  unsigned __int16 err;
	};
`
}

func (nc *NcAvatarCreateFailAck) String() string {
	// todo: refactor to common func
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

//struct PROTO_NC_AVATAR_ERASE_REQ
//{
//char slot;
//};
type NcAvatarEraseReq struct {
	Slot byte `struct:"byte"`
}

func (nc *NcAvatarEraseReq) PdbType() string {
	return `
	struct PROTO_NC_AVATAR_ERASE_REQ
	{
	  char slot;
	};
`
}

//struct PROTO_NC_AVATAR_ERASESUCC_ACK
//{
//char slot;
//};
type NcAvatarEraseSuccAck struct {
	Slot byte `struct:"byte"`
}

func (nc *NcAvatarEraseSuccAck) PdbType() string {
	return `
	struct PROTO_NC_AVATAR_ERASESUCC_ACK
	{
	  char slot;
	};
`
}

func (nc *NcAvatarEraseSuccAck) String() string {
	// todo: refactor to common func
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}
