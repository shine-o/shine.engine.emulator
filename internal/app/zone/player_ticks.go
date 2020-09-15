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
			go ncMiscHeartBeatReq(p)
		}
	}
}

func (p *player) persistPosition() {
	log.Infof("[player_ticks] persistPosition for handle %v", p.handle)
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
			log.Infof("[player_ticks] persisting position for handle %v", p.handle)
			pppe := persistPlayerPositionEvent{
				p: p,
			}
			zoneEvents[persistPlayerPosition] <- &pppe
		}
	}
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
			zm.entities.players.RLock()
			for i := range zm.entities.players.active {

				foreignPlayer := zm.entities.players.active[i]

				foreignPlayer.RLock()
				p.RLock()
				fH := foreignPlayer.handle
				pH := p.handle
				foreignPlayer.RUnlock()
				p.RUnlock()

				if fH == pH {
					continue
				}

				if playerInRange(p, foreignPlayer) {

					nc := foreignPlayer.ncBriefInfoLoginCharacterCmd()

					p.RLock()
					_, exists := p.knownNearbyPlayers[fH]
					p.RUnlock()

					if !exists {
						go ncBriefInfoLoginCharacterCmd(p, &nc)
					}

				}
			}

			zm.entities.players.RUnlock()

			p.RLock()

			log.Infof("[player_ticks] removing out of range entities for handle %v", p.handle)

			for i := range p.knownNearbyPlayers {

				go func(pHandle uint16) {

					var removePlayer = func(pHandle uint16, vp * player) {
						vp.Lock()
						delete(vp.knownNearbyPlayers, pHandle)
						vp.Unlock()
					}

					foreignPlayer := p.knownNearbyPlayers[pHandle]

					foreignPlayer.RLock()
					lastHeartBeat := time.Since(foreignPlayer.conn.lastHeartBeat).Seconds()
					foreignPlayer.RUnlock()

					if lastHeartBeat > playerHeartbeatLimit {
						go removePlayer(pHandle, p)
						return
					}

					if !playerInRange(p, foreignPlayer) {
						nc := structs.NcBriefInfoDeleteHandleCmd{
							Handle: foreignPlayer.handle,
						}
						go ncBriefInfoDeleteHandleCmd(p, &nc)
						go removePlayer(pHandle, p)
					}
				}(i)

			}
			p.RUnlock()
		}
	}
}
