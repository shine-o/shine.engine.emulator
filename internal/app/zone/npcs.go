package zone

import (
	"fmt"
	"sync"
)

const npcHandleMin uint16 = 17000
const npcHandleMax uint16 = 27000
const npcAttemptsMax uint16 = 1500

type npcs struct {
	active      map[uint16]*npc
	sync.RWMutex
}

type handler struct {
	handleIndex uint16
	usedHandles map[uint16]bool
	sync.RWMutex
}

func (h * handler) new(min, max, attempts uint16) (uint16, error) {
	h.RLock()
	index := h.handleIndex
	h.RUnlock()

	for attempts != 0 {

		index++

		if index == max {
			index = min
		}

		h.Lock()
		h.handleIndex = index
		h.Unlock()

		h.RLock()
		_, used := h.usedHandles[index]
		h.RUnlock()

		attempts--

		if used {
			continue
		}

		return index, nil
	}

	return 0, fmt.Errorf("\nmaximum number of attempts reached, no handle is available")

}