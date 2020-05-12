package service

import "time"

func (p *player) heartbeat() {
	log.Infof("[player_ticks] heartbeat ticker/worker for player %v", p.view.name)
	tick := time.NewTicker(5 * time.Second)
	for {
		select {
		case <- p.recv[heartbeatUpdate]:
			p.conn.lastHeartBeat = time.Now()
		case <- p.recv[heartbeatMissing]:
			tick.Stop()
			select {
			case p.conn.close <- true:
			default:
				return
			}
		case <- tick.C:
			log.Infof("sending heartbeat for player %v",p.view.name)
			ncMiscHeartBeatReq(p)
		}
	}
}