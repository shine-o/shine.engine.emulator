package service

import (
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
)

// NC_BRIEFINFO_LOGINCHARACTER_CMD
func NcBriefInfoLoginCharacterCmd(conn playerConnection, nc * structs.NcBriefInfoLoginCharacterCmd) {
	pc := networking.Command{
		Base:     networking.CommandBase{
			OperationCode: 7174,
		},
		NcStruct: nc,
	}
	pc.SendDirectly(conn.data)
}