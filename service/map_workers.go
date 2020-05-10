package service

import (
	"reflect"
	"time"
)

func (zm *zoneMap) mapHandles() {
	for {
		select {
		case e := <-zm.recv[handleCleanUp]:
			ev, ok := e.(*handleCleanUpEvent)
			if !ok {
				log.Errorf("expected event type %vEvent but got %v", handleCleanUp, reflect.TypeOf(ev).String())
				break
			}
			zm.entities.players.Lock()
			for _, ap := range zm.entities.players.active {
				if time.Since(ap.conn.lastHeartBeat).Seconds() > 15 {
					ap.conn.close <- true
					ap.send[heartbeatMissing] <- &emptyEvent{}
				}
			}
			zm.entities.players.Unlock()

		case e := <-zm.recv[registerPlayerHandle]:
			ev, ok := e.(*registerPlayerHandleEvent)
			if !ok {
				log.Errorf("expected event type %vEvent but got %v", registerPlayerHandle, reflect.TypeOf(ev).String())
				break
			}
			//
			zm.entities.players.Lock()
			err := zm.entities.players.manager.newHandle()
			if err != nil {
				ev.err <- err
			}
			ev.player.handle = zm.entities.players.manager.index
			ev.session.handle = zm.entities.players.manager.index
			zm.entities.players.Unlock()
		}
	}
}

func (zm *zoneMap) playerActivity() {
	for {
		select {
		case e := <-zm.recv[playerAppeared]:
			// notify all nearby entities about it
			// players will get packet data
			// mobs will check if player is in range for attack
			ev, ok := e.(*playerAppearedEvent)
			if !ok {
				log.Errorf("expected event type %vEvent but got %v", playerAppeared, reflect.TypeOf(ev).String())
				break
			}
			zm.entities.players.Lock()
			zm.entities.players.active[ev.player.handle] = ev.player
			zm.entities.players.Unlock()

			go newPlayer(ev.player, zm.entities.players.active)
			go nearbyPlayers(ev.player, zm.entities.players.active)

		case e := <-zm.recv[playerDisappeared]:
			log.Info(e)
		case e := <-zm.recv[playerMoved]:
			log.Info(e)
		case e := <-zm.recv[playerStopped]:
			log.Info(e)
		case e := <-zm.recv[playerJumped]:
			log.Info(e)
		}
	}
}
