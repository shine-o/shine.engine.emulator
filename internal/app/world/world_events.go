package world

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

type serverTimeEvent struct {
	np *networking.Parameters
}

type serverSelectEvent struct {
	nc *structs.NcUserLoginWorldReq
	np *networking.Parameters
}

type serverSelectTokenEvent struct {
	np *networking.Parameters
}
