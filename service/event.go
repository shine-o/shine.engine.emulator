package service

// all events are something that either the player triggers or it should be broadcast to nearby players or mobs
// all processes can define event structures with more channels on which to receive data
type event interface {
	// notify the caller about an error while processing event
	// the process triggering the event should handle next steps in case of error
	erroneous() <- chan error
}