package structs

// struct PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ
type NcCharOptionSetShortcutDataReq struct {
	Count     byte
	Shortcuts []ShortCutData `struct:"sizefrom=Count"`
}

// struct PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_ACK
type NcCharOptionImproveShortcutDataAck struct {
	ErrCode uint16
}
