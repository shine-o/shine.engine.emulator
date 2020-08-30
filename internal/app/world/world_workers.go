package world

import (
	"reflect"
)

func (w *world) session() {
	log.Info("[world_worker] session worker")
	for {
		select {
		case e := <-w.recv[serverTime]:
			go func() {
				ev, ok := e.(*serverTimeEvent)

				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverTimeEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				nc := worldTime()

				NcMiscGameTimeAck(ev.np, &nc)
			}()
		case e := <-w.recv[serverSelect]:
			go func() {
				ev, ok := e.(*serverSelectEvent)

				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverSelectEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				serverSelectLogic(ev, w)
			}()
		case e := <-w.recv[serverSelectToken]:
			go func() {
				ev, ok := e.(*serverSelectTokenEvent)

				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverSelectTokenEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				serverSelectTokenLogic(ev)
			}()
		}
	}
}
