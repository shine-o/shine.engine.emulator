package service

import "time"

func (p *player) heartbeat() {
	log.Infof("[player_ticks] heartbeat ticker/worker for player %v", p.view.name)
	tick := time.Tick(5 * time.Second)
	for {
		select {
		case <- p.recv[heartbeatUpdate]:
			p.conn.lastHeartBeat = time.Now()
		case <- p.recv[heartbeatMissing]:
			return
		case <- tick:
			ncMiscHeartBeatReq(p)
		}
	}
}