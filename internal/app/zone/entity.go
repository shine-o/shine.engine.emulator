package zone

import (
	"fmt"
)

const (
	lengthX = 250
	lengthY = 250
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
	handle   uint16
	fallback location
	location
	events
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

func entityInRange(e1, e2 baseEntity) bool {
	viewerX := (e1.x * 8) / 50
	viewerY := (e1.y * 8) / 50

	targetX := (e2.x * 8) / 50
	targetY := (e2.y * 8) / 50

	vertical := targetY <= viewerY+lengthY && targetY >= viewerY || targetY >= (viewerY-lengthY) && targetY <= viewerY
	horizontal := targetX <= (viewerX+lengthX) && targetX >= viewerX || targetX >= (viewerX-lengthX) && targetX <= viewerX

	if vertical && horizontal {
		return true
	}

	return false
}

type mover struct {
	baseEntity
}

type npc struct{}
