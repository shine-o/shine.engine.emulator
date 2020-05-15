package service

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

func (z *zone) playerSession() {
	log.Infof("[zone_worker] playerSession worker")
	for {
		select {
		case e := <-z.recv[loadPlayerData]:
			go func() {
				ev, ok := e.(*playerDataEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerDataEvent{}).String(), reflect.TypeOf(ev).String())
				}

				p := &player{
					conn: playerConnection{
						lastHeartBeat: time.Now(),
						close:         ev.net.NetVars.CloseConnection,
						outboundData:  ev.net.NetVars.OutboundSegments.Send,
					},
				}

				events := []eventIndex{heartbeatUpdate, heartbeatStop}
				p.recv = make(recvEvents)
				p.send = make(sendEvents)

				for _, index := range events {
					c := make(chan event, 2)
					p.recv[index] = c
					p.send[index] = c
				}

				err := p.load(ev.playerName)
				if err != nil {
					log.Error(err)
					ev.err <- err
				}
				ev.player <- p
			}()
		}
	}
}

func (z *zone) mapQueries() {
	log.Infof("[zone_worker] mapQueries worker")
	for {
		select {
		case e := <-z.recv[queryMap]:
			go func() {
				ev, ok := e.(*queryMapEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(queryMapEvent{}).String(), reflect.TypeOf(ev).String())
				}
				zm, ok := z.rm[ev.id]
				if !ok {
					ev.err <- errors.New(fmt.Sprintf("map with id %v is not running on this zone", ev.id))
				}
				ev.zm <- zm
			}()
		}
	}
}