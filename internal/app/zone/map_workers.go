package zone

import (
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
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
		case e := <-zm.recv[playerWalks]:
			// player has a fifo queue for the last 30 movements
			// for every movement
			//		verify collision
			//			if fails, return to previous movement
			// 		verify speed ( default 30 for unmounted/unbuffed player)
			//			if fails return to position 1 in queue
			//		broadcast to players within range
			go func() {
				ev, ok := e.(*playerWalksEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				rX := (ev.nc.To.X * 8) / 50
				rY := (ev.nc.To.X * 8) / 50

				zm.entities.players.RLock()
				defer zm.entities.players.RUnlock()
				player, ok := zm.entities.players.active[ev.handle]

				err := player.move(zm, rX, rY)
				if  err != nil {
					// store player position
					log.Error(err)
					return
				}

				nc := structs.NcActSomeoneMoveWalkCmd{
					Handle:   player.handle,
					From:     ev.nc.From,
					To:       ev.nc.To,
					Speed:    60,
				}
				for i, _ := range zm.entities.players.active {
					go ncActSomeoneMoveWalkCmd(zm.entities.players.active[i], &nc)
				}
			}()
		case e := <-zm.recv[playerRuns]:
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

				rX := (ev.nc.To.X * 8) / 50
				rY := (ev.nc.To.X * 8) / 50

				zm.entities.players.RLock()
				player, ok := zm.entities.players.active[ev.handle]

				err := player.move(zm, rX, rY)
				if  err != nil {
					// store player position
					log.Error(err)
					zm.entities.players.RUnlock()
					return
				}

				nc := structs.NcActSomeoneMoveRunCmd{
					Handle:   player.handle,
					From:     ev.nc.From,
					To:       ev.nc.To,
					Speed:    120,
				}
				for i, _ := range zm.entities.players.active {
					go ncActSomeoneMoveRunCmd(zm.entities.players.active[i], &nc)
				}
				zm.entities.players.RUnlock()

			}()
		case e := <-zm.recv[playerStopped]:
			// movements triggered by keys inmediately send a STOP packet to the server
			// movements triggered by mouse do not send a STOP packet
			// for every stop
			//		verify collision
			//		broadcast to players within range
			go func() {
				ev, ok := e.(*playerStoppedEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerStoppedEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				rX := (ev.nc.Location.X * 8) / 50
				rY := (ev.nc.Location.Y * 8) / 50

				zm.entities.players.RLock()
				player, ok := zm.entities.players.active[ev.handle]

				err := player.move(zm, rX, rY)
				if  err != nil {
					// store player position
					log.Error(err)
					zm.entities.players.RUnlock()
					return
				}
				nc := structs.NcActSomeoneStopCmd{
					Handle:   player.handle,
					Location: ev.nc.Location,
				}
				for i, _ := range zm.entities.players.active {
					go ncActSomeoneStopCmd(zm.entities.players.active[i], &nc)
				}
				zm.entities.players.RUnlock()

			}()
			log.Info(e)
		case e := <-zm.recv[playerJumped]:
			go func() {
				ev, ok := e.(*playerJumpedEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerJumpedEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				zm.entities.players.RLock()
				nc := structs.NcActSomeoneJumpCmd{
					Handle:   ev.handle,
				}
				for i, _ := range zm.entities.players.active {
					go ncActSomeoneJumpCmd(zm.entities.players.active[i], &nc)
				}
				zm.entities.players.RUnlock()
			}()
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
