package zone

import (
	"fmt"
	"sync"
)

const (
	maxAttempts = 1500
)

type handler struct {
	handleIndex uint16
	usedHandles map[uint16]bool
	sync.RWMutex
}

func (h *handler) remove(hid uint16) {
	h.Lock()
	delete(h.usedHandles, hid)
	h.Unlock()
}

func (h *handler) add(ap *npc) {
	handle := ap.getHandle()
	h.Lock()
	h.usedHandles[handle] = true
	h.Unlock()
}

func (h *handler) new() (uint16, error) {
	h.RLock()
	index := h.handleIndex
	h.RUnlock()
	attempts := maxAttempts
	for attempts != 0 {

		index++
		h.RLock()
		_, used := h.usedHandles[index]
		h.RUnlock()

		attempts--

		if used {
			continue
		}

		h.Lock()
		h.handleIndex = index
		h.Unlock()

		return index, nil
	}

	return 0, fmt.Errorf("\nmaximum number of attempts reached, no handle is available")
}
