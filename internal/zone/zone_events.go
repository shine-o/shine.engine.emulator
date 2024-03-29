package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

type playerSHNEvent struct {
	inboundNC structs.NcMapLoginReq
	ok        chan bool
	err       chan error
}

// set player data, send it to the client and return the player through the channel
type playerDataEvent struct {
	player     chan *player
	playerName string
	err        chan error
	net        *networking.Parameters
}

type playerMapLoginEvent struct {
	nc structs.NcMapLoginReq
	np *networking.Parameters
}

type heartbeatUpdateEvent struct {
	*session
}

type playerLogoutStartEvent struct {
	sessionID string
	mapID     int
	handle    uint16
	err       chan error
}

type playerLogoutCancelEvent struct {
	sessionID string
	err       chan error
}

type playerLogoutConcludeEvent struct {
	sessionID string
	err       chan error
}

type persistPlayerPositionEvent struct {
	p *player
}

type changeMapEvent struct {
	p    *player
	s    *session
	next *zoneMap
	prev *zoneMap
}
