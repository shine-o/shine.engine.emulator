package zone

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// NC_BAT_TARGETTING_REQ
// 9217
func ncBatTargetingReq(ctx context.Context, np *networking.Parameters) {
	var (
		psee playerSelectsEntityEvent
		mqe  queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	psee = playerSelectsEntityEvent{
		nc:     &structs.NcBatTargetInfoReq{},
		handle: session.handle,
	}

	err := structs.Unpack(np.Command.Base.Data, psee.nc)
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

	zm.send[playerSelectsEntity] <- &psee
}

func ncBatUntargetReq(ctx context.Context, np *networking.Parameters) {

	// remove current SelectionOrder value for player
	var (
		psee playerUnselectsEntityEvent
		mqe  queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	psee = playerUnselectsEntityEvent{
		handle: session.handle,
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

	zm.send[playerUnselectsEntity] <- &psee
}

// NC_BAT_TARGETINFO_CMD
// 9218
func ncBatTargetInfoCmd(p *player, nc *structs.NcBatTargetInfoCmd) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 9218,
		},
		NcStruct: nc,
	}
	pc.Send(p.conn.outboundData)
}
