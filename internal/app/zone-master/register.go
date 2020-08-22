package zone_master

type ZoneInfo struct {
	IP   string
	Port int32
}

type registeredMaps map[int32]ZoneInfo
