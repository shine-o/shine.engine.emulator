package world

import "sync"

const (
	serverTime eventIndex = iota
	serverSelect
	serverSelectToken

	createCharacter
	deleteCharacter

	characterLogin

	characterSettings
	characterKeymap
	characterShortcuts

	// items that the player adds or removes from the quick access bars
	updateShortcuts
	// features that the player can activate / deactivate at Esc > Options > Game Settings
	updateGameSettings
	// keyboard keys that the player can map to game functions
	updateKeymap

	// when the player wants to return to the character select screen
	characterSelect
)

// all events are something that either the player triggers or it should be broadcast to nearby players or mobs
type event interface{}

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
