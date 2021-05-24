package login

import "sync"

// all events are something that either the player triggers or it should be broadcast to nearby players or mobs
// all processes can define event structures with more channels on which to receive data
// the reason for events and workers is to define access points for data.
// a worker is typically a method which has access to data (map, mobs, players)
type event interface { // notify the caller about an error while processing event
	// the process triggering the event should handle next steps in case of error
}

type eventIndex uint32

type sendEvents map[eventIndex]chan<- event

type recvEvents map[eventIndex]<-chan event

type events struct {
	send sendEvents
	recv recvEvents
}

// dynamically add channels which are linked to a session UUID
// this will allow to launch routines which can be revisited from another event,
// e.g playerLogoutStart event starts a routine that will automatically close the connection in 10 seconds
// but the client can send a cancel signal, therefore we need to notify the routine to abort the countdown
type dynamic struct {
	sync.RWMutex
	events map[string]events
}

func (d *dynamic) add(sid string, i eventIndex) chan event {
	d.Lock()
	_, ok := d.events[sid]
	if !ok {
		d.events[sid] = events{
			send: make(sendEvents),
			recv: make(recvEvents),
		}
	}
	c := make(chan event)
	d.events[sid].send[i] = c
	d.events[sid].recv[i] = c
	d.Unlock()
	return c
}

// to use when no particular data is needed
type emptyEvent struct {
	err chan error
}

// todo: separate with different iotas, for now its simpler to have it like this, but in the future we'll have hundreds of events
const (
	clientVersion eventIndex = iota
	credentialsLogin
	credentialsOk
	worldManagerStatus
	serverList
	serverSelect
	tokenLogin
)
