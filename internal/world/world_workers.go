package world

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"reflect"
	"time"
)

func (w *world) session() {
	log.Info("[world_worker] session worker")
	for {
		select {
		case e := <-w.recv[serverTime]:
			go serverTimeLogic(e)
		case e := <-w.recv[serverSelect]:
			go serverSelectLogic(e, w)
		case e := <-w.recv[serverSelectToken]:
			go serverSelectTokenLogic(e)
		}
	}
}

func serverTimeLogic(e event) {
	ev, ok := e.(*serverTimeEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverTimeEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	t := time.Now()
	hour := byte(t.Hour())
	minute := byte(t.Minute())
	second := byte(t.Second())

	nc :=  &structs.NcMiscGameTimeAck{
		Hour:   hour,
		Minute: minute,
		Second: second,
	}

	networking.Send(ev.np.OutboundSegments.Send, networking.NC_MISC_GAMETIME_ACK, nc)

}

func serverSelectLogic(e event, w *world) {
	ev, ok := e.(*serverSelectEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverSelectEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

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

	nc, err := userCharacters(s)

	if err != nil {
		log.Error(err)
		return
	}
	networking.Send(ev.np.OutboundSegments.Send, networking.NC_USER_LOGINWORLD_ACK, &nc)
	//networking.Send(ev.np.OutboundSegments.Send, networking.NC_USER_LOGOUT_DB, &nc)
}

func serverSelectTokenLogic(e event) {
	ev, ok := e.(*serverSelectTokenEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverSelectTokenEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	nc, err := returnToServerSelect()

	if err != nil {
		log.Error(err)
		return
	}
	networking.Send(ev.np.OutboundSegments.Send, networking.NC_USER_WILL_WORLD_SELECT_ACK, &nc)
}
