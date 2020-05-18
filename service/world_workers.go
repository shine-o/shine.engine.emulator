package service

import (
	"fmt"
	"reflect"
)

func (w *world) session()  {
	log.Info("[world_worker] session worker)")
	for {
		select {
		case e := <- w.recv[serverTime]:
			go func() {
				ev, ok := e.(*serverTimeEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverTimeEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				nc := worldTime()
				NcMiscGameTimeAck(ev.np, &nc)
			}()
		case e := <- w.recv[serverSelect]:
			go func() {
				ev, ok := e.(*serverSelectEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverSelectEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				s, ok := ev.np.Session.(*session)
				if !ok {
					ev.err <- fmt.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
				}
				err := verifyUser(s, ev.nc)
				if err != nil {
					ev.err <- err
				}
				nc, err := userCharacters(s)
				if err != nil {
					ev.err <- err
				}
				ncUserLoginWorldAck(ev.np, &nc)
			}()
		case e := <- w.recv[serverSelectToken]:
			go func() {
				ev, ok := e.(*serverSelectTokenEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverSelectTokenEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				nc, err := returnToServerSelect()
				if err != nil {
					ev.err <- err
				}
				ncUserWillWorldSelectAck(ev.np, &nc)
			}()
		}
	}
}