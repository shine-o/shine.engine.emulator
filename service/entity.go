package service

import "fmt"

type entity interface {
	getHandle() uint16
	getLocation() (int, int)
	inbox() chan<- event
}

type basicActions interface {
	move(x, y int) error
}

type baseEntity struct {
	handle uint16
	location struct {
		x,y int
	}
	events chan <- event
}

func (b baseEntity) getHandle() uint16  {
	return b.handle
}

func (b baseEntity) getLocation() (int, int){
	return b.location.x, b.location.y
}

func (b baseEntity) inbox() chan<- event {
	return b.events
}

func (b baseEntity) move(m zoneMap, x, y int) error {
	if canWalk(m.walkableX, m.walkableY, x, y) {
		return nil
	}
	return fmt.Errorf("entity %v cannot move to x %v  y %v", b.getHandle(), x, y)
}

type player struct {
	baseEntity
	conn playerConnection
}

type playerConnection struct {
	close chan <- bool
	data chan <- []byte
}

type mover struct{
	baseEntity
}

type monster struct {
	baseEntity
}

type npc struct {}


