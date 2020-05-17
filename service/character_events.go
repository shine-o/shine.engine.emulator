package service

import (
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/shine-o/shine.engine.world/service/character"
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
	char     chan *character.Character
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