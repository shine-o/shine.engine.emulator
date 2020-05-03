package service

import (
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

type playerEventError struct {
	code int
	message string
}

type playerAppearedEvent struct {
	nc structs.NcBriefInfoLoginCharacterCmd
}

const (
	playerAppeared uint32 = iota
	playerDisappeared
	playerMoved
	playerStopped
	playerJumped
)

var errMissingZoneParams = playerEventError{
	code:    0,
	message: "missing zoneParameters struct",
}

var errInvalidMap = playerEventError{
	code:    0,
	message: "character is located in an map that is not running on this zone",
}

func (pae playerAppearedEvent) eventType() uint32 {
	return playerAppeared
}

func (pae playerAppearedEvent) process(np *networking.Parameters, char character.Character) error {
	zp, ok := np.Extra.(zoneParameters)
	if ok {
		m, ok := zp.rm[int(char.Location.MapID)]
		if ok {
			entity := &player{
				baseEntity: baseEntity{
					handle: m.handles.newHandle(),
					location: struct {
						x, y int
					}{
						x: char.Location.X,
						y: char.Location.Y,
					},
					events: make(chan event),
				},
			}

			m.send[playerAppeared] <- playerAppearedEvent{
				nc: structs.NcBriefInfoLoginCharacterCmd{
					Handle: entity.getHandle(),
					CharID: structs.Name5{
						Name: char.Name,
					},
					Coordinates:     structs.ShineCoordType{},
					Mode:            0,
					Class:           char.Appearance.Class,
					Shape:           char.Appearance.NcRepresentation(),
					ShapeData:       structs.NcBriefInfoLoginCharacterCmdShapeData{},
					Polymorph:       0,
					Emoticon:        structs.StopEmoticonDescript{},
					CharTitle:       structs.CharTitleBriefInfo{},
					AbstateBit:      structs.AbstateBit{},
					MyGuild:         0,
					Type:            0,
					IsAcademyMember: 0,
					IsAutoPick:      0,
					Level:           char.Attributes.Level,
					Animation:       [32]byte{},
					MoverHandle:     0,
					MoverSlot:       0,
					KQTeamType:      0,
					UsingMinipet:    0,
					Unk:             0,
				},
			}
		} else {
			return errInvalidMap
		}
	} else {
		return errMissingZoneParams
	}
	return nil
}

func (e playerEventError) Error() string {
	return e.message
}
