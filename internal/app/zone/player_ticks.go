package zone

import "time"

func (p *player) heartbeatTicker() {
	log.Infof("[player_ticks] heartbeatTicker  for player %v", p.view.name)
	tick := time.NewTicker(5 * time.Second)
	p.Lock()
	p.tickers = append(p.tickers, tick)
	p.Unlock()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if p == nil {
				return
			}
			log.Infof("sending heartbeat for player %v", p.view.name)
			ncMiscHeartBeatReq(p)
		}
	}
}

func (p *player) persistPositionTicker() {
	log.Infof("[player_ticks] heartbeatTicker for player %v", p.view.name)
	tick := time.NewTicker(4 * time.Second)

	p.Lock()
	p.tickers = append(p.tickers, tick)
	p.Unlock()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if p == nil {
				return
			}
			log.Infof("persisting position for player %v", p.view.name)
			pppe := persistPlayerPositionEvent{
				p: p,
			}
			zoneEvents[persistPlayerPosition] <- &pppe
		}
	}
}
