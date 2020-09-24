package zone

import (
	"sync"
)

const playerHandleMin uint16 = 8000
const playerHandleMax uint16 = 12000
const playerAttemptsMax uint16 = 500

type players struct {
	handler
	active      map[uint16]*player
	sync.RWMutex
}

func (p *players) all() <-chan *player {

	p.RLock()
	ch := make(chan *player, len(p.active))
	p.RUnlock()

	go func(send chan<- *player) {
		p.RLock()
		for _, ap := range p.active {
			send <- ap
		}
		p.RUnlock()
		close(send)
	}(ch)

	return ch
}

func (p *players) get(h uint16) *player {
	p.RLock()
	player := p.active[h]
	p.RUnlock()
	return player
}

func (p *players) remove(h uint16) {
	p.Lock()
	delete(p.active, h)
	delete(p.usedHandles, h)
	p.Unlock()
}

func (p *players) add(ap *player) {
	p.Lock()
	p.active[ap.handle] = ap
	p.usedHandles[ap.handle] = true
	p.Unlock()
}

func playerInRange(v, t *player) bool {
	v.RLock()
	t.RLock()

	yes := entityInRange(v.baseEntity, t.baseEntity)

	v.RUnlock()
	t.RUnlock()

	if yes {
		v.Lock()
		v.players[t.handle] = t
		v.Unlock()
		return true
	}
	return false
}
