package login

import (
	"reflect"
)

func (l *login) authentication() {
	for {
		select {
		case e := <-l.events.recv[clientVersion]:
			go func() {
				ev, ok := e.(*clientVersionEvent)

				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&clientVersionEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				clientVersionLogic(ev)
			}()

		case e := <-l.events.recv[credentialsLogin]:
			go func() {
				ev, ok := e.(*credentialsLoginEvent)

				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&credentialsLoginEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				credentialsLoginLogic(ev, e, l)
			}()
		case e := <-l.events.recv[worldManagerStatus]:
			go func() {
				ev, ok := e.(*worldManagerStatusEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&worldManagerStatusEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				worldManagerStatusLogic(ev)
			}()
		//case e := <- l.events.recv[serverList]:

		case e := <-l.events.recv[serverSelect]:
			go func() {
				ev, ok := e.(*serverSelectEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverSelectEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				serverSelectLogic(l, ev)
			}()
		case e := <-l.events.recv[tokenLogin]:
			go func() {
				ev, ok := e.(*tokenLoginEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&tokenLoginEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				tokenLoginLogic(l, ev)
			}()
		}
	}
}
