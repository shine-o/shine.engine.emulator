package zonemaster

type zoneInfo struct {
	IP   string
	Port int32
}

type registeredMaps map[int32]zoneInfo
