package zone

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
		case e := <-z.recv[playerData]:
			go func() {
				ev, ok := e.(*playerDataEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerDataEvent{}).String(), reflect.TypeOf(ev).String())
				}

				p := &player{
					conn: playerConnection{
						lastHeartBeat: time.Now(),
						close:         ev.net.CloseConnection,
						outboundData:  ev.net.OutboundSegments.Send,
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
		case e := <-z.recv[playerLogoutStart]:
			go func() {
				ev, ok := e.(*playerLogoutStartEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerLogoutStartEvent{}).String(), reflect.TypeOf(ev).String())
				}
				m, ok := z.rm[ev.mapID]
				if !ok {
					ev.err <- fmt.Errorf("map with id %v not available", ev.mapID)
					return
				}
				m.entities.players.Lock()
				p, ok := m.entities.players.active[ev.handle]
				m.entities.players.Unlock()
				if !ok {
					ev.err <- fmt.Errorf("map with id %v not available", ev.mapID)
				}
				cancel := z.dynamic.add(ev.sessionID, dLogoutCancel)
				conclude := z.dynamic.add(ev.sessionID, dLogoutConclude)
				go playerLogout(cancel, conclude, m, p)
			}()
		case e := <-z.recv[playerLogoutCancel]:
			go func() {
				ev, ok := e.(*playerLogoutCancelEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerLogoutCancelEvent{}).String(), reflect.TypeOf(ev).String())
				}
				z.dynamic.Lock()
				defer z.dynamic.Unlock()
				select {
				case z.dynamic.events[ev.sessionID].send[dLogoutCancel] <- &emptyEvent{}:
					return
				default:
					log.Error("failed to send event")
					return
				}
			}()
		case e := <-z.recv[playerLogoutConclude]:
			go func() {
				ev, ok := e.(*playerLogoutConcludeEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerLogoutConcludeEvent{}).String(), reflect.TypeOf(ev).String())
				}
				z.dynamic.Lock()
				defer z.dynamic.Unlock()
				select {
				case z.dynamic.events[ev.sessionID].send[dLogoutConclude] <- &emptyEvent{}:
					return
				default:
					log.Error("failed to send event")
					return
				}
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
