package structs

//struct PROTO_NC_MAP_LOGIN_REQ
//{
//  PROTO_NC_CHAR_ZONE_CHARDATA_REQ chardata;
//  Name8 checksum[50];
//};
type NcMapLoginReq struct {
	CharData NcZoneCharDataReq
	CheckSum [56]Name8
}

//struct PROTO_NC_CHAR_ZONE_CHARDATA_REQ
//{
//  unsigned __int16 wldmanhandle;
//  Name5 charid;
//};
type NcZoneCharDataReq struct {
	WorldManager uint16
	CharID       Name5
}

//struct PROTO_NC_MAP_LOGINCOMPLETE_CMD
//{
//  char dummy[1];
//};
type NcMapLoginCompleteCmd struct {
	//Dummy byte
}

//struct PROTO_NC_MAP_LOGOUT_CMD
//{
//  unsigned __int16 handle;
//};
type MapLogoutCmd struct {
	Handle uint16
}

//struct PROTO_NC_MAP_FIELD_ATTRIBUTE_CMD
//{
//  FIELD_MAP_TYPE eFieldMapType;
//};
type NcMapFieldAttributeCmd struct {
	FieldMapType uint32
}

//struct PROTO_NC_MAP_CAN_USE_REVIVEITEM_CMD
//{
//  char bCanUseReviveItem;
//};
type NcMapCanUseReviveItemCmd struct {
	CanUse byte
}

//struct PROTO_NC_MAP_TOWNPORTAL_REQ
//{
//  char portalindex;
//};
type NcMapTownPortalReq struct {
	PortalIndex byte
}

//struct PROTO_NC_MAP_TOWNPORTAL_ACK
//{
//  unsigned __int16 err;
//};
type NcMapTownPortalAck struct {
	Err uint16
}

//NC_MAP_LOGIN_ACK
type NcMapLoginAck NcCharMapLoginAck
