package service

// all events are something that either the player triggers or it should be broadcast to nearby players or mobs
// all processes can define event structures with more channels on which to receive data
type event interface {
	// notify the caller about an error while processing event
	// the process triggering the event should handle next steps in case of error
	erroneous() <-chan error
}

type eventIndex uint32

type sendEvents map[eventIndex]chan<- event

type recvEvents map[eventIndex]<-chan event

// to use when no particular data is needed
type emptyEvent struct {
	err chan error
}

const (
	loadPlayerData eventIndex = iota
	registerPlayerHandle
	playerAppeared
	playerDisappeared
	playerMoved
	playerStopped
	playerJumped

	queryMap
	queryPlayer
	queryMonster

	clientSHN

	handleCleanUp

	heartbeatMissing
)

func (e * emptyEvent) erroneous() <-chan error  {
	return e.err
}
