package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

type playerMapLoginEvent struct {
	nc structs.NcMapLoginReq
	np *networking.Parameters
}

type playerHandleEvent struct {
	player  *player
	session *session
	done    chan bool
	err     chan error
}
