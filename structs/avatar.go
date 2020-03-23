package structs

// struct PROTO_NC_AVATAR_CREATE_REQ
// {
//	char slotnum;
//	Name5 name;
//	PROTO_AVATAR_SHAPE_INFO char_shape;
// };
type NcAvatarCreateReq struct {
	SlotNum   byte `struct:"byte"`
	Name      Name5
	CharShape ProtoAvatarShapeInfo
}

// struct PROTO_NC_AVATAR_CREATEDATASUC_ACK
// {
//	char numofavatar;
//	PROTO_AVATARINFORMATION avatar[];
//};
