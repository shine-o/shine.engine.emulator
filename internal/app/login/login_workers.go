package login

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"google.golang.org/grpc/connectivity"
	"reflect"
)

func (l *login) authentication() {
	for {
		select {
		case e := <-l.events.recv[clientVersion]:
			go clientVersionLogic(e)
		case e := <-l.events.recv[credentialsLogin]:
			go credentialsLoginLogic(e, l)
		case e := <-l.events.recv[worldManagerStatus]:
			go worldManagerStatusLogic(e)
		case e := <-l.events.recv[serverSelect]:
			go serverSelectLogic(l, e)
		case e := <-l.events.recv[tokenLogin]:
			go tokenLoginLogic(l, e)
		}
	}
}

func clientVersionLogic(e event) {
	ev, ok := e.(*clientVersionEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&clientVersionEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	err := checkClientVersion(ev.nc)
	if err != nil {
		log.Error(err)
		ncUserClientWrongVersionCheckAck(ev.np)
		return
	}
	ncUserClientRightVersionCheckAck(ev.np)
}

func worldManagerStatusLogic(e event) {
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
		log.Error("connection with the world master was not Ready")
		return
	}

	ncUserWorldStatusAck(ev.np)
}

func credentialsLoginLogic(e event, l *login) {
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
}

func tokenLoginLogic(l *login, e event) {
	ev, ok := e.(*tokenLoginEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&tokenLoginEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	_, err := redisClient.Get(ev.nc.Otp.Name).Result()
	if err != nil {
		log.Error(err)
		ncUserLoginFailAck(ev.np, 69)
		return
	}

	loginSuccessful(l, ev.np)
}

func loginSuccessful(l *login, np *networking.Parameters) {
	nc := structs.NcUserLoginAck{}
	for _, w := range l.worlds {
		nc.Worlds = append(nc.Worlds, structs.WorldInfo{
			WorldNumber: byte(w.id),
			WorldName: structs.Name4{
				Name: w.name,
			},
			WorldStatus: 6,
		})
	}
	nc.NumOfWorld = byte(len(l.worlds))
	ncUserLoginAck(np, nc)
}

func serverSelectLogic(l *login, e event) {
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
}
