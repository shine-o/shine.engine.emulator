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

				go func(vp, fp * player) {
					fp.RLock()
					fH := fp.handle
					fp.RUnlock()

					vp.RLock()
					pH := vp.handle
					vp.RUnlock()

					if fH == pH {
						return
					}

					if playerInRange(vp, fp) {
						fp.RLock()
						nc := fp.ncBriefInfoLoginCharacterCmd()
						fp.RUnlock()

						vp.RLock()
						_, exists := vp.knownNearbyPlayers[fH]
						vp.RUnlock()

						if !exists {
							go ncBriefInfoLoginCharacterCmd(vp, &nc)
						}
					}
				}(p, foreignPlayer)

			}

			zm.entities.players.RUnlock()

			p.RLock()
			log.Infof("[player_ticks] removing out of range entities for handle %v", p.handle)
			for i := range p.knownNearbyPlayers {
				go checkRemoval(i, p)
			}
			p.RUnlock()
		}
	}
}

func removeNearbyPlayer(targetHandle uint16, vp * player) {
	vp.Lock()
	delete(vp.knownNearbyPlayers, targetHandle)
	vp.Unlock()
}

// if foreign player timed out or is not in range
// send packet to the client to notify of player disappearance
func checkRemoval(fHandle uint16, viewerPlayer * player) {

	viewerPlayer.RLock()
	foreignPlayer := viewerPlayer.knownNearbyPlayers[fHandle]
	viewerPlayer.RUnlock()

	if lastHeartbeat(foreignPlayer) > playerHeartbeatLimit {
		go removeNearbyPlayer(fHandle, viewerPlayer)
		return
	}

	if !playerInRange(viewerPlayer, foreignPlayer) {
		nc := structs.NcBriefInfoDeleteHandleCmd{
			Handle: foreignPlayer.handle,
		}
		go ncBriefInfoDeleteHandleCmd(viewerPlayer, &nc)
		go removeNearbyPlayer(fHandle, viewerPlayer)
	}
}

func lastHeartbeat(p *player) float64 {
	p.RLock()
	lastHeartBeat := time.Since(p.conn.lastHeartBeat).Seconds()
	p.RUnlock()
	return lastHeartBeat
}
