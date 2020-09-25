package zone

import (
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"reflect"
	"sync"
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

func (zm *zoneMap) monsterActivity() {
	log.Infof("[map_worker] monsterActivity worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case e := <-zm.recv[monsterAppeared]:
			log.Info(e)
		case e := <-zm.recv[monsterDisappeared]:
			log.Info(e)
		case e := <-zm.recv[monsterWalks]:
			go func() {
				ev, ok := e.(*monsterWalksEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&monsterWalksEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				for ap := range zm.entities.players.all() {
					go func(p *player, m *monster) {
						if monsterInRange(p, m) {
							ncActSomeoneMoveWalkCmd(p, ev.nc)
						}
					}(ap, ev.m)
				}
			}()
		case e := <-zm.recv[monsterRuns]:
			go func() {
				ev, ok := e.(*monsterRunsEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&monsterRunsEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				for ap := range zm.entities.players.all() {
					go func(p *player, m *monster) {
						if monsterInRange(p, m) {
							go ncActSomeoneMoveRunCmd(p, ev.nc)
						}
					}(ap, ev.m)
				}
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
				log.Errorf("expected event type %v but got %v", reflect.TypeOf(&queryPlayerEvent{}).String(), reflect.TypeOf(ev).String())
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
	for p := range zm.entities.players.all() {

		if p.spawned() {
			continue
		}

		lhb := lastHeartbeat(p)

		if lhb < playerHeartbeatLimit {
			continue
		}

		p.RLock()

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

		go zm.entities.players.remove(p.handle)

		p.RUnlock()

	}
}

func playerHandleLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerHandleEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerHandleEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	handle, err := zm.entities.players.handler.new(playerHandleMin, playerHandleMax, playerAttemptsMax)

	if err != nil {
		ev.err <- err
		return
	}

	ev.player.Lock()
	ev.player.handle = handle
	ev.player.Unlock()

	zm.entities.players.add(ev.player)

	ev.session.handle = handle
	ev.session.mapID = ev.player.current.mapID

	ev.done <- true
}

func playerAppearedLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerAppearedEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerAppearedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	
	p1 := zm.entities.players.get(ev.handle)

	if p1 == nil {
		log.Error("player not found")
		return
	}

	var wg sync.WaitGroup

	wg.Add(4)
	go func() {
		defer wg.Done()
		newPlayer(p1, zm)
	}()

	go func() {
		defer wg.Done()
		nearbyPlayers(p1, zm)
	}()

	go func() {
		defer wg.Done()
		nearbyMonsters(p1, zm)
	}()

	go func() {
		defer wg.Done()
		p1.allNPC(zm)
	}()

	wg.Wait()

	go p1.heartbeat()
	go p1.persistPosition()

	go p1.nearbyPlayersMaintenance(zm)

	go p1.nearbyMonstersMaintenance(zm)

	go p1.nearbyNPCMaintenance(zm)

	//go adjacentMonstersInform(p1, zm)
}

func (p * player) allNPC(zm * zoneMap)  {
	var npcs structs.NcBriefInfoMobCmd

	for n := range zm.entities.npcs.all() {
		npcs.Mobs = append(npcs.Mobs, n.ncBriefInfoRegenMobCmd())
	}

	npcs.MobNum = byte(len(npcs.Mobs))

	ncBriefInfoMobCmd(p, &npcs)

}

func playerDisappearedLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerDisappearedEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerDisappearedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	for ap := range zm.entities.players.all() {
		go func(p2 *player) {
			if p2.getHandle() == ev.handle {
				return
			}
			ncMapLogoutCmd(p2, &structs.NcMapLogoutCmd{
				Handle: ev.handle,
			})
		}(ap)
	}
}

const (
	runSpeed = 120
	walkSpeed = 60
	//runSpeed = 300
	//walkSpeed = 150
)
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

	p1 := zm.entities.players.get(ev.handle)

	if p1 == nil {
		log.Error("player not found")
		return
	}

	if !ok {
		log.Error("player not found during playerStoppedLogic")
		return
	}

	igX := int(ev.nc.To.X)
	igY := int( ev.nc.To.Y)

	rX, rY := igCoordToBitmap(igX, igY)

	err := p1.move(zm, rX, rY)

	if err != nil {
		// extra validation steps to avoid speed hacks
		log.Error(err)
		return
	}

	p1.Lock()
	p1.current.x = igX
	p1.current.y = igY
	p1.Unlock()

	nc := structs.NcActSomeoneMoveWalkCmd{
		Handle: ev.handle,
		From:   ev.nc.From,
		To:     ev.nc.To,
		Speed:  walkSpeed,
	}

	for ap := range zm.entities.players.all() {
		go func(p2 *player) {
			if p2.getHandle() != ev.handle {
				if playerInRange(p2, p1) {
					go ncActSomeoneMoveWalkCmd(p2, &nc)
				}
			}
		}(ap)
	}
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

	p1 := zm.entities.players.get(ev.handle)

	if p1 == nil {
		log.Error("player not found")
		return
	}

	igX := int(ev.nc.To.X)
	igY := int(ev.nc.To.Y)

	rX, rY := igCoordToBitmap(igX, igY)

	err := p1.move(zm, rX, rY)
	if err != nil {
		// extra validation steps to avoid speed hacks
		log.Error(err)
		return
	}

	p1.Lock()
	p1.current.x = igX
	p1.current.y = igY
	p1.Unlock()

	nc := structs.NcActSomeoneMoveRunCmd{
		Handle: ev.handle,
		From:   ev.nc.From,
		To:     ev.nc.To,
		Speed:  runSpeed,
	}

	for ap := range zm.entities.players.all() {
		go func(p2 *player) {
			if p2.getHandle() != ev.handle {
				if playerInRange(p2, p1) {
					ncActSomeoneMoveRunCmd(p2, &nc)
				}
			}
		}(ap)
	}

	//go adjacentMonstersInform(p1, zm)

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

	p1 := zm.entities.players.get(ev.handle)

	if p1 == nil {
		log.Error("player not found")
		return
	}

	igX := int(ev.nc.Location.X)
	igY := int( ev.nc.Location.Y)

	rX, rY := igCoordToBitmap(igX, igY)

	err := p1.move(zm, rX, rY)

	if err != nil {
		log.Error(err)
		return
	}

	p1.Lock()
	p1.current.x = igX
	p1.current.y = igY
	p1.Unlock()

	nc := structs.NcActSomeoneStopCmd{
		Handle:   ev.handle,
		Location: ev.nc.Location,
	}

	for ap := range zm.entities.players.all() {
		go func(p2 *player) {
			if p2.getHandle() != ev.handle {
				if playerInRange(p2, p1) {
					ncActSomeoneStopCmd(p2, &nc)
				}
			}
		}(ap)
	}

}

func playerJumpedLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerJumpedEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerJumpedEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p1 := zm.entities.players.get(ev.handle)

	if p1 == nil {
		log.Error("player not found")
		return
	}

	nc := structs.NcActSomeoneJumpCmd{
		Handle: ev.handle,
	}

	for ap := range zm.entities.players.all() {
		go func(p2 *player) {
			if p2.getHandle() != ev.handle {
				if playerInRange(p2, p1) {
					go ncActSomeoneJumpCmd(p2, &nc)
				}
			}
		}(ap)
	}
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

	p1 := zm.entities.players.get(ev.handle)

	if p1 == nil {
		log.Error("player not found")
		return
	}

	//TODO: could also be NPC, Item on the ground or a Monster

	p2 := zm.entities.players.get(ev.nc.ForeignHandle)

	if p2 == nil {
		return
	}

	if playerInRange(p1, p2) {

		p1.RLock()
		nc1 := p1.ncBriefInfoLoginCharacterCmd()
		p1.RUnlock()
		go ncBriefInfoLoginCharacterCmd(p2, &nc1)

		p2.RLock()
		nc2 := p2.ncBriefInfoLoginCharacterCmd()
		p2.RUnlock()
		go ncBriefInfoLoginCharacterCmd(p1, &nc2)
	}


	m := zm.entities.monsters.get(ev.nc.ForeignHandle)

	if m == nil {
		return
	}

	if monsterInRange(p1, m) {
		nc := m.ncBriefInfoRegenMobCmd()
		go ncBriefInfoRegenMobCmd(p1, &nc)
	}

}

