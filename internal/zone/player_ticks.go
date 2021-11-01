package zone

import (
	"time"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
)

func (p *player) heartbeat() {
	log.Infof("[player_ticks] heartbeat for player %v", p.view.name)
	tick := time.NewTicker(5 * time.Second)

	p.ticks.Lock()
	p.ticks.list = append(p.ticks.list, tick)
	p.ticks.Unlock()

	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			if p == nil {
				return
			}
			log.Infof("[player_ticks] sending heartbeat for player %v", p.view.name)
			go networking.Send(p.conn.outboundData, networking.NC_MISC_HEARTBEAT_REQ, nil)
		}
	}
}

func (p *player) persistPosition() {
	h := p.getHandle()
	log.Infof("[player_ticks] persistPosition for handle %v", h)
	tick := time.NewTicker(4 * time.Second)

	p.ticks.Lock()
	p.ticks.list = append(p.ticks.list, tick)
	p.ticks.Unlock()

	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			log.Infof("[player_ticks] persisting position for handle %v", h)
			e := persistPlayerPositionEvent{
				p: p,
			}
			zoneEvents[persistPlayerPosition] <- &e
		}
	}
}
