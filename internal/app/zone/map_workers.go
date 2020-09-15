package zone

import (
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"reflect"
	"time"
)

const playerHeartbeatLimit = 10

func (zm *zoneMap) mapHandles() {
	log.Infof("[map_worker] mapHandles worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case <-zm.recv[playerHandleMaintenance]:
			go playerHandleMaintenanceLogic(zm)
		case e := <-zm.recv[playerHandle]:
			go playerHandleLogic(e, zm)
		}
	}
}

func (zm *zoneMap) playerActivity() {
	log.Infof("[map_worker] playerActivity worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case e := <-zm.recv[playerAppeared]:
			go playerAppearedLogic(e, zm)
		case e := <-zm.recv[playerDisappeared]:
			go playerDisappearedLogic(e, zm)
		case e := <-zm.recv[playerWalks]:
			go playerWalksLogic(e, zm)
		case e := <-zm.recv[playerRuns]:
			go playerRunsLogic(e, zm)
		case e := <-zm.recv[playerStopped]:
			go playerStoppedLogic(e, zm)
		case e := <-zm.recv[playerJumped]:
			go playerJumpedLogic(e, zm)
		case e := <-zm.recv[unknownHandle]:
			go unknownHandleLogic(e, zm)
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

func playerHandleMaintenanceLogic(zm *zoneMap) {
	zm.entities.players.RLock()
	for i := range zm.entities.players.active {

		var deleteHandle = func(mapInx uint16, zm * zoneMap) {
			zm.entities.players.Lock()
			delete(zm.entities.players.active, mapInx)
			zm.entities.players.Unlock()
		}

		p := zm.entities.players.active[i]

		p.RLock()

		lastHeartBeat := time.Since(p.conn.lastHeartBeat).Seconds()

		if lastHeartBeat < playerHeartbeatLimit {
			p.RUnlock()
			continue
		}

		pde := &playerDisappearedEvent{
			handle: p.handle,
		}

		select {
		case zm.events.send[playerDisappeared] <- pde:
			break
		default:
			log.Error("failed to stop heartbeat")
			break
		}

		for _, t := range p.tickers {
			t.Stop()
		}
		p.RUnlock()
		go deleteHandle(i, zm)
	}
	zm.entities.players.RUnlock()
}

func playerHandleLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerHandleEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerHandleEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	handle, err := zm.entities.players.newHandle()

	if err != nil {
		ev.err <- err
		return
	}

	zm.entities.players.Lock()
	zm.entities.players.active[handle] = ev.player
	zm.entities.players.Unlock()

	ev.player.Lock()
	ev.player.handle = handle
	ev.player.Unlock()

	ev.session.handle = handle
	ev.session.mapID = ev.player.mapID

	ev.done <- true
}

func playerAppearedLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerAppearedEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	zm.entities.players.RLock()
	player, ok := zm.entities.players.active[ev.handle]
	zm.entities.players.RUnlock()

	if !ok {
		log.Error("player not found during playerAppearedLogic")
		return
	}

	go player.heartbeat()
	go player.persistPosition()
	go player.nearbyPlayers(zm)

	go newPlayer(player, zm.entities.players)
	go nearbyPlayers(player, zm.entities.players)
}

func playerDisappearedLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerDisappearedEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerDisappearedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	zm.entities.players.RLock()

	for _, p := range zm.entities.players.active {
		if p.handle == ev.handle {
			continue
		}

		go ncMapLogoutCmd(p, &structs.NcMapLogoutCmd{
			Handle: ev.handle,
		})
	}

	zm.entities.players.RUnlock()
}

