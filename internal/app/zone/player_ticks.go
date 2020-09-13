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
	//log.Infof("[player_ticks] persistPosition for handle %v", p.handle)
	//tick := time.NewTicker(4 * time.Second)
	//
	//p.Lock()
	//p.tickers = append(p.tickers, tick)
	//p.Unlock()
	//defer tick.Stop()
	//
	//for {
	//	select {
	//	case <-tick.C:
	//		if p == nil {
	//			return
	//		}
	//		log.Infof("[player_ticks] persisting position for handle %v", p.handle)
	//		pppe := persistPlayerPositionEvent{
	//			p: p,
	//		}
	//		zoneEvents[persistPlayerPosition] <- &pppe
	//	}
	//}
}

// remove entities that are outside the view range of the player
func (p *player) nearbyPlayers(zm *zoneMap) {
	log.Infof("[player_ticks] nearbyEntities for handle %v", p.handle)
	tick := time.NewTicker(200 * time.Millisecond)

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

			// also send data to the player that is running about nearby players
			for i := range zm.entities.players.active {
				foreignPlayer := zm.entities.players.active[i]
				foreignPlayer.Lock()
				if foreignPlayer.handle == p.handle {
					foreignPlayer.Unlock()
					continue
				}

				if playerInRange(p, foreignPlayer) {
					nc := foreignPlayer.ncBriefInfoLoginCharacterCmd()
					_, ok := p.knownNearbyPlayers[foreignPlayer.handle]
					if ok {
						foreignPlayer.Unlock()
						continue
					}
					go ncBriefInfoLoginCharacterCmd(p, &nc)
				}

				foreignPlayer.Unlock()
			}

			log.Infof("[player_ticks] removing out of range entities for handle %v", p.handle)
			p.Lock()
			for i := range p.knownNearbyPlayers {
				foreignPlayer := p.knownNearbyPlayers[i]

				lastHeartBeat := time.Since(foreignPlayer.conn.lastHeartBeat).Seconds()
				if lastHeartBeat > playerHeartbeatLimit {
					delete(p.knownNearbyPlayers, i)
					continue
				}

				if !playerInRange(p, foreignPlayer) {
					nc := structs.NcBriefInfoDeleteHandleCmd{
						Handle: foreignPlayer.handle,
					}

					go ncBriefInfoDeleteHandleCmd(p, &nc)

					delete(p.knownNearbyPlayers, i)
				}
			}
			p.Unlock()
		}
	}
}
