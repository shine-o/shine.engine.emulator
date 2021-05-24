package world

import (
	"context"
	"encoding/hex"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

// NcCharOptionGetShortcutSizeReq
// NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ
func ncCharOptionGetShortcutSizeReq(ctx context.Context, np *networking.Parameters) {
	ncCharOptionGetShortcutSizeAck(np)
}

// NcCharOptionGetShortcutSizeAck
// NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK
func ncCharOptionGetShortcutSizeAck(np *networking.Parameters) {
	// not sure what this data is
	hd, err := hex.DecodeString("0105000318000005000c00000c01000c02000c03000c040000")
	if err != nil {
		log.Error(err)
		return
	}
	nc := structs.NcCharOptionGetShortcutSizeAck{}
	err = structs.Unpack(hd, &nc)
	if err != nil {
		return
	}
	networking.Send(np.OutboundSegments.Send, networking.NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK, &nc)
}

// NcCharOptionGetWindowPosReq
// NC_CHAR_OPTION_GET_WINDOWPOS_REQ
func ncCharOptionGetWindowPosReq(ctx context.Context, np *networking.Parameters) {
	ncCharOptionGetWindowPosAck(np)
}

// NcCharOptionGetWindowPosAck
// NC_CHAR_OPTION_GET_WINDOWPOS_ACK
func ncCharOptionGetWindowPosAck(np *networking.Parameters) {
	hd, err := hex.DecodeString("01d707011800001c000000cdcc443d00000000cdccf03e00000000000000008ee3783f9a997b3f398ee33d000000001cc7713f00000000abaa6a3f00000000398e633f00000000c7715c3f9a99773f398ee33d9a99733f398ee33d9a996f3f398ee33d9a996b3f398ee33d0000903d2222023ecdcc663e4a9f943e0000983d4444843e9a99233e2222623ecd4c333f4444b43e00000000176c213e9a99993a8ee3b83e00000000610ba63e66e62e3f2ed8823e66662f3f610ba63ecd4c633f000000000000303ef549373fcd4c733ff549573f00000000832d203fcd4c143f4444243e3333803ededdbd3eb0010000730100008200000072000000000a0000a00500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	if err != nil {
		log.Error(err)
		return
	}
	nc := structs.NcCharOptionGetWindowPosAck{}
	err = structs.Unpack(hd, &nc)
	if err != nil {
		return
	}

	networking.Send(np.OutboundSegments.Send, networking.NC_CHAR_OPTION_GET_WINDOWPOS_ACK, &nc)
}

// NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_REQ
// 28727
func ncCharOptionImproveSetShortcutDataReq(ctx context.Context, np *networking.Parameters) {
	var e updateShortcutsEvent

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	e = updateShortcutsEvent{
		nc:          structs.NcCharOptionSetShortcutDataReq{},
		np:          np,
		characterID: session.characterID,
	}

	err := structs.Unpack(np.Command.Base.Data, e.nc)
	if err != nil {
		return
	}

	worldEvents[updateShortcuts] <- &e
}