func playerWalksLogic(e event, zm *zoneMap) {
	// player has a fifo queue for the last 30 movements
	// for every movement
	//		verify collision
	//			if fails, return to previous movement
	// 		verify speed ( default 30 for unmounted/unbuffed player)
	//			if fails return to position 1 in queue
	//		broadcast to players within range
	ev, ok := e.(*playerWalksEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	rX := (ev.nc.To.X * 8) / 50
	rY := (ev.nc.To.Y * 8) / 50

	zm.entities.players.RLock()
	player, ok := zm.entities.players.active[ev.handle]
	zm.entities.players.RUnlock()

	if !ok {
		log.Error("player not found during playerStoppedLogic")
		return
	}

	err := player.move(zm, rX, rY)

	if err != nil {
		// extra validation steps to avoid speed hacks
		log.Error(err)
		return
	}

	player.Lock()
	player.x = ev.nc.To.X
	player.y = ev.nc.To.Y
	player.Unlock()

	nc := structs.NcActSomeoneMoveWalkCmd{
		Handle: player.handle,
		From:   ev.nc.From,
		To:     ev.nc.To,
		Speed:  60,
	}
	zm.entities.players.RLock()
	for i := range zm.entities.players.active {
		foreignPlayer := zm.entities.players.active[i]

		if foreignPlayer.handle != player.handle {
			if playerInRange(foreignPlayer, player) {
				go ncActSomeoneMoveWalkCmd(foreignPlayer, &nc)
			}
		}
	}
	zm.entities.players.RUnlock()

}

func playerRunsLogic(e event, zm *zoneMap) {
	// player has a fifo queue for the last 30 movements
	// for every movement
	//		verify collision
	//			if fails, return to previous movement
	// 		verify speed ( default 30 for unmounted/unbuffed player)
	//			if fails return to position 1 in queue
	//		broadcast to players within range
	// 		add to movements array
	ev, ok := e.(*playerRunsEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	rX := (ev.nc.To.X * 8) / 50
	rY := (ev.nc.To.Y * 8) / 50

	zm.entities.players.RLock()
	player, ok := zm.entities.players.active[ev.handle]
	zm.entities.players.RUnlock()

	if !ok {
		log.Error("player not found during playerStoppedLogic")
		return
	}

	err := player.move(zm, rX, rY)
	if err != nil {
		// extra validation steps to avoid speed hacks
		log.Error(err)
		return
	}

	player.Lock()
	player.x = ev.nc.To.X
	player.y = ev.nc.To.Y
	player.Unlock()

	nc := structs.NcActSomeoneMoveRunCmd{
		Handle: player.handle,
		From:   ev.nc.From,
		To:     ev.nc.To,
		Speed:  120,
	}

	zm.entities.players.RLock()
	for i := range zm.entities.players.active {
		foreignPlayer := zm.entities.players.active[i]

		if foreignPlayer.handle != player.handle {
			if playerInRange(foreignPlayer, player) {
				go ncActSomeoneMoveRunCmd(foreignPlayer, &nc)
			}
		}
	}
	zm.entities.players.RUnlock()
}

func playerStoppedLogic(e event, zm *zoneMap) {
	// movements triggered by keys inmediately send a STOP packet to the server
	// movements triggered by mouse do not send a STOP packet
	// for every stop
	//		verify collision
	//		broadcast to players within range
	ev, ok := e.(*playerStoppedEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerStoppedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	rX := (ev.nc.Location.X * 8) / 50
	rY := (ev.nc.Location.Y * 8) / 50

	zm.entities.players.RLock()
	player, ok := zm.entities.players.active[ev.handle]
	zm.entities.players.RUnlock()

	if !ok {
		log.Error("player not found during playerStoppedLogic")
		return
	}

	err := player.move(zm, rX, rY)

	if err != nil {
		log.Error(err)
		return
	}

	nc := structs.NcActSomeoneStopCmd{
		Handle:   player.handle,
		Location: ev.nc.Location,
	}

	zm.entities.players.RLock()
	for i := range zm.entities.players.active {
		foreignPlayer := zm.entities.players.active[i]

		if foreignPlayer.handle != player.handle {
			if playerInRange(foreignPlayer, player) {
				go ncActSomeoneStopCmd(foreignPlayer, &nc)
			}
		}
	}
	zm.entities.players.RUnlock()
}

func playerJumpedLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerJumpedEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerJumpedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	zm.entities.players.RLock()
	player, ok := zm.entities.players.active[ev.handle]
	zm.entities.players.RUnlock()

	if !ok {
		log.Error("player not found during playerStoppedLogic")
		return
	}

	nc := structs.NcActSomeoneJumpCmd{
		Handle: ev.handle,
	}

	zm.entities.players.RLock()

	for i := range zm.entities.players.active {
		foreignPlayer := zm.entities.players.active[i]
		if foreignPlayer.handle != player.handle {
			if playerInRange(foreignPlayer, player) {
				go ncActSomeoneJumpCmd(foreignPlayer, &nc)
			}
		}
	}
	zm.entities.players.RUnlock()
}

func unknownHandleLogic(e event, zm *zoneMap) {
	ev, ok := e.(*unknownHandleEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&unknownHandleEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	if ev.handle != ev.nc.AffectedHandle {
		log.Errorf("mismatched handles %v %v", ev.handle, ev.nc.AffectedHandle)
		return
	}

	zm.entities.players.RLock()
	player, ok := zm.entities.players.active[ev.handle]
	zm.entities.players.RUnlock()

	if !ok {
		log.Errorf("player with handle %v not found", ev.handle)
		return
	}

	//TODO: could also be NPC, Item on the ground or a Monster
	zm.entities.players.RLock()
	foreignPlayer, ok := zm.entities.players.active[ev.nc.ForeignHandle]
	zm.entities.players.RUnlock()

	if !ok {
		log.Errorf("player with handle %v not found", ev.handle)
		return
	}

	if playerInRange(player, foreignPlayer) {

		foreignPlayer.RLock()
		nc1 := foreignPlayer.ncBriefInfoLoginCharacterCmd()
		foreignPlayer.RUnlock()

		player.RLock()
		nc2 := player.ncBriefInfoLoginCharacterCmd()
		player.RUnlock()

		go ncBriefInfoLoginCharacterCmd(player, &nc1)
		go ncBriefInfoLoginCharacterCmd(foreignPlayer, &nc2)
	}
}

// notify every player in proximity about player that logged in
func newPlayer(p *player, nearbyPlayers *players) {
	nearbyPlayers.RLock()
	for i := range nearbyPlayers.active {
		foreignPlayer := nearbyPlayers.active[i]

		if p.handle != foreignPlayer.handle {
			if playerInRange(foreignPlayer, p) {
				nc := p.ncBriefInfoLoginCharacterCmd()
				go ncBriefInfoLoginCharacterCmd(foreignPlayer, &nc)
			}
		}
	}
	nearbyPlayers.RUnlock()
}

// send info to player about nearby players
func nearbyPlayers(p *player, nearbyPlayers *players) {
	var characters []structs.NcBriefInfoLoginCharacterCmd

	nearbyPlayers.RLock()

	for i := range nearbyPlayers.active {
		foreignPlayer := nearbyPlayers.active[i]

		if p.handle != foreignPlayer.handle {
			if playerInRange(foreignPlayer, p) {
				nc := foreignPlayer.ncBriefInfoLoginCharacterCmd()
				characters = append(characters, nc)
			}
		}

	}

	nearbyPlayers.RUnlock()

	ncBriefInfoCharacterCmd(p, &structs.NcBriefInfoCharacterCmd{
		Number:     byte(len(characters)),
		Characters: characters,
	})
}
