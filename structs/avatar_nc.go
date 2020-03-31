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

//struct PROTO_NC_AVATAR_ERASE_REQ
//{
//char slot;
//};
type NcAvatarEraseReq struct {
	Slot byte `struct:"byte"`
}

func (nc * NcAvatarEraseReq) PdbAnalog() string {
	return `
	struct PROTO_NC_AVATAR_ERASE_REQ
	{
	  char slot;
	};
`
}

func (nc * NcAvatarEraseReq) Pack() ([]byte, error) {
	return Pack(nc)
}

func (nc * NcAvatarEraseReq) Unpack(data []byte) error {
	return Unpack(data, nc)
}

//struct PROTO_NC_AVATAR_ERASESUCC_ACK
//{
//char slot;
//};
type NcAvatarEraseSuccAck struct {
	Slot byte `struct:"byte"`
}

func (nc * NcAvatarEraseSuccAck) PdbAnalog() string {
	return `
	struct PROTO_NC_AVATAR_ERASESUCC_ACK
	{
	  char slot;
	};
`
}

func (nc * NcAvatarEraseSuccAck) String() string {
	// todo: refactor to common func
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc * NcAvatarEraseSuccAck) Pack() ([]byte, error) {
	return Pack(nc)
}

func (nc * NcAvatarEraseSuccAck) Unpack(data []byte) error {
	return Unpack(data, nc)
}