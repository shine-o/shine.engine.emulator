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

type playerAppearedEvent struct {
	*player
	nc structs.NcBriefInfoLoginCharacterCmd
}

const (
	playerEntersMap uint32 = iota
	playerAppeared
	playerDisappeared
	playerMoved
	playerStopped
	playerJumped
)

var errInvalidMap = playerEventError{
	code:    0,
	message: "character is located in an map that is not running on this zone",
}

func (pae *playerAppearedEvent) process(np *networking.Parameters, char *character.Character) error {
	m, ok := rm[int(char.Location.MapID)]
	if ok {
		p := &player{
			baseEntity: baseEntity{
				handle: m.handles.newHandle(),
				location: location{
					x: char.Location.X,
					y: char.Location.Y,
				},
				events: make(chan event),
			},
			conn: playerConnection{
				close: np.NetVars.CloseConnection,
				data:  np.NetVars.OutboundSegments.Send,
			},
			characterID: char.ID,
		}
		pae.player = p
		pae.nc = p.ncLoginRepresentation(char)
		m.send[playerAppeared] <- pae
	} else {
		return errInvalidMap
	}
	return nil
}

func (pae playerAppearedEvent) eventType() uint32 {
	return playerAppeared
}

func (e playerEventError) Error() string {
	return e.message
}
