package zone

import (
	"fmt"
	"sync"
)

const (
	lengthX = 256
	lengthY = 256
)

type entity interface {
	basicActions
	getHandle() uint16
	getLocation() (uint32, uint32)
}

type handler struct {
	handleIndex uint16
	usedHandles map[uint16]bool
	sync.RWMutex
}

func (h *handler) new(min, max, attempts uint16) (uint16, error) {
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

type basicActions interface {
	move(x, y int) error
}

type location struct {
	mapID     int
	mapName   string
	x, y      int
	d         int
	movements [15]movement
}

type movement struct {
	x, y uint32
}

type baseEntity struct {
	handle   uint16
	fallback location
	current  location
	events
}

type targeting struct {
	selectionOrder byte
	selectingP     *player
	selectingM     *monster
	selectingN     *npc
	selectedByP    []*player
	selectedByM    []*monster
	selectedByN    []*npc
}

type status struct {
	idling   chan bool
	fighting chan bool
	chasing  chan bool
	fleeing  chan bool
}

func (b *baseEntity) getHandle() uint16 {
	return b.handle
}

func (b *baseEntity) getLocation() (int, int) {
	return b.current.x, b.current.y
}

func (b *baseEntity) move(m *zoneMap, x, y int) error {
	if canWalk(m.walkableX, m.walkableY, x, y) {
		return nil
	}
	return fmt.Errorf("entity %v cannot move to x %v  y %v", b.getHandle(), x, y)
}

func entityInRange(e1, e2 baseEntity) bool {
	targetX, targetY := igCoordToBitmap(e2.current.x, e2.current.y)
	viewerX, viewerY := igCoordToBitmap(e1.current.x, e1.current.y)

	maxY := viewerY + lengthY
	minY := viewerY - lengthY

	maxX := viewerX + lengthX
	minX := viewerX - lengthX

	vertical := (targetY <= maxY && targetY >= viewerY) || (targetY >= minY && targetY <= viewerY)
	horizontal := (targetX <= maxX && targetX >= viewerX) || (targetX >= minX && targetX <= viewerX)

	if vertical && horizontal {
		return true
	}

	return false
}

type mover struct {
	baseEntity
}
