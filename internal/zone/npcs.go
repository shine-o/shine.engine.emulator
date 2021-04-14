package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"sync"
)

type npcs struct {
	*handler
	active map[uint16]*npc
	sync.RWMutex
}

func (n *npcs) all() <-chan *npc {

	n.RLock()
	ch := make(chan *npc, len(n.active))
	n.RUnlock()

	go func(send chan<- *npc) {
		n.RLock()
		for _, ap := range n.active {
			send <- ap
		}
		n.RUnlock()
		close(send)
	}(ch)

	return ch
}

func (n *npcs) get(h uint16) *npc {
	n.RLock()
	npc := n.active[h]
	n.RUnlock()
	return npc
}

func (n *npcs) remove(h uint16) {
	n.Lock()
	delete(n.active, h)
	n.Unlock()
}

func (n *npcs) add(ap *npc) {
	h := ap.getHandle()
	n.Lock()
	n.active[h] = ap
	n.handler.usedHandles[h] = true
	n.Unlock()
}

func npcInRange(p *player, n *npc) bool {
	h := n.getHandle()
	yes := entityInRange(p.baseEntity.current, n.baseEntity.current)

	if yes {
		p.proximity.Lock()
		p.proximity.npcs[h] = n
		p.proximity.Unlock()
		return true
	}
	return false
}

func knownNpc(p *player, nh uint16) bool {
	p.proximity.RLock()
	_, ok := p.proximity.npcs[nh]
	p.proximity.RUnlock()
	if ok {
		return true
	}
	return false
}

func adjacentNpcsInform(p *player, zm *zoneMap) {
	for m := range zm.entities.npcs.all() {
		go func(p *player, n *npc) {
			if !knownNpc(p, n.getHandle()) {
				if npcInRange(p, n) {
					nc := ncBriefInfoRegenMobCmd(n)
					networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_REGENMOB_CMD, &nc)
				}
			}
		}(p, m)
	}
}

