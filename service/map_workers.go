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
		case <-zm.recv[playerHandleMaintenance]:
			go func() {
				zm.entities.players.Lock()
				for i, p := range zm.entities.players.active {
					p.Lock()
					if time.Since(p.conn.lastHeartBeat).Seconds() < playerHeartbeatLimit {
						p.Unlock()
						continue
					}
					p.Unlock()
					select {
					case zm.entities.players.active[i].send[heartbeatStop] <- &emptyEvent{}:
						time.Sleep(500 * time.Millisecond)
						p.Lock()
						delete(zm.entities.players.active, i)
						p.Unlock()
						// send event that notifies all players about the logout
					default:
						log.Error("failed to stop heartbeat")
						break
					}
				}
				zm.entities.players.Unlock()
			}()

		case e := <-zm.recv[playerHandle]:
			go func() {
				ev, ok := e.(*playerHandleEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerHandleEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				zm.entities.players.Lock()
				handle, err := zm.entities.players.newHandle()
				if err != nil {
					zm.entities.players.Unlock()
					ev.err <- err
					return
				}
				ev.player.Lock()
				zm.entities.players.active[handle] = ev.player
				zm.entities.players.Unlock()
				ev.player.handle = handle
				ev.session.handle = handle
				ev.session.mapID = ev.player.mapID
				ev.player.Unlock()
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
				zm.entities.players.Lock() // TODO: check if its necessary
				player, ok := zm.entities.players.active[ev.handle]
				if !ok {
					return
				}
				zm.entities.players.Unlock()
				go player.heartbeat()
				go newPlayer(player, &zm.entities.players)
				go nearbyPlayers(player, &zm.entities.players)
			}()

		case e := <-zm.recv[playerDisappeared]:
			log.Info(e)
		case e := <-zm.recv[playerRuns]:
			log.Info(e)

			// player has a fifo queue for the last 30 movements
			// for every movement
			//		verify collision
			//			if fails, return to previous movement
			// 		verify speed ( default 30 for unmounted/unbuffed player)
			//			if fails return to position 1 in queue
			//		broadcast to players within range
			// 		add to movements array
			go func() {
				ev, ok := e.(*playerRunsEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				// find player
				//ev.nc.
				//player, ok := zm.entities.players.active[ev.handle]
				//ev.nc.From
			}()
		case e := <-zm.recv[playerWalks]:
			// player has a fifo queue for the last 30 movements
			// for every movement
			//		verify collision
			//			if fails, return to previous movement
			// 		verify speed ( default 30 for unmounted/unbuffed player)
			//			if fails return to position 1 in queue
			//		broadcast to players within range

			log.Info(e)
		case e := <-zm.recv[playerStopped]:
			// movements triggered by keys inmediately send a STOP packet to the server
			// movements triggered by mouse do not send a STOP packet
			// for every stop
			//		verify collision
			//		broadcast to players within range

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
		case e := <-zm.recv[queryPlayer]:
			ev, ok := e.(*queryPlayerEvent)
			if !ok {
				log.Errorf("expected event type %v but got %v", reflect.TypeOf(queryPlayerEvent{}).String(), reflect.TypeOf(ev).String())
				return
			}
		}
	}
}

func (zm *zoneMap) monsterQueries() {
	log.Infof("[map_worker] monsterQueries worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case e := <-zm.recv[queryMonster]:
			log.Info(e)
		}
	}
}
