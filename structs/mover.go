package structs

import (
	"encoding/json"
	"reflect"
)

//struct PROTO_NC_MOVER_MOVESPEED_CMD
//{
//  unsigned __int16 nMoverHandle;
//  unsigned __int16 nWalk;
//  unsigned __int16 nRun;
//};
type NcMoverMoveSpeedCmd struct {
	MoverHandle uint16
	Walk        uint16
	Run         uint16
}

func (nc *NcMoverMoveSpeedCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMoverMoveSpeedCmd) PdbType() string {
	return `
	struct PROTO_NC_MOVER_MOVESPEED_CMD
	{
	  unsigned __int16 nMoverHandle;
	  unsigned __int16 nWalk;
	  unsigned __int16 nRun;
	};
`
}

//struct PROTO_NC_MOVER_HUNGRY_CMD
//{
//  unsigned __int16 nHungry;
//};
type NcMoverHungryCmd struct {
	Hungry uint16
}

func (nc *NcMoverHungryCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMoverHungryCmd) PdbType() string {
	return `
	struct PROTO_NC_MOVER_HUNGRY_CMD
	{
	  unsigned __int16 nHungry;
	};
`
}

//struct PROTO_NC_MOVER_RIDE_ON_CMD
//{
//  unsigned __int16 nMoverHandle;
//  char nSlot;
//  char nGrade;
//};
type NcMoverRideOnCmd struct {
	MoverHandle uint16
	Slot        byte
	Grade       byte
}

func (nc *NcMoverRideOnCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMoverRideOnCmd) PdbType() string {
	return `
	struct PROTO_NC_MOVER_RIDE_ON_CMD
	{
	  unsigned __int16 nMoverHandle;
	  char nSlot;
	  char nGrade;
	};
`
}

//struct PROTO_NC_MOVER_SOMEONE_RIDE_ON_CMD
//{
//  unsigned __int16 nHandle;
//  unsigned __int16 nMoverHandle;
//  char nSlot;
//  char nGrade;
//};
type NcMoverSomeoneRideOnCmd struct {
	Handle      uint16
	MoverHandle uint16
	Slot        byte
	Grade       byte
}

func (nc *NcMoverSomeoneRideOnCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMoverSomeoneRideOnCmd) PdbType() string {
	return `
	struct PROTO_NC_MOVER_SOMEONE_RIDE_ON_CMD
	{
	  unsigned __int16 nHandle;
	  unsigned __int16 nMoverHandle;
	  char nSlot;
	  char nGrade;
	};
`
}

//struct PROTO_NC_MOVER_SOMEONE_RIDE_OFF_CMD
//{
//  unsigned __int16 nHandle;
//};
type NcMoverSomeoneRideOffCmd struct {
	Handle uint16
}

func (nc *NcMoverSomeoneRideOffCmd) String() string {
	sd, err := json.Marshal(nc)
	if err != nil {
		log.Errorf("converting struct %v to json resulted in error: %v", reflect.TypeOf(nc).String(), err)
	}
	return string(sd)
}

func (nc *NcMoverSomeoneRideOffCmd) PdbType() string {
	return `
	struct PROTO_NC_MOVER_SOMEONE_RIDE_OFF_CMD
	{
	  unsigned __int16 nHandle;
	};
`
}
