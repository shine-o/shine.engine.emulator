package structs

// struct PROTO_NC_AVATAR_CREATE_REQ
type NcAvatarCreateReq struct {
	SlotNum byte
	Name    Name5
	Shape   ProtoAvatarShapeInfo
}

// struct PROTO_NC_AVATAR_CREATESUCC_ACK
type NcAvatarCreateSuccAck struct {
	NumOfAvatar byte
	Avatar      AvatarInformation
}

// struct PROTO_NC_AVATAR_CREATEFAIL_ACK
type NcAvatarCreateFailAck struct {
	Err uint16
}

// struct PROTO_NC_AVATAR_ERASE_REQ
type NcAvatarEraseReq struct {
	Slot byte
}

// struct PROTO_NC_AVATAR_ERASESUCC_ACK
type NcAvatarEraseSuccAck struct {
	Slot byte
}
