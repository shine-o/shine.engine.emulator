package world

import (
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"reflect"
	"time"
)

func worldTime() structs.NcMiscGameTimeAck {
	t := time.Now()
	hour := byte(t.Hour())
	minute := byte(t.Minute())
	second := byte(t.Second())

	return structs.NcMiscGameTimeAck{
		Hour:   hour,
		Minute: minute,
		Second: second,
	}
}

func serverSelectLogic(ev *serverSelectEvent, w *world) {
	s, ok := ev.np.Session.(*session)

	if !ok {
		log.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
		return
	}

	err := verifyUser(s, ev.nc)

	if err != nil {
		log.Error(err)
		return
	}

	nc, err := userCharacters(w.db, s)

	if err != nil {
		log.Error(err)
		return
	}

	ncUserLoginWorldAck(ev.np, &nc)
}

func serverSelectTokenLogic(ev *serverSelectTokenEvent) {
	nc, err := returnToServerSelect()

	if err != nil {
		log.Error(err)
		return
	}

	ncUserWillWorldSelectAck(ev.np, &nc)
}
