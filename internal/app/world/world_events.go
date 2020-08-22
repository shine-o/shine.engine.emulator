package world

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

type serverTimeEvent struct {
	np  *networking.Parameters
	err chan error
}

type serverSelectEvent struct {
	nc  *structs.NcUserLoginWorldReq
	np  *networking.Parameters
	err chan error
}

type serverSelectTokenEvent struct {
	np  *networking.Parameters
	err chan error
}

func (s serverTimeEvent) erroneous() <-chan error {
	return s.err
}

func (s serverSelectEvent) erroneous() <-chan error {
	return s.err
}

func (s serverSelectTokenEvent) erroneous() <-chan error {
	return s.err
}
