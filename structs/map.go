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
