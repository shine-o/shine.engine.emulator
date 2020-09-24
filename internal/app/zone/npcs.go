package zone

import (
	"sync"
)

const npcHandleMin uint16 = 17000
const npcHandleMax uint16 = 27000
const npcAttemptsMax uint16 = 1500

type npcs struct {
	handler
	active      map[uint16]*npc
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
	delete(n.handler.usedHandles, h)
	n.Unlock()
}

func (n *npcs) add(ap *npc) {
	n.Lock()
	n.active[ap.handle] = ap
	n.handler.usedHandles[ap.handle] = true
	n.Unlock()
}