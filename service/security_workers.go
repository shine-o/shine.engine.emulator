package service

import "reflect"

func (z *zone) security() {
	for {
		select{
		case e := <- z.recv[clientSHN]:
			log.Info(e)
			ev, ok := e.(*clientSHNEvent)
			if !ok {
				log.Errorf("expected event type %vEvent but got %v", clientSHN, reflect.TypeOf(ev).String())
				break
			}
			ev.ok <- true
		}
	}
}