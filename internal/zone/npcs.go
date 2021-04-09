package zone

import (
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
	n.Lock()
	n.active[ap.handle] = ap
	n.handler.usedHandles[ap.handle] = true
	n.Unlock()
}

func npcInRange(p *player, n *npc) bool {
	p.RLock()
	n.RLock()
	yes := entityInRange(p.baseEntity, n.baseEntity)
	p.RUnlock()
	n.RUnlock()

	if yes {
		p.Lock()
		p.npcs[n.handle] = n
		p.Unlock()
		return true
	}
	return false
}
