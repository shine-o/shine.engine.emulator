package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"time"
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

// remove entities that are outside the view range of the player
func (p *player) nearbyPlayersMaintenance(zm *zoneMap) {

	log.Infof("[player_ticks] nearbyEntities for handle %v", p.getHandle())
	tick := time.NewTicker(200 * time.Millisecond)
	p.ticks.Lock()
	p.ticks.list = append(p.ticks.list, tick)
	p.ticks.Unlock()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			for ap := range zm.entities.players.all() {
				go func(p1, p2 *player) {
					if p2.getHandle() == p1.getHandle() {
						return
					}

					p1.proximity.RLock()
					_, exists := p1.proximity.players[p2.getHandle()]
					p1.proximity.RUnlock()

					if !exists && playerInRange(p1, p2) {
						nc := ncBriefInfoLoginCharacterCmd(p2)
						networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_LOGINCHARACTER_CMD, &nc)
					}
				}(p, ap)
			}

			for ap := range p.adjacentPlayers() {
				go checkRemoval(p, ap)
			}

		}
	}
}

func (p *player) nearbyNpcsMaintenance(zm *zoneMap) {
	log.Infof("[player_ticks] nearbyNpcsMaintenance for handle %v", p.getHandle())
	tick := time.NewTicker(200 * time.Millisecond)

	p.ticks.Lock()
	p.ticks.list = append(p.ticks.list, tick)
	p.ticks.Unlock()

	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			// for each monster
			// if nearby, add to known nearby
			for nn := range zm.entities.npcs.all() {
				go func(p *player, n *npc) {
					p.proximity.RLock()
					_, exists := p.proximity.npcs[n.getHandle()]
					p.proximity.RUnlock()
					if !exists && npcInRange(p, n) {
						nc := ncBriefInfoRegenMobCmd(n)
						networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_REGENMOB_CMD, &nc)
					}
				}(p, nn)
			}

			for an := range p.adjacentNpcs() {
				if an.monster {
					go func(p *player, n *npc) {
						if !npcInRange(p, n) {
							mh := n.getHandle()
							p.proximity.Lock()
							delete(p.proximity.npcs, mh)
							p.proximity.Unlock()
							nc := structs.NcBriefInfoDeleteHandleCmd{
								Handle: mh,
							}
							networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_BRIEFINFODELETE_CMD, nc)
						}
					}(p, an)
				}
			}
		}
	}
}

// if foreign player timed out or is not in range
// send packet to the client to notify of player disappearance
func checkRemoval(p1, p2 *player) {
	fh := p2.getHandle()

	if justSpawned(p2) {
		return
	}

	if lastHeartbeat(p2) > playerHeartbeatLimit {
		p1.removeAdjacentPlayer(fh)
		return
	}

	if !playerInRange(p1, p2) {
		nc := structs.NcBriefInfoDeleteHandleCmd{
			Handle: p2.getHandle(),
		}
		networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_BRIEFINFODELETE_CMD, &nc)
		p1.removeAdjacentPlayer(fh)
	}
}

func justSpawned(p *player) bool {
	p.state.RLock()
	js := p.state.justSpawned
	p.state.RUnlock()
	return js
}
