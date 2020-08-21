package service

// notify running go routines to do something periodically by sending a signal to a channel
// each tick type will trigger some action concerning a running map, an active mob or a player
type tick interface {
	tickType() uint32
}
