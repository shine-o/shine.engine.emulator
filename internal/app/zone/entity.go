package zone

import (
	"fmt"
)

const (
	lengthX = 1200
	lengthY = 1200
)

type entity interface {
	basicActions
	getHandle() uint16
	getLocation() (uint32, uint32)
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

	if minY > 2048 {
		//log.Infof("minY=%v  maxY=%v minX=%v maxX=%v; viewer at X=%v, Y=%v ;target at X=%v Y=%v", minY, maxY, minX, maxX, viewerX, viewerY, targetX, targetY)
	}

	vertical   := (targetY <= maxY && targetY >= viewerY) || (targetY >= minY && targetY <= viewerY)
	horizontal := (targetX <= maxX && targetX >= viewerX) || (targetX >= minX && targetX <= viewerX)

	if vertical && horizontal {
		return true
	}

	return false
}

type mover struct {
	baseEntity
}

type npc struct{}
