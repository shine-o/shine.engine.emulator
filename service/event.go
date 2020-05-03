package service

// all events that can happen inside a sector
// EntityMoved
type event interface {
	// all events are something that either the player triggers or it should be broadcast to nearby players or mobs
	eventType() uint32
}
