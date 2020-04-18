package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/game/character"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

//NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD
//28724
func gameOptionCmd(ctx context.Context, options * character.ClientOptions) {
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
		NcStruct: nc,
	}
	go pc.Send(ctx)
}

//NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD
//28723
func keymapCmd(ctx context.Context, options * character.ClientOptions) {
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
		NcStruct: nc,
	}
	go pc.Send(ctx)
}

//NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD
//28722
func shortcutDataCmd(ctx context.Context, options * character.ClientOptions) {
	// fighter: 040000040000000000010400010000000a0100ac0d00000b0100b10d0000
	nc := structs.NcCharGetShortcutDataCmd{}
	err := structs.Unpack(options.GameOptions, &nc)
	if err != nil {
		return
	}
	pc := networking.Command{
		Base:     networking.CommandBase{
			OperationCode:    28722,
		},
		NcStruct: nc,
	}
	go pc.Send(ctx)
}

//NC_CHAR_OPTION_GET_SHORTCUTSIZE_REQ
func getShortcutSizeReq(ctx context.Context, pc * networking.Command) {}

//NC_CHAR_OPTION_GET_SHORTCUTSIZE_ACK
func getShortcutSizeAck(ctx context.Context) {}

//NC_CHAR_OPTION_GET_WINDOWPOS_REQ
func getWindowPosReq(ctx context.Context, pc * networking.Command) {}

//NC_CHAR_OPTION_GET_WINDOWPOS_ACK
func getWindowPosAck(ctx context.Context) {}