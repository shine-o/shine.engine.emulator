package service

import (
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

type playerEventError struct {
	code    int
	message string
}

func (e playerEventError) Error() string {
	return e.message
}

// set player data, send it to the client and return the player through the channel
type playerDataEvent struct {
	player     chan *player
	playerName string
	err        chan error
	net        *networking.Parameters
}

func (e *playerDataEvent) erroneous() <-chan error {
	return e.err
}

type playerAppearedEvent struct {
	np         *networking.Parameters
	char       *character.Character
	player     *player
	outboundNC structs.NcBriefInfoLoginCharacterCmd
	err        chan error
}

func (e *playerAppearedEvent) erroneous() <-chan error {
	return e.err
}

var errInvalidMap = playerEventError{
	code:    0,
	message: "character is located in an map that is not running on this zone",
}
