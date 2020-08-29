package structs

// struct PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ
//{
//  char nShortCutDataCnt;
//  SHORT_CUT_DATA ShortCutData[];
//};
type NcCharOptionSetShortcutDataReq struct {
	Count     byte
	Shortcuts []ShortCutData `struct:"sizefrom=Count"`
}

// struct PROTO_NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_ACK
//{
//  unsigned __int16 nError;
//};
type NcCharOptionImproveShortcutDataAck struct {
	ErrCode uint16
}
