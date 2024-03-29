package zone

import "time"

func (zm *zoneMap) removeInactiveHandles() {
	log.Infof("[map_ticks] removeInactiveHandles ticker/worker for map %v", zm.data.Info.MapName)
	tick := time.Tick(100 * time.Millisecond)
	for {
		select {
		case <-tick:
			select {
			case zm.events.send[playerHandleMaintenance] <- &emptyEvent{}:
				break
			default:
				log.Infof("failed executing playerHandleMaintenance for %v", zm.data.Info.MapName)
			}
		}
	}
}
