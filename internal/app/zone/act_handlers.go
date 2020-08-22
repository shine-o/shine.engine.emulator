package zone

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// NC_ACT_MOVERUN_CMD
// run
func ncActMoveRunCmd(ctx context.Context, np *networking.Parameters) {
	var (
		pre playerRunsEvent
		mqe queryMapEvent
	)

	session, ok := np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	pre = playerRunsEvent{
		nc: &structs.NcActMoveRunCmd{},
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

	zm.send[playerRuns] <- &pre
}

// NC_ACT_MOVEWALK_CMD
// walk
func ncActMoveWalkCmd(ctx context.Context, np *networking.Parameters) {}

// NC_ACT_STOP_REQ
// stop walk/run, a.k.a last known position of entity
func ncActStopReq(ctx context.Context, np *networking.Parameters) {}

// NC_ACT_JUMP_CMD
// jump
func ncActJumpCmd(ctx context.Context, np *networking.Parameters) {

}


// NC_ACT_SOMEONEMOVEWALK_CMD
// someone walked
func ncActSomeoneMoveWalkCmd(np *networking.Parameters, nc * structs.NcActSomeoneMoveWalkCmd) {}

// NC_ACT_SOMEONEMOVERUN_CMD
// someone has run
func ncActSomeoneMoveRunCmd(np *networking.Parameters, nc * structs.NcActSomeoneMoveWalkCmd) {}


// NC_ACT_SOMEONESTOP_CMD
// someone stopped
func ncActSomeoneStopCmd(np *networking.Parameters, nc * structs.NcActSomeoneStopCmd) {}
