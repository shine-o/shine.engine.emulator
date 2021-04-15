package zone

import (
	"fmt"
	"sync"
)

const (
	//lengthX = 512
	//lengthY = 512
	//lengthX = 256
	//lengthY = 256
	lengthX     = 225
	lengthY     = 225
	maxAttempts = 1500
)

type entity interface {
	basicActions
	getHandle() uint16
	getLocation() (int, int)
}

type handler struct {
	handleIndex uint16
	usedHandles map[uint16]bool
	sync.RWMutex
}

type basicActions interface {
	move(m *zoneMap, x, y int) error
}

type location struct {
	mapID     int
	mapName   string
	x, y, d   int
	movements []movement
}

type movement struct {
	x, y uint32
}

type baseEntity struct {
	info     entityInfo
	fallback location
	current  location
	next     location
	events   events
	// dangerZone: only to be used when loading or other situation!!
	sync.RWMutex
}

type entityInfo struct {
	handle uint16
	monster bool
}

type targeting struct {
	selectionOrder byte
	selectingP     *player
	selectingN     *npc
	selectedByP    []*player
	selectedByN    []*npc
	sync.RWMutex
}

type entityState struct {
	idling   chan bool
	fighting chan bool
	chasing  chan bool
	fleeing  chan bool
}

type mover struct {
	baseEntity
}

var _ entity = (*player)(nil)

var _ entity = (*npc)(nil)

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

func (b *baseEntity) getHandle() uint16 {
	b.RLock()
	h := b.info.handle
	b.RUnlock()
	return h
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

func entityInRange(e1, e2 location) bool {
	viewerX, viewerY := igCoordToBitmap(e1.x, e1.y)
	targetX, targetY := igCoordToBitmap(e2.x, e2.y)

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
