package zone

import "reflect"

func (z *zone) security() {
	log.Infof("[worker] security worker")
	for {
		select {
		case e := <-z.recv[playerSHN]:
			go playerSHNLogic(e)
		}
	}
}

func playerSHNLogic(e event) {
	ev, ok := e.(*playerSHNEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerSHNEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	// u.u'
	ev.ok <- true
}
