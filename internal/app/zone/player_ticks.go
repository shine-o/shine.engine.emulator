package zone

import "time"

func (p *player) heartbeat() {
	log.Infof("[player_ticks] heartbeat ticker/worker for player %v", p.view.name)
	tick := time.NewTicker(5 * time.Second)
	for {
		select {
		case <- p.recv[heartbeatUpdate]:
			p.Lock()
			p.conn.lastHeartBeat = time.Now()
			log.Infof("updating heartbeat for player %v", p.view.name)
			p.Unlock()
		case <-p.recv[heartbeatStop]:
			select {
			case p.conn.close <- true:
				tick.Stop()
				return
			default:
				tick.Stop()
				return
			}
		case <-tick.C:
			log.Infof("sending heartbeat for player %v", p.view.name)
			if p == nil {
				tick.Stop()
				return
			}
			ncMiscHeartBeatReq(p)
		//default:
		//	continue
		}
	}
}
