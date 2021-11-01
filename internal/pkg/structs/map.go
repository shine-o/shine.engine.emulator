package structs

// struct PROTO_NC_MAP_LOGIN_REQ
type NcMapLoginReq struct {
	CharData NcZoneCharDataReq
	CheckSum [54]Name8
}

// struct PROTO_NC_CHAR_ZONE_CHARDATA_REQ
type NcZoneCharDataReq struct {
	WorldManager uint16
	CharID       Name5
}

// struct PROTO_NC_MAP_LOGINCOMPLETE_CMD
type NcMapLoginCompleteCmd struct { // Dummy byte
}

// struct PROTO_NC_MAP_LOGOUT_CMD
type MapLogoutCmd struct {
	Handle uint16
}

// struct PROTO_NC_MAP_FIELD_ATTRIBUTE_CMD
type NcMapFieldAttributeCmd struct {
	FieldMapType uint32
}

// struct PROTO_NC_MAP_CAN_USE_REVIVEITEM_CMD
type NcMapCanUseReviveItemCmd struct {
	CanUse byte
}

// struct PROTO_NC_MAP_TOWNPORTAL_REQ
type NcMapTownPortalReq struct {
	PortalIndex byte
}

// struct PROTO_NC_MAP_TOWNPORTAL_ACK
type NcMapTownPortalAck struct {
	Err uint16
}

// NC_MAP_LOGIN_ACK
type NcMapLoginAck NcCharMapLoginAck

// struct PROTO_NC_MAP_LOGOUT_CMD
type NcMapLogoutCmd struct {
	Handle uint16
}
