package service

import "time"

func (zm *zoneMap) removeInactiveHandles() {
	log.Infof("[map_ticks] heartbeat ticker/worker for map %v", zm.data.Info.MapName)
	tick := time.Tick(5 *  time.Second)
	for {
		select {
		case <- tick:
			zm.send[playerHandleMaintenance] <- &emptyEvent{}
		}
	}
}