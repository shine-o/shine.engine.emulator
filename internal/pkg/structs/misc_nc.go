package structs

// struct PROTO_NC_MISC_SEED_ACK
type NcMiscSeedAck struct {
	Seed uint16
}

// struct PROTO_NC_MISC_GAMETIME_ACK
type NcMiscGameTimeAck struct {
	Hour   byte
	Minute byte
	Second byte
}

// struct PROTO_NC_MISC_HEARTBEAT_ACK
type NcMiscHeartBeatAck struct{}

//struct PROTO_NC_MISC_SERVER_TIME_NOTIFY_CMD
type NcMiscServerTimeNotifyCmd struct {
	Time     TM
	TimeZone byte
}
