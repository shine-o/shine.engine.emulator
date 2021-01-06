package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
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
			go networking.Send(p.conn.outboundData, networking.NC_MISC_HEARTBEAT_REQ, nil)
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
			log.Infof("[player_ticks] persisting position for handle %v", p.handle)
			pppe := persistPlayerPositionEvent{
				p: p,
			}
			zoneEvents[persistPlayerPosition] <- &pppe
		}
	}
}

// remove entities that are outside the view range of the player
func (p *player) nearbyPlayersMaintenance(zm *zoneMap) {

	log.Infof("[player_ticks] nearbyEntities for handle %v", p.handle)
	tick := time.NewTicker(200 * time.Millisecond)
	p.Lock()
	p.tickers = append(p.tickers, tick)
	p.Unlock()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			for ap := range zm.entities.players.all() {
				go func(p1, p2 *player) {
					if p2.getHandle() == p1.getHandle() {
						return
					}

					p1.RLock()
					_, exists := p1.players[p2.getHandle()]
					p1.RUnlock()

					if !exists && playerInRange(p1, p2) {
						p2.RLock() //todo move locks into the method
						nc := p2.ncBriefInfoLoginCharacterCmd()
						p2.RUnlock()
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

func (p *player) nearbyMonstersMaintenance(zm *zoneMap) {
	log.Infof("[player_ticks] nearbyMonstersMaintenance for handle %v", p.handle)
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
			for am := range zm.entities.monsters.all() {
				go func(p *player, m *monster) {
					p.RLock()
					_, exists := p.monsters[m.getHandle()]
					p.RUnlock()
					if !exists && monsterInRange(p, m) {
						nc := m.ncBriefInfoRegenMobCmd()
						networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_REGENMOB_CMD, &nc)
					}
				}(p, am)
			}

			for am := range p.adjacentMonsters() {
				go func(p *player, m *monster) {
					if !monsterInRange(p, m) {
						mh := m.getHandle()
						p.Lock()
						delete(p.monsters, mh)
						p.Unlock()
						nc := structs.NcBriefInfoDeleteHandleCmd{
							Handle: mh,
						}
						networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_BRIEFINFODELETE_CMD, nc)
					}
				}(p, am)
			}
		}
	}
}

func (p *player) nearbyNPCMaintenance(zm *zoneMap) {
	log.Infof("[player_ticks] nearbyNPCs for handle %v", p.handle)
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
			for an := range zm.entities.npcs.all() {
				go func(p *player, n *npc) {
					p.RLock()
					_, exists := p.npcs[n.getHandle()]
					p.RUnlock()
					if !exists && npcInRange(p, n) {
						nc := n.ncBriefInfoRegenMobCmd()
						networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_REGENMOB_CMD, &nc)
					}
				}(p, an)
			}

			//for am := range p.adjacentMonsters() {
			//	go func(p *player, m *monster) {
			//		if !monsterInRange(p, m) {
			//			mh := m.getHandle()
			//			p.Lock()
			//			delete(p.monsters, mh)
			//			p.Unlock()
			//			nc := structs.NcBriefInfoDeleteHandleCmd{
			//				Handle: mh,
			//			}
			//			ncBriefInfoDeleteHandleCmd(p, &nc)
			//		}
			//	}(p, am)
			//}
		}
	}
}

// if foreign player timed out or is not in range
// send packet to the client to notify of player disappearance
func checkRemoval(p1, p2 *player) {
	fh := p2.getHandle()

	if p2.spawned() {
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
