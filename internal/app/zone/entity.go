package zone

import (
	"fmt"
	"sync"
)

type entity interface {
	getHandle() uint16
	getLocation() (uint32, uint32)
}

type basicActions interface {
	move(x, y int) error
}

type location struct {
	mapID     int
	mapName   string
	x, y      uint32
	d         uint8
	movements [15]movement
}

type movement struct {
	x, y uint32
}

type baseEntity struct {
	handle uint16
	location
	events
	sync.Mutex
}

const (
	lengthX = 100
	lengthY = 100
)

func playerInRange(viewer, target *player) bool {

	targetX := (target.x * 8) / 50
	targetY := (target.y * 8) / 50

	viewerX := (viewer.x * 8) / 50
	viewerY := (viewer.y * 8) / 50

	vertical := targetY <= viewerY+lengthY && targetY > viewerY || targetY >= (viewerY-lengthY) && targetY < viewerY
	horizontal := targetX <= (viewerX+lengthX) && targetX > viewerX || targetX >= (viewerX-lengthX) && targetX < viewerX

	if vertical && horizontal {

		viewer.Lock()
		viewer.knownNearbyPlayers[target.handle] = target
		viewer.Unlock()

		log.Infof("%v is in range of %v", target.handle, viewer.handle)
		return true
	}

	log.Infof("%v is in not in range of %v", target.handle, viewer.handle)

	return false
}

func (b *baseEntity) getHandle() uint16 {
	return b.handle
}

func (b *baseEntity) getLocation() (uint32, uint32) {
	return b.location.x, b.location.y
}

func (b *baseEntity) move(m *zoneMap, x, y uint32) error {
	if canWalk(m.walkableX, m.walkableY, x, y) {
		return nil
	}
	return fmt.Errorf("entity %v cannot move to x %v  y %v", b.getHandle(), x, y)
}

type mover struct {
	baseEntity
}

type monster struct {
	baseEntity
}

type npc struct{}
