package zone

import (
	"fmt"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"sync"
)

const monsterHandleMin uint16 = 17000
const monsterHandleMax uint16 = 27000
const monsterAttemptsMax uint16 = 1500

type monsters struct {
	handleIndex uint16
	active      map[uint16]*monster
	sync.RWMutex
}

func (m * monsters) activeMonsters() <- chan *monster {

	m.RLock()
	ch := make(chan *monster, len(m.active))
	m.RUnlock()

	go func(send chan <- *monster) {
		m.RLock()
		for _, ap := range m.active {
			send <- ap
		}
		m.RUnlock()
		close(send)
	}(ch)

	return ch
}

func (m * monsters) removeHandle(h uint16)  {
	m.Lock()
	delete(m.active, h)
	m.Unlock()
}

func (m * monsters) addHandle(h uint16, ap * monster)  {
	m.Lock()
	m.active[h] = ap
	m.Unlock()
}

func (m *monsters) newHandle() (uint16, error) {
	var attempts uint16 = 0
	min := monsterHandleMin
	max := monsterHandleMax
	maxAttempts := monsterAttemptsMax

	m.RLock()
	index := m.handleIndex
	m.RUnlock()

	for {

		if attempts == maxAttempts {
			return 0, fmt.Errorf("\nmaximum number of attempts reached, no handle is available")
		}

		index++

		if index == max {
			index = min
		}

		m.Lock()
		m.handleIndex = index
		m.Unlock()

		m.RLock()
		_, used := m.active[index]
		m.RUnlock()

		if used {
			attempts++
			continue
		}

		return index, nil
	}
}

func knownMonster(p * player, mh uint16) bool {
	p.RLock()
	_, ok := p.monsters[mh]
	p.RUnlock()
	if ok {
		return true
	}
	return false
}

func adjacentMonstersInform(p * player, zm * zoneMap)  {
	for m := range zm.entities.activeMonsters() {
		go func(p * player, m * monster) {
			if !knownMonster(p, m.getHandle()) {
				if monsterInRange(p, m) {
					nc := m.ncBriefInfoRegenMobCmd()
					ncBriefInfoRegenMobCmd(p, &nc)
				}
			}

		}(p, m)
	}
}

func (m * monster) ncBriefInfoRegenMobCmd() structs.NcBriefInfoRegenMobCmd {
	m.RLock()
	nc := structs.NcBriefInfoRegenMobCmd{
		Handle:         m.handle,
		Mode:           byte(m.mobInfoServer.EnemyDetect),
		MobID:          m.mobInfo.ID,
		Coord:          structs.ShineCoordType{
			XY:        structs.ShineXYType{
				X: m.x,
				Y: m.y,
			},
			Direction: m.d,
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

// for every movement a player makes, launch a routine that:
//		iterates over every monster

// a monster can have many routines linked to it

// all these routines should be started when a monster spawns
// all these routines should be stopped when a monster dies

// when a monster dies, it should respawn again in a random number of seconds between RegMin and RegMax
// a monster should spawn at the defined random coordinates product of  X,Y, Width, Height
