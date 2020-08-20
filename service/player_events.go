package service

import (
	"github.com/shine-o/shine.engine.core/networking"
)

// set player data, send it to the client and return the player through the channel
type playerDataEvent struct {
	player     chan *player
	playerName string
	err        chan error
	net        *networking.Parameters
}

type playerAppearedEvent struct {
	handle 		uint16
	err        	 chan error
}

type playerDisappearedEvent struct {
	handle uint16
	err        	 chan error
}
