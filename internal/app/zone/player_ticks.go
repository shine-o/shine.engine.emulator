package zone

import (
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"time"
)

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

// remove entities that are outside the view range of the player
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
			log.Infof("[player_ticks] removing out of range entities for player %v", p.view.name)
			p.Lock()
			for i := range p.baseEntity.nearbyEntities {
				e := p.baseEntity.nearbyEntities[i]
				if !inRange(&p.baseEntity, e) {
					nc := structs.NcBriefInfoDeleteHandleCmd{
						Handle: e.handle,
					}

					ncBriefInfoDeleteHandleCmd(p, &nc)

					delete(p.baseEntity.nearbyEntities, i)
				}
			}
			p.Unlock()
		}
	}
}
