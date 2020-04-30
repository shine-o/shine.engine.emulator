package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/game/character"
	zm "github.com/shine-o/shine.engine.core/grpc/zone-master"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NcCharLoginReq handles a petition to login to the zone where the character's location map is running
// NC_CHAR_LOGIN_REQ
func NcCharLoginReq(ctx context.Context, np * networking.Parameters) {
	nc := structs.NcCharLoginReq{}
	if err := structs.Unpack(np.Command.Base.Data, &nc); err != nil {
		return
	}

	// get character where user_id and slot match
	var char character.Character
	err := db.Model(&char).
		Relation("Location").
		Relation("Options").Where("slot = ?", nc.Slot).Select()

	if err != nil {
		return
	}

	conn, err := newRPCClient("zone_master")
	c := zm.NewMasterClient(conn)
	rpcCtx, _ := context.WithTimeout(context.Background(), gRPCTimeout)

	ci, err := c.WhereIsMap(rpcCtx, &zm.MapQuery{
		Name: char.Location.MapName,
	})

	if err != nil {
		return
	}

	ack := structs.NcCharLoginAck{
		ZoneIP:   structs.Name4{
			Name: ci.IP,
		},
		ZonePort: uint16(ci.Port),
	}

	// not sure if these should be ordered
	NcCharOptionImproveGetGameOptionCmd(ctx, char.Options)
	NcCharOptionImproveGetKeymapCmd(ctx, char.Options)
	NcCharOptionImproveGetShortcutDataCmd(ctx, char.Options)

	go NcCharLoginAck(ctx, &ack)
}

// NcCharLoginAck sends zone connection info to the client
// NC_CHAR_LOGIN_ACK
func NcCharLoginAck(ctx context.Context, ack *structs.NcCharLoginAck) {
	// query the zone master for connection info for the map
	// send it to the client
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 4099,
		},
		NcStruct: ack,
	}
	pc.Send(ctx)
}
