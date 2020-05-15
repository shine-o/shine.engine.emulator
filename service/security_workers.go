package service

import "reflect"

func (z *zone) security() {
	log.Infof("[worker] security worker")
	for {
		select {
		case e := <-z.recv[clientSHN]:
			ev, ok := e.(*clientSHNEvent)
			if !ok {
				log.Errorf("expected event type %v but got %v", reflect.TypeOf(clientSHNEvent{}).String(), reflect.TypeOf(ev).String())
				break
			}
			ev.ok <- true
		}
	}
}
