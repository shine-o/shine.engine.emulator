package structs

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

// struct PROTO_NC_AVATAR_CREATESUCC_ACK
//{
//  char numofavatar;
//  PROTO_AVATARINFORMATION avatar;
//};
type NcAvatarCreateSuccAck struct {
	NumOfAvatar byte `struct:"byte"`
	Avatar      AvatarInformation
}

//struct PROTO_NC_AVATAR_CREATEFAIL_ACK
//{
//  unsigned __int16 err;
//};
type NcAvatarCreateFailAck struct {
	Err uint16 `struct:"uint16"`
}

//struct PROTO_NC_AVATAR_ERASE_REQ
//{
//char slot;
//};
type NcAvatarEraseReq struct {
	Slot byte `struct:"byte"`
}

//struct PROTO_NC_AVATAR_ERASESUCC_ACK
//{
//char slot;
//};
type NcAvatarEraseSuccAck struct {
	Slot byte `struct:"byte"`
}