package service

import (
	"errors"
	"github.com/shine-o/shine.engine.core/structs"
	"google.golang.org/grpc/connectivity"
	"reflect"
)

func (l *login) authentication ()  {
	for {
		select {
			case e := <- l.events.recv[clientVersion]:
				go func() {
					ev, ok := e.(*clientVersionEvent)
					if !ok {
						log.Errorf("expected event type %v but got %v", reflect.TypeOf(&clientVersionEvent{}).String(), reflect.TypeOf(ev).String())
						return
					}
					 err := checkClientVersion(ev.nc)
					if err != nil {
						log.Error(e)
						ncUserClientWrongVersionCheckAck(ev.np)
						return
					}
					ncUserClientRightVersionCheckAck(ev.np)
				}()

			case e := <- l.events.recv[credentialsLogin]:
				go func() {
					ev, ok := e.(*credentialsLoginEvent)
					if !ok {
						log.Errorf("expected event type %v but got %v", reflect.TypeOf(&credentialsLoginEvent{}).String(), reflect.TypeOf(ev).String())
						return
					}
					err := checkCredentials(ev.nc)
					if err != nil {
						log.Error(e)
						ncUserLoginFailAck(ev.np, 69)
						return
					}
					loginSuccessful(l, ev.np)
				}()
		case e := <-l.events.recv[worldManagerStatus]:
			go func() {
				ev, ok := e.(*worldManagerStatusEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&worldManagerStatusEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				conn, err := newRPCClient("world_master")

				if err != nil {
					log.Error(err)
					return
				}
				defer conn.Close()
				if conn.GetState() != connectivity.Ready {
					ncUserWorldStatusAck(ev.np)
					return
				} else {
					log.Error(errors.New("connection with the world master was not Ready"))
					return
				}
			}()
		//case e := <- l.events.recv[serverList]:

		case e := <-l.events.recv[serverSelect]:
			go func() {
				ev, ok := e.(*serverSelectEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&serverSelectEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				for _, w := range l.worlds {
					if byte(w.id) == ev.nc.WorldNo {
						nc := structs.NcUserWorldSelectAck{
							WorldStatus: 6,
							Ip: structs.Name4{
								Name: w.ip,
							},
							Port: uint16(w.port),
						}
						ncUserWorldSelectAck(ev.np, nc)
						return
					}
				}
				log.Errorf("failed to find world with ID %v in %v", ev.nc.WorldNo, l.worlds)
			}()
		case e := <-l.events.recv[tokenLogin]:
			go func() {
				ev, ok := e.(*tokenLoginEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&tokenLoginEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				err := loginByCode(ev.nc)
				if err != nil {
					log.Error(err)
					ncUserLoginFailAck(ev.np, 69)
					return
				}
				//TODO: move to own event; loginSuccessful
				loginSuccessful(l, ev.np)
				return
			}()
		}
	}
}
