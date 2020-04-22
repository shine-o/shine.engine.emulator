package entities

type Mob interface {
	Handle() int
	Location() (int, int)
	Move(x,y int) error
}

type Player interface {
	Mob
}

type NPC interface {
	Mob
}

type Monster interface {
	Mob
}