package zone

import "github.com/shine-o/shine.engine.emulator/pkg/structs"

type playerHandleEvent struct {
	player  *player
	session *session
	done    chan bool
	err     chan error
}

type playerAppearedEvent struct {
	handle uint16
}

type playerDisappearedEvent struct {
	handle uint16
}

type playerWalksEvent struct {
	handle uint16
	nc     *structs.NcActMoveRunCmd
}

type playerRunsEvent struct {
	handle uint16
	nc     *structs.NcActMoveRunCmd
}

type playerJumpedEvent struct {
	handle uint16
}

type playerStoppedEvent struct {
	handle uint16
	nc     *structs.NcActStopReq
}

type unknownHandleEvent struct {
	handle uint16
	nc     *structs.NcBriefInfoInformCmd
}

type monsterWalksEvent struct {
	nc *structs.NcActSomeoneMoveWalkCmd
	m  *monster
}

type monsterRunsEvent struct {
	nc *structs.NcActSomeoneMoveRunCmd
	m  *monster
}

type playerSelectsEntityEvent struct {
	nc     *structs.NcBatTargetInfoReq
	handle uint16
}

type playerUnselectsEntityEvent struct {
	handle uint16
}

type playerClicksOnNpcEvent struct {
	nc     *structs.NcActNpcClickCmd
	handle uint16
}

type playerPromptReplyEvent struct {
	nc *structs.NcServerMenuAck
	s  *session
}
