package service

import (
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

type serverTimeEvent struct {
	np *networking.Parameters
	err chan error
}

type serverSelectEvent struct {
	nc *structs.NcUserLoginWorldReq
	np *networking.Parameters
	err chan error
}

type serverSelectTokenEvent struct {
	np *networking.Parameters
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

