package zone

import (
	"fmt"
	"sync"
)

const playerHandleMin uint16 = 8000
const playerHandleMax uint16 = 12000
const playerAttemptsMax uint16 = 500

type players struct {
	handleIndex uint16
	active      map[uint16]*player
	sync.RWMutex
}

func (p *players) activePlayers() <-chan *player {

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

func (p *players) removeHandle(h uint16) {
	p.Lock()
	delete(p.active, h)
	p.Unlock()
}

func (p *players) addHandle(h uint16, ap *player) {
	p.Lock()
	p.active[h] = ap
	p.Unlock()
}

func (p *players) newHandle() (uint16, error) {
	var attempts uint16 = 0
	min := playerHandleMin
	max := playerHandleMax
	maxAttempts := playerAttemptsMax

	p.RLock()
	index := p.handleIndex
	p.RUnlock()

	for {
		if attempts == maxAttempts {
			return 0, fmt.Errorf("\nmaximum number of attempts reached, no handle is available")
		}

		index++

		if index == max {
			index = min
		}

		p.Lock()
		p.handleIndex = index
		p.Unlock()

		p.RLock()
		if _, used := p.active[index]; used {
			attempts++
			continue
		}
		p.RUnlock()

		return index, nil
	}
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
