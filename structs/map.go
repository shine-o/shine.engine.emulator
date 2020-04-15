package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_MAP_LOGIN_REQ
//{
//  PROTO_NC_CHAR_ZONE_CHARDATA_REQ chardata;
//  Name8 checksum[50];
//};
type NcMapLoginReq struct {
	CharData NcZoneCharDataReq
	CheckSum [49]Name8
}

func (nc *NcMapLoginReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMapLoginReq) PdbType() string {
	return `
	struct PROTO_NC_MAP_LOGIN_REQ
	{
	  PROTO_NC_CHAR_ZONE_CHARDATA_REQ chardata;
	  Name8 checksum[50];
	};
`
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

func (nc *NcMapLoginCompleteCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMapLoginCompleteCmd) PdbType() string {
	return `
	struct PROTO_NC_MAP_LOGINCOMPLETE_CMD
	{
	  char dummy[1];
	};
`
}

//struct PROTO_NC_MAP_LOGOUT_CMD
//{
//  unsigned __int16 handle;
//};
type MapLogoutCmd struct {
	Handle uint16
}

func (nc *MapLogoutCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *MapLogoutCmd) PdbType() string {
	return `
	struct PROTO_NC_MAP_LOGOUT_CMD
	{
	  unsigned __int16 handle;
	};
`
}

//struct PROTO_NC_MAP_FIELD_ATTRIBUTE_CMD
//{
//  FIELD_MAP_TYPE eFieldMapType;
//};
type NcMapFieldAttributeCmd struct {
	FieldMapType uint32
}

func (nc *NcMapFieldAttributeCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMapFieldAttributeCmd) PdbType() string {
	return `
	struct PROTO_NC_MAP_FIELD_ATTRIBUTE_CMD
	{
	  FIELD_MAP_TYPE eFieldMapType;
	};
`
}

//struct PROTO_NC_MAP_CAN_USE_REVIVEITEM_CMD
//{
//  char bCanUseReviveItem;
//};
type NcMapCanUseReviveItemCmd struct {
	CanUse byte
}

func (nc *NcMapCanUseReviveItemCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMapCanUseReviveItemCmd) PdbType() string {
	return `
	struct PROTO_NC_MAP_CAN_USE_REVIVEITEM_CMD
	{
	  char bCanUseReviveItem;
	};
`
}

//struct PROTO_NC_MAP_TOWNPORTAL_REQ
//{
//  char portalindex;
//};
type NcMapTownPortalReq struct {
	PortalIndex byte
}

func (nc *NcMapTownPortalReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMapTownPortalReq) PdbType() string {
	return `
	struct PROTO_NC_MAP_TOWNPORTAL_REQ
	{
	  char portalindex;
	};
`
}

//struct PROTO_NC_MAP_TOWNPORTAL_ACK
//{
//  unsigned __int16 err;
//};
type NcMapTownPortalAck struct {
	Err uint16
}

func (nc *NcMapTownPortalAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMapTownPortalAck) PdbType() string {
	return `
	struct PROTO_NC_MAP_TOWNPORTAL_ACK
	{
	  unsigned __int16 err;
	};
`
}