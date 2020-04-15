package structs

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

//struct PROTO_NC_MOVER_HUNGRY_CMD
//{
//  unsigned __int16 nHungry;
//};
type NcMoverHungryCmd struct {
	Hungry uint16
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

//struct PROTO_NC_MOVER_SOMEONE_RIDE_OFF_CMD
//{
//  unsigned __int16 nHandle;
//};
type NcMoverSomeoneRideOffCmd struct {
	Handle uint16
}