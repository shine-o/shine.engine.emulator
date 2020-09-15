package zone

import "time"

func (zm *zoneMap) removeInactiveHandles() {
	log.Infof("[map_ticks] heartbeat ticker/worker for map %v", zm.data.Info.MapName)
	tick := time.Tick(7 * time.Second)
	for {
		select {
		case <-tick:
			select {
			case zm.send[playerHandleMaintenance] <- &emptyEvent{}:
				log.Infof("executing playerHandleMaintenance for %v", zm.data.Info.MapName)
				break
			default:
				log.Infof("failed toexecuting playerHandleMaintenance for %v", zm.data.Info.MapName)
			}
		}
	}
}
