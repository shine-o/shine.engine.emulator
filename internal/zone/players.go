package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"sync"
	"time"
)

type players struct {
	*handler
	active map[uint16]*player
	*sync.RWMutex
}

func (p *players) all() <-chan *player {

	p.RLock()
	ch := make(chan *player, len(p.active))
	p.RUnlock()

	go func(p *players, send chan<- *player) {
		p.RLock()
		for _, ap := range p.active {
			send <- ap
		}
		p.RUnlock()
		close(send)
	}(p, ch)

	return ch
}

func (p *players) get(h uint16) *player {
	p.RLock()
	defer 	p.RUnlock()
	player, ok  := p.active[h]
	if !ok {
		log.Error(errors.Err{
			Code:    errors.ZoneMissingPlayer,
			Details: errors.ErrDetails{
				"handle": h,
			},
		})
		return nil
	}
	return player
}

func (p *players) remove(h uint16) {
	p.Lock()
	delete(p.active, h)
	p.Unlock()
}

func (p *players) add(ap *player) {
	h := ap.getHandle()

	ap.state.Lock()
	ap.state.justSpawned = true
	ap.state.Unlock()

	p.Lock()
	p.active[h] = ap
	p.handler.usedHandles[h] = true
	p.Unlock()

	go func(p *player) {
		time.Sleep(15 * time.Second)
		p.state.Lock()
		p.state.justSpawned = false
		p.state.Unlock()
	}(ap)
}

func playerInRange(v, t *player) bool {
	h := t.getHandle()

	v.baseEntity.RLock()
	t.baseEntity.RLock()
	vc := v.baseEntity.current
	tc := t.baseEntity.current
	v.baseEntity.RUnlock()
	t.baseEntity.RUnlock()

	yes := entityInRange(vc, tc)

	if yes {
		v.proximity.Lock()
		v.proximity.players[h] = t
		v.proximity.Unlock()
		return true
	}
	return false
}
