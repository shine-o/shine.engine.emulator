package zone

import (
	"fmt"
	"sync"
)

const monsterHandleMin uint16 = 17000
const monsterHandleMax uint16 = 27000
const monsterAttemptsMax uint16 = 50

type monstersData struct {
	mobInfo interface{}
	mobInfoServer interface{}
}

type monsters struct {
	handleIndex uint16
	active      map[uint16]*monster
	sync.RWMutex
}

// a monster can have many routines linked to it

// all these routines should be started when a monster spawns
// all these routines should be stopped when a monster dies

// when a monster dies, it should respawn again in a random number of seconds between RegMin and RegMax
// a monster should spawn at the defined random coordinates product of  X,Y, Width, Height


func (m *monsters) newHandle() (uint16, error) {
	var attempts uint16 = 0
	min := monsterHandleMin
	max := monsterHandleMax
	maxAttempts := monsterAttemptsMax

	index := m.handleIndex

	for {

		if attempts == maxAttempts {
			return 0, fmt.Errorf("\nmaximum number of attempts reached, no handle is available")
		}

		index++

		if index == max {
			index = min
		}

		m.handleIndex = index

		if _, used := m.active[index]; used {
			attempts++
			continue
		}

		return index, nil
	}
}
