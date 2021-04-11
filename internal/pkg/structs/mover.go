package structs

//struct PROTO_NC_MOVER_MOVESPEED_CMD
type NcMoverMoveSpeedCmd struct {
	MoverHandle uint16
	Walk        uint16
	Run         uint16
}

//struct PROTO_NC_MOVER_HUNGRY_CMD
type NcMoverHungryCmd struct {
	Hungry uint16
}

//struct PROTO_NC_MOVER_RIDE_ON_CMD
type NcMoverRideOnCmd struct {
	MoverHandle uint16
	Slot        byte
	Grade       byte
}

//struct PROTO_NC_MOVER_SOMEONE_RIDE_ON_CMD
type NcMoverSomeoneRideOnCmd struct {
	Handle      uint16
	MoverHandle uint16
	Slot        byte
	Grade       byte
}

//struct PROTO_NC_MOVER_SOMEONE_RIDE_OFF_CMD
type NcMoverSomeoneRideOffCmd struct {
	Handle uint16
}
