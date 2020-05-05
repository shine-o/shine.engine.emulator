package service

type event interface {
	// all events are something that either the player triggers or it should be broadcast to nearby players or mobs
	eventType() uint32
}
