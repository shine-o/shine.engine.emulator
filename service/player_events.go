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

type playerEntersMapEvent struct {
	session *session
	netParams networking.Network
	inboundNC structs.NcMapLoginReq
	err chan error
}

type playerAppearedEvent struct {
	np *networking.Parameters
	char *character.Character
	player *player
	outboundNC structs.NcBriefInfoLoginCharacterCmd
	err chan error
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

func (pae *playerAppearedEvent) erroneous() <- chan error {
	return pae.err
}

func (pem *playerEntersMapEvent) erroneous() <- chan error {
	return pem.err
}
//func (pae *playerAppearedEvent) process() error {
//	return nil
//}
//
//func (pem playerEntersMapEvent) process() error {
//	// given the character, create handle for it in the map he's located at and register it
//	// we save the map ID and handle ID in the session
//	char, err := character.GetByName(db, pem.inboundNC.CharData.CharID.Name)
//	if err != nil {
//		return err
//	}
//
//	m, ok := rm[int(char.Location.MapID)]
//	if ok {
//		//p := &player{
//		//	baseEntity: baseEntity{
//		//		handle: m.handles.newHandle(),
//		//		location: location{
//		//			x: char.Location.X,
//		//			y: char.Location.Y,
//		//		},
//		//		events: make(chan event),
//		//	},
//		//	conn: playerConnection{
//		//		close: pem.netParams.CloseConnection,
//		//		outboundData:  pem.netParams.OutboundSegments.Send,
//		//	},
//		//	characterID: char.ID,
//		//}
//
//	} else {
//		return errInvalidMap
//	}
//}

func (e playerEventError) Error() string {
	return e.message
}
