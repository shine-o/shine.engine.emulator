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
	playerHandle uint16
	mapID 		 int
	err        	 chan error
}

func (e *playerDataEvent) erroneous() <-chan error {
	return e.err
}

func (e *playerAppearedEvent) erroneous() <-chan error {
	return e.err
}
