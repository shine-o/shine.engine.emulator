package service

import (
	"context"
	"github.com/shine-o/shine.engine.core/networking"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/spf13/viper"
)

func ncMiscSeedAck(ctx context.Context, np *networking.Parameters) {
	xov := ctx.Value(networking.XorOffset)
	xc := xov.(chan uint16)

	xorLimit := uint16(viper.GetInt("crypt.xorLimit"))

	xorOffset := networking.RandomXorKey(xorLimit)
	log.Infof("[xor offset] %v", xorOffset)

	nc := structs.NcMiscSeedAck{
		Seed: xorOffset,
	}

	np.Command.NcStruct = &nc

	np.Command.Send(ctx)

	xc <- xorOffset
}

// for the zone service it is the server that makes use of this handler
func ncMiscHeartBeatReq(p *player) {
	pc := networking.Command{
		Base: networking.CommandBase{
			OperationCode: 2052,
		},
	}
	pc.SendDirectly(p.conn.outboundData)
}

// for the zone service it is the client that makes use of this handler
func ncMiscHeartBeatAck(ctx context.Context, np *networking.Parameters) {
	var (
		mqe      queryMapEvent
		eventErr = make(chan error)
	)

	sv := ctx.Value(networking.ShineSession)
	session, ok := sv.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	var (
		mapResult = make(chan *zoneMap)
		zm        *zoneMap
	)

	mqe = queryMapEvent{
		id:  session.mapID,
		zm:  mapResult,
		err: eventErr,
	}

	zoneEvents[queryMap] <- &mqe

	select {
	case zm = <-mapResult:
		break
	case e := <-eventErr:
		log.Error(e)
		return
	}
	zm.entities.players.Lock()
	p, ok := zm.entities.players.active[session.handle]
	zm.entities.players.Unlock()
	if !ok {
		log.Error("player handle not available for character %v on map %v", session.characterName, zm.data.Info.MapName)
	}

	p.send[heartbeatUpdate] <- &emptyEvent{}
}
