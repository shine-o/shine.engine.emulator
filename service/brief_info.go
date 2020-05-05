package service

import (
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NC_BRIEFINFO_LOGINCHARACTER_CMD
func ncBriefInfoLoginCharacterCmd(p *player, nc *structs.NcBriefInfoLoginCharacterCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 7174,
		},
		NcStruct: nc,
	}
	pc.SendDirectly(p.conn.data)
}

// NC_BRIEFINFO_CHARACTER_CMD
func ncBriefInfoCharacterCmd(p *player, nc *structs.NcBriefInfoCharacterCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 7175,
		},
		NcStruct: nc,
	}
	pc.SendDirectly(p.conn.data)
}
