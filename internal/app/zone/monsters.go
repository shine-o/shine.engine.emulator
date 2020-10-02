package zone

import (
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"sync"
)

type monsters struct {
	*handler
	handleIndex uint16
	active      map[uint16]*monster
	sync.RWMutex
}

func (m *monsters) all() <-chan *monster {

	m.RLock()
	ch := make(chan *monster, len(m.active))
	m.RUnlock()

	go func(send chan<- *monster) {
		m.RLock()
		for _, ap := range m.active {
			send <- ap
		}
		m.RUnlock()
		close(send)
	}(ch)

	return ch
}

func (m *monsters) get(h uint16) *monster {
	m.RLock()
	monster := m.active[h]
	m.RUnlock()
	return monster
}

func (m *monsters) remove(h uint16) {
	m.Lock()
	delete(m.active, h)
	m.Unlock()
}

func (m *monsters) add(ap *monster) {
	m.Lock()
	m.active[ap.handle] = ap
	m.handler.usedHandles[ap.handle] = true
	m.Unlock()
}

func (m *monster) ncBriefInfoRegenMobCmd() structs.NcBriefInfoRegenMobCmd {
	m.RLock()
	nc := structs.NcBriefInfoRegenMobCmd{
		Handle: m.handle,
		Mode:   byte(m.mobInfoServer.EnemyDetect),
		MobID:  m.mobInfo.ID,
		Coord: structs.ShineCoordType{
			XY: structs.ShineXYType{
				X: uint32(m.current.x),
				Y: uint32(m.current.y),
			},
			Direction: uint8(m.current.d),
		},
	}
	m.RUnlock()
	return nc
}

func adjacentMonstersInform(p *player, zm *zoneMap) {
	for m := range zm.entities.monsters.all() {
		go func(p *player, m *monster) {
			if !knownMonster(p, m.getHandle()) {
				if monsterInRange(p, m) {
					nc := m.ncBriefInfoRegenMobCmd()
					ncBriefInfoRegenMobCmd(p, &nc)
				}
			}
		}(p, m)
	}
}

func knownMonster(p *player, mh uint16) bool {
	p.RLock()
	_, ok := p.monsters[mh]
	p.RUnlock()
	if ok {
		return true
	}
	return false
}

func monsterInRange(p *player, m *monster) bool {
	p.RLock()
	m.RLock()
	yes := entityInRange(p.baseEntity, m.baseEntity)
	p.RUnlock()
	m.RUnlock()

	if yes {
		p.Lock()
		p.monsters[m.handle] = m
		p.Unlock()
		return true
	}
	return false
}
