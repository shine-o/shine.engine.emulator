package service

import (
	"context"
	"encoding/hex"
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NcCharOptionImproveGetGameOptionCmd sends the character's game options to the client
// NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD
func NcCharOptionImproveGetGameOptionCmd(ctx context.Context, options * character.ClientOptions) {
	// fighter: 20000000010100010200010300010400000500010600000700010800010900010a00010b00010c00000d00000e00000f00001000011100001200001300001400011500001600001700001800011900011a00011b00011c00011d00001e00001f0001
	nc := structs.NcCharOptionImproveGetGameOptionCmd{}
	err := structs.Unpack(options.GameOptions, &nc)
	if err != nil {
		return
	}
	pc := networking.Command{
		Base:     networking.CommandBase{
			OperationCode:    28724,
		},
		NcStruct: &nc,
	}
	pc.Send(ctx)
}

// NcCharOptionImproveGetKeymapCmd sends the character's key map settings to the client
// NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD
func NcCharOptionImproveGetKeymapCmd(ctx context.Context, options * character.ClientOptions) {
	// fighter: 5f00000000790100001b02000043030000490400004b0500004c0600004607000048080000560900000d0a00114e0b0011470c0011500d0011570e0000de0f000058100000471100000012000000130000001400000015000052160011411700005718000053190000001a0000411b0000441c00005a1d0000201e0000261f000028200000252100002722000024230000542400005125000045260000f527000042280000502900004d2a0000552b00105a2c0000002d0000002e0000232f000031300000323100003332000034330000353400003635000037360000383700003938000030390000bd3a0000bb3b0010313c0010323d0010333e0010343f0010354000103641001037420010384300103944001030450010bd460010bb4700123148001232490012334a0012344b0012354c0012364d0012374e0012384f00123950001230510012bd520012bb530000005400000055000000560000005700000058000000590000005a0000005b0000005c0000005d0000005e000000
	nc := structs.NcCharGetKeyMapCmd{}
	err := structs.Unpack(options.Keymap, &nc)
	if err != nil {
		return
	}
	pc := networking.Command{
		Base:     networking.CommandBase{
			OperationCode:    28723,
		},
		NcStruct: &nc,
	}
	pc.Send(ctx)
}

// NcCharOptionImproveGetShortcutDataCmd sends the character's shortcut settings to the client
// NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD
func NcCharOptionImproveGetShortcutDataCmd(ctx context.Context, options * character.ClientOptions) {
	nc := structs.NcCharGetShortcutDataCmd{}
	err := structs.Unpack(options.Shortcuts, &nc)
	if err != nil {
		return
	}
	pc := networking.Command{
		Base:     networking.CommandBase{
			OperationCode:    28722,
		},
		NcStruct: &nc,
	}
	pc.Send(ctx)
}

// NcCharOptionGetShortcutSizeReq ~~~ unknown yet what this data is for
// NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ
func NcCharOptionGetShortcutSizeReq(ctx context.Context, pc * networking.Command) {
	// gotta handle this
	go NcCharOptionGetShortcutSizeAck(ctx)
}

// NcCharOptionGetShortcutSizeAck ~~~ unknown yet what this data is for
// NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK
func NcCharOptionGetShortcutSizeAck(ctx context.Context) {
	// not sure what this data is
	hd, err := hex.DecodeString("0105000318000005041000000c01000c02000c03000c040000")
	if err != nil {
		log.Error(err)
		return
	}
	nc := structs.NcCharOptionGetShortcutSizeAck{}
	err = structs.Unpack(hd, &nc)
	if err != nil {
		return
	}
	pc := networking.Command{
		Base:     networking.CommandBase{
			OperationCode:    28677,
		},
		NcStruct: &nc,
	}
	go pc.Send(ctx)
}

// NcCharOptionGetWindowPosReq ~~~ unknown yet what this data is for
// NC_CHAR_OPTION_GET_WINDOWPOS_REQ
func NcCharOptionGetWindowPosReq(ctx context.Context, pc * networking.Command) {
	go NcCharOptionGetWindowPosAck(ctx)
}

// NcCharOptionGetWindowPosAck ~~~ unknown yet what this data is for
// NC_CHAR_OPTION_GET_WINDOWPOS_ACK
func NcCharOptionGetWindowPosAck(ctx context.Context) {
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
	pc := networking.Command{
		Base:     networking.CommandBase{
			OperationCode:    28685,
		},
		NcStruct: &nc,
	}
	go pc.Send(ctx)
}