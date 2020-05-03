package service

import (
	"github.com/shine-o/shine.engine.core/structs"
)

// all events that can happen inside a sector
// EntityMoved
type event interface {
	// all events are something that either the player triggers or it should be broadcast to nearby players or mobs
	eventType() uint32
}

type playerAppearedEvent struct {
	nc structs.NcBriefInfoLoginCharacterCmd
}

func (pae playerAppearedEvent) eventType() uint32 {
	return playerAppeared
}

const (
	playerAppeared uint32 = iota
	playerDisappeared
	playerMoved
	playerStopped
	playerJumped
)