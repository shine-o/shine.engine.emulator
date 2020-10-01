package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

func ncMapLogoutCmd(p *player, nc *structs.NcMapLogoutCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 6149,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}

// NC_MAP_TONORMALCOORD_CMD
// 6172