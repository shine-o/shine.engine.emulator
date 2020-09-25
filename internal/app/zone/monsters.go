package zone

import (
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"sync"
)

const monsterHandleMin uint16 = 17000
const monsterHandleMax uint16 = 27000
const monsterAttemptsMax uint16 = 1500

type monsters struct {
	handler
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
	delete(m.handler.usedHandles, h)
	m.Unlock()
}

func (m *monsters) add(ap *monster) {
	m.Lock()
	m.active[ap.handle] = ap
	m.handler.usedHandles[ap.handle] = true
	m.Unlock()
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


// for every movement a player makes, launch a routine that:
//		iterates over every monster

// a monster can have many routines linked to it

// all these routines should be started when a monster spawns
// all these routines should be stopped when a monster dies

// when a monster dies, it should respawn again in a random number of seconds between RegMin and RegMax
// a monster should spawn at the defined random coordinates product of  X,Y, Width, Height
