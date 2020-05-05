package entities

type Mob interface {
	Handle() int
	Location() (int, int)
}

type Player interface {
	Mob
}

type NPC interface {
	Mob
}

type Monster interface {
	Mob
	Inbox()
}

type BasicActions interface {
	Move(x,y int) error
	Attack(int)
}