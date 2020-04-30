package service

import (
	"github.com/shine-o/shine.engine.core/game/entities"
)

// all events that can happen inside a sector
// EntityMoved
type event interface {
	// all events are something that either the player triggers or it should be broadcast to nearby players or mobs
	// in all cases, a network command is needed to notify the players
	networkCommand() interface{}
}

type entityAppeared struct {
	event
	entities.Mob
}

type entityDisappeared struct {
	event
	entities.Mob
}

type entityMoved struct {
	event
	entities.Mob
}

type entityStopped struct {
	event
	entities.Mob
}

type entityJumped struct {
	event
	entities.Mob
}
