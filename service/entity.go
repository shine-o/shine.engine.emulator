package service

type entity interface {
	handle() int
	location() (int, int)
	inbox() chan <- event
}

type player interface {
	entity
}

type npc interface {
	entity
}

type monster interface {
	entity
}

type basicActions interface {
	move(x,y int) error
	attack(int)
}