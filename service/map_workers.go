package service

import (
	"reflect"
	"time"
)

const playerHeartbeatLimit = 30

func (zm *zoneMap) mapHandles() {
	log.Infof("[map_worker] mapHandles worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case <-zm.recv[handleCleanUp]:
			go func() {
				zm.entities.players.Lock()
				for i, ap := range zm.entities.players.active {
					ap.RLock()
					if time.Since(ap.conn.lastHeartBeat).Seconds() < playerHeartbeatLimit {
						ap.RUnlock()
						continue
					}
					ap.Lock()
					select {
					case zm.entities.players.active[i].send[heartbeatStop] <- &emptyEvent{}:
					default:
						log.Error("failed to stop heartbeat")
						return
					}
					time.Sleep(500 * time.Millisecond)
					delete(zm.entities.players.active, i)
					ap.Unlock()
				}
				zm.entities.players.Unlock()
			}()

		case e := <-zm.recv[registerPlayerHandle]:
			go func() {
				ev, ok := e.(*registerPlayerHandleEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&registerPlayerHandleEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				zm.entities.players.Lock()
				handle, err := zm.entities.players.newHandle()
				if err != nil {
					zm.entities.players.Unlock()
					ev.err <- err
					return
				}
				zm.entities.players.active[handle] = ev.player
				zm.entities.players.Unlock()

				ev.player.handle = handle
				ev.session.handle = handle
				ev.session.mapID = ev.player.mapID
				ev.done <- true
			}()
		}
	}
}

func (zm *zoneMap) playerActivity() {
	log.Infof("[map_worker] playerActivity worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case e := <-zm.recv[playerAppeared]:
			go func() {
				ev, ok := e.(*playerAppearedEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				zm.entities.players.Lock()
				player := zm.entities.players.active[ev.playerHandle]
				zm.entities.players.Unlock()
				player.RLock()
				go player.heartbeat()
				go newPlayer(player, zm.entities.players.active)
				go nearbyPlayers(player, zm.entities.players.active)
				player.RUnlock()
			}()

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
