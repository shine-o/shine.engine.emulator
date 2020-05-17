package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/shine-o/shine.engine.world/service/character"
)

// NcCharLoginReq handles a petition to login to the zone where the character's location map is running
// NC_CHAR_LOGIN_REQ
func NcCharLoginReq(ctx context.Context, np * networking.Parameters) {
	var (
		nc structs.NcCharLoginReq
		cl characterLoginEvent
	)

	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		log.Error(err)
		return
	}

	cl = characterLoginEvent{
		nc: &nc,
		zoneInfo: make(chan * structs.NcCharLoginAck),
		char : make(chan *character.Character),
		err: make(chan error),
	}

	worldEvents[characterLogin] <- &cl

	select {
	case nc := <- cl.zoneInfo:
		ncCharLoginAck(np, nc)
	case char := <- cl.char:
		gameOptions, err := character.NcGameOptions(char.Options.GameOptions)
		keyMap, err := character.NcKeyMap(char.Options.Keymap)
		shortcuts, err := character.NcShortcutData(char.Options.Shortcuts)

		ncCharOptionImproveGetGameOptionCmd(np, &gameOptions)
		ncCharOptionImproveGetKeymapCmd(np, &keyMap)
		ncCharOptionImproveGetShortcutDataCmd(np, &shortcuts)

	case err := <- cl.err:
		log.Error(err)
		return
	}
}

// NcCharLoginAck sends zone connection info to the client
// NC_CHAR_LOGIN_ACK
func ncCharLoginAck(np * networking.Parameters, nc *structs.NcCharLoginAck) {
	// query the zone master for connection info for the map
	// send it to the client
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4099,
		},
		NcStruct: nc,
	}
	pc.Send(np.OutboundSegments.Send)
}