// notify every player in proximity about player that logged in
func newPlayer(p1 *player, zm *zoneMap) {
	for ap := range zm.entities.players.all() {
		go func(p2 *player) {
			if p1.getHandle() != p2.getHandle() {
				if playerInRange(p2, p1) {
					nc := p1.ncBriefInfoLoginCharacterCmd()
					ncBriefInfoLoginCharacterCmd(p2, &nc)
				}
			}
		}(ap)
	}
}

// send info to player about nearby players
func nearbyPlayers(p1 *player, zm *zoneMap) {
	var characters []structs.NcBriefInfoLoginCharacterCmd

	for p2 := range zm.entities.players.all() {
		if p1.getHandle() != p2.getHandle() {
			if playerInRange(p2, p1) {
				nc := p2.ncBriefInfoLoginCharacterCmd()
				characters = append(characters, nc)
			}
		}
	}

	ncBriefInfoCharacterCmd(p1, &structs.NcBriefInfoCharacterCmd{
		Number:     byte(len(characters)),
		Characters: characters,
	})
}

func nearbyMonsters(p *player, zm *zoneMap) {
	for am := range zm.entities.monsters.all() {
		go func(p * player, m *monster) {
			if monsterInRange(p, m) {
				nc := m.ncBriefInfoRegenMobCmd()
				ncBriefInfoRegenMobCmd(p, &nc)
			}
		}(p, am)
	}
}

