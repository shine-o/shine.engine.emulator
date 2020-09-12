package zone

import "time"

func (p *player) heartbeat() {
	log.Infof("[player_ticks] heartbeat for player %v", p.view.name)
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
			log.Infof("[player_ticks] sending heartbeat for player %v", p.view.name)
			ncMiscHeartBeatReq(p)
		}
	}
}

func (p *player) persistPosition() {
	log.Infof("[player_ticks] persistPosition for player %v", p.view.name)
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
			log.Infof("[player_ticks] persisting position for player %v", p.view.name)
			pppe := persistPlayerPositionEvent{
				p: p,
			}
			zoneEvents[persistPlayerPosition] <- &pppe
		}
	}
}


func (p *player) nearbyEntities() {
	log.Infof("[player_ticks] nearbyEntities for player %v", p.view.name)
	tick := time.NewTicker(2 * time.Second)

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
			log.Infof("[player_ticks] persisting position for player %v", p.view.name)
			p.Lock()
			for handle, _ := range p.baseEntity.nearbyEntities {
				e := p.baseEntity.nearbyEntities[handle]
				if !inRange(&p.baseEntity, e) {
					// send NC_BRIEFINFO_BRIEFINFODELETE_CMD
				}
			}
			p.Unlock()
		}
	}
}
