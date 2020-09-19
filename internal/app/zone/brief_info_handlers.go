package zone

import (
	"context"
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

// NC_BRIEFINFO_BRIEFINFODELETE_CMD
// 7182
func ncBriefInfoDeleteHandleCmd(p *player, nc *structs.NcBriefInfoDeleteHandleCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 7182,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}

// NC_BRIEFINFO_INFORM_CMD
// 7169
func ncBriefInfoInformCmd(ctx context.Context, np *networking.Parameters) {
	// trigger handleInfo
	// if targetHandle is within range of affectedHandle
	//		send NC_BRIEFINFO_LOGINCHARACTER_CMD of the targetHandle to the affectedHandle
	var (
		uhe unknownHandleEvent
		mqe queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	uhe = unknownHandleEvent{
		handle: session.handle,
		nc:     &structs.NcBriefInfoInformCmd{},
	}

	err := structs.Unpack(np.Command.Base.Data, uhe.nc)
	if err != nil {
		log.Error(err)
		return
	}

	mqe = queryMapEvent{
		id:  session.mapID,
		zm:  make(chan *zoneMap),
		err: make(chan error),
	}

	zoneEvents[queryMap] <- &mqe

	var zm *zoneMap
	select {
	case zm = <-mqe.zm:
		break
	case e := <-mqe.err:
		log.Error(e)
		return
	}

	zm.send[unknownHandle] <- &uhe
}

// NC_BRIEFINFO_REGENMOB_CMD
// 7176
func ncBriefInfoRegenMobCmd(p * player, nc * structs.NcBriefInfoRegenMobCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 7176,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}
