package world

import (
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

type createCharacterEvent struct {
	nc   *structs.NcAvatarCreateReq
	np   *networking.Parameters
	char chan *character.Character
	err  chan error
}

type deleteCharacterEvent struct {
	nc   *structs.NcAvatarEraseReq
	np   *networking.Parameters
	done chan bool
	err  chan error
}

type characterLoginEvent struct {
	nc       *structs.NcCharLoginReq
	np       *networking.Parameters
	zoneInfo chan *structs.NcCharLoginAck
	err      chan error
}

type characterSettingsEvent struct {
	char *character.Character
	np   *networking.Parameters
	err      chan error
}

func (e createCharacterEvent) erroneous() <-chan error {
	return e.err
}

func (e deleteCharacterEvent) erroneous() <-chan error {
	return e.err
}

func (e characterLoginEvent) erroneous() <-chan error {
	return e.err
}

func (e characterSettingsEvent) erroneous() <-chan error {
	return e.err
}
