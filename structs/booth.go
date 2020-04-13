package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_BOOTH_ENTRY_REQ
//{
//  unsigned __int16 booth;
//};
type NcBoothEntryReq struct {
	Booth uint16
}

func (nc *NcBoothEntryReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBoothEntryReq) PdbType() string {
	return `
	struct PROTO_NC_BOOTH_ENTRY_REQ
	{
	  unsigned __int16 booth;
	};
`
}

//struct PROTO_NC_BOOTH_SOMEONEOPEN_CMD
//{
//  unsigned __int16 handle;
//  CHARBRIEFINFO_CAMP tent;
//  char issell;
//  STREETBOOTH_SIGNBOARD signboard;
//};
type NcBoothSomeoneOpenCmd struct {
	Handle    uint16
	Tent      CharBriefInfoCamp
	IsSelling byte
	Sign      StreetBoothSignBoard
}

func (nc *NcBoothSomeoneOpenCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBoothSomeoneOpenCmd) PdbType() string {
	return `
	struct PROTO_NC_BOOTH_SOMEONEOPEN_CMD
	{
	  unsigned __int16 handle;
	  CHARBRIEFINFO_CAMP tent;
	  char issell;
	  STREETBOOTH_SIGNBOARD signboard;
	};
`
}

//struct PROTO_NC_BOOTH_REFRESH_REQ
//{
//  unsigned __int16 booth;
//};
type NcBoothRefreshReq struct {
	Booth uint16
}

func (nc *NcBoothRefreshReq) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBoothRefreshReq) PdbType() string {
	return `
struct PROTO_NC_BOOTH_REFRESH_REQ
{
  unsigned __int16 booth;
};
`
}

//struct PROTO_NC_BOOTH_ENTRY_SELL_ACK
//{
//  unsigned __int16 err;
//  unsigned __int16 boothhandle;
//  char numofitem;
//  PROTO_NC_BOOTH_ENTRY_SELL_ACK::BoothItemList items[];
//};
type NcBoothEntrySellAck struct {
	Err         uint16
	BoothHandle uint16
	NumOfItems  byte
	Items       []NcBoothEntrySellAckItemList `struct:"sizefrom=NumOfItems"`
}

func (nc *NcBoothEntrySellAck) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBoothEntrySellAck) PdbType() string {
	return `
	struct __unaligned __declspec(align(1)) PROTO_NC_BOOTH_ENTRY_SELL_ACK
	{
	  unsigned __int16 err;
	  unsigned __int16 boothhandle;
	  char numofitem;
	  PROTO_NC_BOOTH_ENTRY_SELL_ACK::BoothItemList items[];
	};
`
}

//struct PROTO_NC_BOOTH_SEARCH_BOOTH_CLOSED_CMD
//{
//  unsigned __int16 nClosedBoothOwnerHandle;
//};
type NcBoothSearchBoothClosedCmd struct {
	ClosedBoothOwnerHandle uint16
}

func (nc *NcBoothSearchBoothClosedCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcBoothSearchBoothClosedCmd) PdbType() string {
	return `
	struct PROTO_NC_BOOTH_SEARCH_BOOTH_CLOSED_CMD
	{
	  unsigned __int16 nClosedBoothOwnerHandle;
	};
`
}
