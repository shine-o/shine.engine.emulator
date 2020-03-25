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
	SlotNum   byte `struct:"byte"`
	Name      Name5
	Shape ProtoAvatarShapeInfo
}

func (nc * NcAvatarCreateReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc * NcAvatarCreateReq) PdbAnalog() string {
	return `
		struct PROTO_NC_AVATAR_CREATE_REQ
		{
		  char slotnum;
		  Name5 name;
		  PROTO_AVATAR_SHAPE_INFO char_shape;
		};
`
}

func (nc * NcAvatarCreateReq) Pack() ([]byte, error) {
	return Pack(nc)
}

func (nc * NcAvatarCreateReq) Unpack(data []byte) error {
	return Unpack(data, nc)
}

// struct PROTO_NC_AVATAR_CREATESUCC_ACK
//{
//  char numofavatar;
//  PROTO_AVATARINFORMATION avatar;
//};
type NcAvatarCreateSuccAck struct {
	NumOfAvatar byte `struct:"byte"`
	Avatar AvatarInformation
}

func (nc * NcAvatarCreateSuccAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc * NcAvatarCreateSuccAck) PdbAnalog() string {
	return `
		struct PROTO_NC_AVATAR_CREATESUCC_ACK
		{
		  char numofavatar;
		  PROTO_AVATARINFORMATION avatar;
		};
`
}

func (nc * NcAvatarCreateSuccAck) Pack() ([]byte, error) {
	return Pack(nc)
}

func (nc * NcAvatarCreateSuccAck) Unpack(data []byte) error {
	return Unpack(data, nc)
}