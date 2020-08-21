package service

import "fmt"

type entity interface {
	getHandle() uint16
	getLocation() (uint32, uint32)
}

type basicActions interface {
	move(x, y int) error
}

type location struct {
	mapID   int
	mapName string
	x, y    uint32
	d       uint8
	movements [15]movement
}

type movement struct {
	x, y uint32
}

type baseEntity struct {
	handle uint16
	location
	events
}

func (b baseEntity) getHandle() uint16 {
	return b.handle
}

func (b baseEntity) getLocation() (uint32, uint32) {
	return b.location.x, b.location.y
}

func (b baseEntity) move(m * zoneMap, x, y int) error {
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