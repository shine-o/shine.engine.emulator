package zone

import (
	"sync"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
)

// all events are something that either the player triggers or it should be broadcast to nearby players or mobs
// all processes can define event structures with more channels on which to receive data
// the reason for events and workers is to define access points for data.
// a worker is typically a method which has access to data (map, mobs, players)
type event interface { // notify the caller about an error while processing event
	// the process triggering the event should handle next steps in case of error
	// erroneous() <-chan error
}

type eventIndex int

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
type dynamicEvents struct {
	sync.RWMutex
	events map[string]events
}

// to use when no particular data is needed
type emptyEvent struct{}

const (
	// zone events
	playerMapLogin eventIndex = iota
	playerSHN
	playerData
	heartbeatUpdate
	playerLogoutStart
	playerLogoutCancel
	playerLogoutConclude
	persistPlayerPosition
	changeMap

	// map events
	playerHandleMaintenance
	playerHandle
	playerAppeared
	playerDisappeared
	playerWalks
	playerRuns
	playerJumped
	playerStopped
	unknownHandle
	monsterAppeared
	monsterDisappeared
	monsterWalks
	monsterRuns
	monsterStopped
	monsterIdle
	queryPlayer
	queryMonster
	playerSelectsEntity
	playerUnselectsEntity
	playerClicksOnNpc
	playerPromptReply
	itemIsMoved
	itemEquip
	itemUnEquip

	// entity data update events
	eduPosition
	eduState
	eduStats
	eduEquipItem
	eduUnEquipItem
	eduSelectEntity
	eduUnselectsEntity

	// dynamically registered events
	// events that are defined at run time
	dLogoutCancel
	dLogoutConclude
)

var (
	zEvents = []eventIndex{
		playerMapLogin, playerSHN, playerData,
		heartbeatUpdate, playerLogoutStart,
		playerLogoutCancel, playerLogoutConclude, persistPlayerPosition,
		changeMap,
	}
	mapEvents = []eventIndex{
		playerHandle,
		playerHandleMaintenance,
		queryPlayer, queryMonster,
		playerAppeared, playerDisappeared, playerJumped, playerWalks, playerRuns, playerStopped,
		unknownHandle, monsterAppeared, monsterDisappeared, monsterWalks, monsterRuns,
		playerSelectsEntity, playerUnselectsEntity, playerClicksOnNpc, playerPromptReply, itemIsMoved, itemEquip, itemUnEquip,
	}

	playerEvents = []eventIndex{
		eduPosition, eduState, eduStats, eduEquipItem, eduUnEquipItem, eduSelectEntity, eduUnselectsEntity,
	}
)

func (de *dynamicEvents) add(sid string, i eventIndex) chan event {
	de.Lock()
	_, ok := de.events[sid]
	if !ok {
		de.events[sid] = events{
			send: make(sendEvents),
			recv: make(recvEvents),
		}
	}
	c := make(chan event)
	de.events[sid].send[i] = c
	de.events[sid].recv[i] = c
	de.Unlock()
	return c
}

func eventTypeCastError(expected, got string) error {
	return errors.Err{
		Code: errors.ZoneUnexpectedEventType,
		Details: errors.Details{
			"expected": expected,
			"got":      got,
		},
	}
}
