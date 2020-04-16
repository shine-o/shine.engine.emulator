package service

type ZoneInfo struct {
	IP string
	Port int32
}

type registeredMaps map[string]ZoneInfo