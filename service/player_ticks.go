package service

import "time"

func (p *player) heartbeat() {
	tick := time.Tick(5 * time.Second)
	for {
		select {
		case <- p.recv[heartbeatMissing]:
			return
		case <- tick:
			// NC_MISC_HEARTBEAT_REQ
			ncMiscHeartBeatReq(p)
		}
	}
}
