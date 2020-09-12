package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// set player data, send it to the client and return the player through the channel
type playerDataEvent struct {
	player     chan *player
	playerName string
	err        chan error
	net        *networking.Parameters
}

type playerAppearedEvent struct {
	handle uint16
	err    chan error
}

type playerDisappearedEvent struct {
	handle uint16
}

type playerWalksEvent struct {
	handle uint16
	nc     *structs.NcActMoveRunCmd
}

type playerRunsEvent struct {
	handle uint16
	nc     *structs.NcActMoveRunCmd
}

type playerJumpedEvent struct {
	handle uint16
}

type playerStoppedEvent struct {
	handle uint16
	nc     *structs.NcActStopReq
}
