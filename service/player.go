package service

import (
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/structs"
)

type player struct {
	baseEntity
	conn        playerConnection
	stats playerStats
	state playerState
}

type playerConnection struct {
	close chan<- bool
	outboundData  chan<- []byte
}



type playerState struct {
	level uint8
	autoPickup bool
	class uint8
	gender uint8
	hairType uint8
	hairColour uint8
	faceType uint8
	polymorph uint16
	moverHandle uint16
	moverSlot uint8
	miniPet bool
}

func (p *player) ncLoginRepresentation(char *character.Character) structs.NcBriefInfoLoginCharacterCmd {
	nc := structs.NcBriefInfoLoginCharacterCmd{
		Handle: p.getHandle(),
		CharID: structs.Name5{
			Name: char.Name,
		},
		Coordinates: structs.ShineCoordType{
			XY: structs.ShineXYType{
				X: char.Location.X,
				Y: char.Location.Y,
			},
			Direction: char.Location.D,
		},
		Mode:            0,
		Class:           char.Appearance.Class,
		Shape:           char.Appearance.NcRepresentation(),
		ShapeData:       structs.NcBriefInfoLoginCharacterCmdShapeData{},
		Polymorph:       65535,
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
	}
	return nc
}
