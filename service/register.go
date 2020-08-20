package service

type WorldInfo struct {
	ID int
	Name string
	IP string
	Port int32
}

type registeredWorlds map[int32]WorldInfo