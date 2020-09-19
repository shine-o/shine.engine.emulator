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

			for ap := range zm.entities.activePlayers() {
				go func(p1, p2 * player) {
					if p2.getHandle() == p1.getHandle() {
						return
					}

					if playerInRange(p1, p2) {

						p1.RLock()
						_, exists := p1.players[p2.getHandle()]
						p1.RUnlock()

						if !exists {
							p2.RLock() //todo move locks into the method
							nc := p2.ncBriefInfoLoginCharacterCmd()
							p2.RUnlock()
							ncBriefInfoLoginCharacterCmd(p1, &nc)
						}
					}
				}(p, ap)
			}

			for ap := range p.adjacentPlayers() {
				go checkRemoval(p, ap)
			}

		}
	}
}

func (p *player) nearbyMonsters(zm *zoneMap) {
	log.Infof("[player_ticks] nearbyMonsters for handle %v", p.handle)
	tick := time.NewTicker(200 * time.Millisecond)

	p.Lock()
	p.tickers = append(p.tickers, tick)
	p.Unlock()
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			// for each monster
			// if nearby, add to known nearby
		}
	}
}

func removeNearbyPlayer(targetHandle uint16, vp *player) {
	vp.Lock()
	delete(vp.players, targetHandle)
	vp.Unlock()
}

// if foreign player timed out or is not in range
// send packet to the client to notify of player disappearance
func checkRemoval(p1, p2 *player) {
	fh := p2.getHandle()

	if lastHeartbeat(p2) > playerHeartbeatLimit {
		go removeNearbyPlayer(fh, p1)
		return
	}

	if !playerInRange(p1, p2) {
		nc := structs.NcBriefInfoDeleteHandleCmd{
			Handle: p2.getHandle(),
		}
		go ncBriefInfoDeleteHandleCmd(p1, &nc)
		go removeNearbyPlayer(fh, p1)
	}
}

