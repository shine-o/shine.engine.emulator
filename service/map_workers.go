package service

import (
	"reflect"
	"time"
)

func (zm *zoneMap) mapHandles() {
	log.Infof("[map_worker] mapHandles worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case e := <-zm.recv[handleCleanUp]:
			ev, ok := e.(*handleCleanUpEvent)
			if !ok {
				log.Errorf("expected event type %v but got %v", reflect.TypeOf(handleCleanUpEvent{}).String(), reflect.TypeOf(ev).String())
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
				log.Errorf("expected event type %v but got %v", reflect.TypeOf(&registerPlayerHandleEvent{}).String(), reflect.TypeOf(ev).String())
				break
			}
			//
			zm.entities.players.Lock()
			err := zm.entities.players.manager.newHandle()
			if err != nil {
				ev.err <- err
				break
			}
			handle := zm.entities.players.manager.index
			zm.entities.players.Unlock()

			ev.player.handle = handle
			ev.session.handle = handle
			go ev.player.heartbeat()
			ev.done <- true
		}
	}
}

func (zm *zoneMap) playerActivity() {
	log.Infof("[map_worker] playerActivity worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case e := <-zm.recv[playerAppeared]:
			// notify all nearby entities about it
			// players will get packet data
			// mobs will check if player is in range for attack
			ev, ok := e.(*playerAppearedEvent)
			if !ok {
				log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
				break
			}
			zm.entities.players.Lock()
			player := zm.entities.players.active[ev.playerHandle]
			zm.entities.players.Unlock()

			go newPlayer(player, zm.entities.players.active)
			go nearbyPlayers(player, zm.entities.players.active)

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

func (zm *zoneMap) playerQueries() {
	log.Infof("[map_worker] playerQueries worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case eq := <-zm.recv[queryPlayer]:
			log.Info(eq)
		}
	}
}

func (zm *zoneMap) monsterQueries() {
	log.Infof("[map_worker] monsterQueries worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case eq := <-zm.recv[queryMonster]:
			log.Info(eq)
		}
	}
}
