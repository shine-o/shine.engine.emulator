package zone

import "reflect"

func (z *zone) security() {
	log.Infof("[worker] security worker")
	for {
		select {
		case e := <-z.recv[playerSHN]:
			ev, ok := e.(*playerSHNEvent)
			if !ok {
				log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerSHNEvent{}).String(), reflect.TypeOf(ev).String())
				break
			}
			ev.ok <- true
		}
	}
}
