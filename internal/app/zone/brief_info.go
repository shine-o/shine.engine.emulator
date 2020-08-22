package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// NC_BRIEFINFO_LOGINCHARACTER_CMD
func ncBriefInfoLoginCharacterCmd(p *player, nc *structs.NcBriefInfoLoginCharacterCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 7174,
		},
		NcStruct: &nc,
	}
	pc.Send(p.conn.outboundData)
}

// NC_BRIEFINFO_CHARACTER_CMD
func ncBriefInfoCharacterCmd(p *player, nc *structs.NcBriefInfoCharacterCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 7175,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}
