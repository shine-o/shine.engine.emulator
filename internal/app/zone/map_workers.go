package zone

import (
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/world"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"reflect"
	"strings"
	"sync"
)

const (
	runSpeed             = 120
	walkSpeed            = 60
	playerHeartbeatLimit = 10
)

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
		case e := <-zm.recv[playerSelectsEntity]:
			go playerSelectsEntityLogic(zm, e)
		case e := <-zm.recv[playerUnselectsEntity]:
			go playerUnselectsEntityLogic(zm, e)
		}
	}
}

func (zm *zoneMap) npcInteractions() {
	log.Infof("[map_worker] npcInteractions worker for map %v", zm.data.Info.MapName)
	for {
		select {
		case e := <-zm.recv[playerClicksOnNpc]:
			go playerClicksOnNpcLogic(zm, e)
		case e := <-zm.recv[playerPromptReply]:
			go playerPromptReplyLogic(zm, e)
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
			go monsterWalksLogic(zm, e)
		case e := <-zm.recv[monsterRuns]:
			go monsterRunsLogic(zm, e)
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

func playerClicksOnNpcLogic(zm *zoneMap, e event) {
	//log.Info(e)
	ev, ok := e.(*playerClicksOnNpcEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerClicksOnNpcEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	p := zm.entities.players.get(ev.handle)
	if p == nil {
		log.Errorf("player not found %v", ev.handle)
		return
	}
	// find npc with handle in ev.nc.Handle
	// send id of mob
	var nc structs.NcActNpcMenuOpenReq
	for n := range zm.entities.npcs.all() {
		if n.handle == ev.nc.NpcHandle {
			nc.MobID = n.mobInfo.ID
			if isPortal(n) {

				var md *world.Map

				for i, m := range mapData {
					if m.Info.MapName.Name == n.npcData.ShinePortal.ServerMapIndex {
						md = mapData[i]
						break
					}
				}

				var mapName string
				if md != nil {
					mapName = md.Info.Name
				} else {
					mapName = "UNAVAILABLE"
				}

				mapName = strings.Replace(mapName, "#", " ", -1)

				title := fmt.Sprintf("Do you want to move to %v", mapName)

				nc := structs.NcServerMenuReq{
					Title:     title,
					Priority:  0,
					NpcHandle: n.getHandle(),
					NpcPosition: structs.ShineXYType{
						X: uint32(n.current.x),
						Y: uint32(n.current.y),
					},
					LimitRange: 350,
					MenuNumber: 2,
					Menu: []structs.ServerMenu{
						{
							Reply:   1,
							Content: "Yes.",
						},
						{
							Reply:   0,
							Content: "No.",
						},
					},
				}

				go networking.Send(p.conn.outboundData, networking.NC_MENU_SERVERMENU_REQ, &nc)
				if md == nil {
					return
				}

				p.Lock()
				p.next = &location{
					mapID:     md.ID,
					mapName:   md.Info.MapName.Name,
					x:         n.npcData.ShinePortal.X,
					y:         n.npcData.ShinePortal.Y,
					movements: [15]movement{},
				}
				p.Unlock()
				return
			}
			break
		}
	}
	networking.Send(p.conn.outboundData, networking.NC_ACT_NPCMENUOPEN_REQ, &nc)
}

func playerPromptReplyLogic(zm *zoneMap, e event) {
	log.Info(e)
	ev, ok := e.(*playerPromptReplyEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerPromptReplyEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	if ev.nc.Reply == 0 {
		return
	}

	p := zm.entities.players.get(ev.s.handle)

	if p == nil {
		log.Errorf("player not found %v", ev.s.handle)
		return
	}

	if p.targeting.selectingN == nil {
		log.Warning("prompt cannot be answered, player is no longer selecting an NPC")
		return
	}

	// for now, its only about portals
	if isPortal(p.targeting.selectingN) && portalMatchesLocation(p.targeting.selectingN.npcData.ShinePortal, p.next) {
		// move player to map
		mqe := queryMapEvent{
			id:  p.next.mapID,
			zm:  make(chan *zoneMap),
			err: make(chan error),
		}

		zoneEvents[queryMap] <- &mqe

		var nzm *zoneMap
		select {
		case nzm = <-mqe.zm:
			break
		case e := <-mqe.err:
			log.Error(e)
			return
		}

		cme := changeMapEvent{
			p:    p,
			s:    ev.s,
			prev: zm,
			next: nzm,
		}

		zoneEvents[changeMap] <- &cme
	}
}

func playerSelectsEntityLogic(zm *zoneMap, e event) {
	ev, ok := e.(*playerSelectsEntityEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerSelectsEntityEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	vp := zm.entities.players.get(ev.handle)
	if vp == nil {
		log.Errorf("player not found %v", ev.handle)
		return
	}

	p, m, n := findFirstEntity(zm, ev.nc.TargetHandle)

	// set timeout in case of nonexistent handle
	// or use bool channels, and in the default case check if all three are false and return if so
	var nc *structs.NcBatTargetInfoCmd

	var notP, notM, notN bool
	for {
		select {
		case ap := <-p:
			if ap == nil {
				notP = true
				break
			}

			nc = ap.ncBatTargetInfoCmd()

			order := vp.selectsPlayer(ap)

			nc.Order = order
			networking.Send(vp.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nc)

			if ap.targeting.selectingP != nil {
				nextNc := ap.targeting.selectingP.ncBatTargetInfoCmd()
				nextNc.Order = order + 1
				networking.Send(vp.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nextNc)
			}

			if ap.targeting.selectingM != nil {
				nextNc := ap.targeting.selectingM.ncBatTargetInfoCmd()
				nextNc.Order = order + 1
				networking.Send(vp.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nextNc)
			}

			if ap.targeting.selectingN != nil {
				nextNc := ap.targeting.selectingN.ncBatTargetInfoCmd()
				nextNc.Order = order + 1
				networking.Send(vp.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nextNc)
			}

			for p := range vp.selectedByPlayers() {
				nextNc := *nc
				nextNc.Order++
				networking.Send(p.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nextNc)
			}

			return
		case am := <-m:
			if am == nil {
				notM = true
				break
			}

			order := vp.selectsMonster(am)

			nc = am.ncBatTargetInfoCmd()

			nc.Order = order

			networking.Send(vp.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nc)

			//if vp is being selected by player
			//send them information about the monster
			for p := range vp.selectedByPlayers() {
				nextNc := *nc
				nextNc.Order++
				networking.Send(p.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nextNc)
			}
			return
		case an := <-n:
			if an == nil {
				notN = true
				break
			}
			order := vp.selectsNPC(an)

			nc = an.ncBatTargetInfoCmd()

			nc.Order = order

			networking.Send(vp.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nc)

			for p := range vp.selectedByPlayers() {
				nextNc := *nc
				nextNc.Order++
				networking.Send(p.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, nextNc)
			}
			return
		default:
			if notP && notM && notN {
				return
			}
		}
	}
}

func playerUnselectsEntityLogic(zm *zoneMap, e event) {
	ev, ok := e.(*playerUnselectsEntityEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerUnselectsEntityEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	vp := zm.entities.players.get(ev.handle)
	if vp == nil {
		log.Errorf("player not found %v", ev.handle)
		return
	}

	var order byte
	vp.Lock()
	order = vp.targeting.selectionOrder
	vp.targeting.selectingP = nil
	vp.targeting.selectingM = nil
	vp.targeting.selectingN = nil
	vp.Unlock()

	nc := structs.NcBatTargetInfoCmd{
		Order:         order + 1,
		Handle:        65535,
		TargetHP:      0,
		TargetMaxHP:   0,
		TargetSP:      0,
		TargetMaxSP:   0,
		TargetLP:      0,
		TargetMaxLP:   0,
		TargetLevel:   0,
		HpChangeOrder: 0,
	}
	vp.RLock()
	for _, p := range vp.targeting.selectedByP {
		go networking.Send(p.conn.outboundData, networking.NC_BAT_TARGETINFO_CMD, &nc)
	}
	vp.RUnlock()
}

func monsterWalksLogic(zm *zoneMap, e event) {
	ev, ok := e.(*monsterWalksEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&monsterWalksEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	for ap := range zm.entities.players.all() {
		go func(p *player, m *monster) {
			if monsterInRange(p, m) {
				networking.Send(p.conn.outboundData, networking.NC_ACT_SOMEONEMOVEWALK_CMD, ev.nc)
			}
		}(ap, ev.m)
	}
}

func monsterRunsLogic(zm *zoneMap, e event) {
	ev, ok := e.(*monsterRunsEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&monsterRunsEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	for ap := range zm.entities.players.all() {
		go func(p *player, m *monster) {
			if monsterInRange(p, m) {
				go networking.Send(p.conn.outboundData, networking.NC_ACT_SOMEONEMOVERUN_CMD, ev.nc)
			}
		}(ap, ev.m)
	}
}

func playerHandleMaintenanceLogic(zm *zoneMap) {
	for ap := range zm.entities.players.all() {

		go func(p *player) {

			if p.spawned() {
				return
			}

			lhb := lastHeartbeat(p)

			if lhb < playerHeartbeatLimit {
				return
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
			go zm.entities.players.handler.remove(p.handle)

			p.RUnlock()

		}(ap)
	}
	zm.metrics.players.Set(float64(len(zm.entities.players.active)))
}

func playerHandleLogic(e event, zm *zoneMap) {
	ev, ok := e.(*playerHandleEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&playerHandleEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	handle, err := zm.entities.players.handler.new()

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
		showAllNPC(p1, zm)
	}()

	wg.Wait()

	go p1.heartbeat()

	go p1.persistPosition()
	go p1.nearbyPlayersMaintenance(zm)
	go p1.nearbyMonstersMaintenance(zm)
	go p1.nearbyNPCMaintenance(zm)

	//go adjacentMonstersInform(p1, zm)
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
			nc := &structs.NcBriefInfoDeleteHandleCmd{
				Handle: ev.handle,
			}
			networking.Send(p2.conn.outboundData, networking.NC_BRIEFINFO_BRIEFINFODELETE_CMD, nc)
		}(ap)
	}

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
					go networking.Send(p2.conn.outboundData, networking.NC_ACT_SOMEONEMOVEWALK_CMD, &nc)
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
					networking.Send(p2.conn.outboundData, networking.NC_ACT_SOMEONEMOVERUN_CMD, &nc)
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
	igY := int(ev.nc.Location.Y)

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
					networking.Send(p2.conn.outboundData, networking.NC_ACT_SOMEONESTOP_CMD, &nc)
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
					go networking.Send(p2.conn.outboundData, networking.NC_ACT_SOMEEONEJUMP_CMD, &nc)
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

	// handle is persisted across maps X_x
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
		go networking.Send(p2.conn.outboundData, networking.NC_BRIEFINFO_LOGINCHARACTER_CMD, &nc1)

		p2.RLock()
		nc2 := p2.ncBriefInfoLoginCharacterCmd()
		p2.RUnlock()
		go networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_LOGINCHARACTER_CMD, &nc2)

	}

	m := zm.entities.monsters.get(ev.nc.ForeignHandle)

	if m == nil {
		return
	}

	if monsterInRange(p1, m) {
		nc := m.ncBriefInfoRegenMobCmd()
		go networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_REGENMOB_CMD, &nc)
	}

}

// notify every player in proximity about player that logged in
func newPlayer(p1 *player, zm *zoneMap) {
	for ap := range zm.entities.players.all() {
		go func(p2 *player) {
			if p1.getHandle() != p2.getHandle() {
				if playerInRange(p2, p1) {
					nc := p1.ncBriefInfoLoginCharacterCmd()
					networking.Send(p2.conn.outboundData, networking.NC_BRIEFINFO_LOGINCHARACTER_CMD, &nc)
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

	nc := &structs.NcBriefInfoCharacterCmd{
		Number:     byte(len(characters)),
		Characters: characters,
	}

	networking.Send(p1.conn.outboundData, networking.NC_BRIEFINFO_CHARACTER_CMD, nc)
}

func nearbyMonsters(p *player, zm *zoneMap) {
	for am := range zm.entities.monsters.all() {
		go func(p *player, m *monster) {
			if monsterInRange(p, m) {
				nc := m.ncBriefInfoRegenMobCmd()
				networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_REGENMOB_CMD, &nc)
			}
		}(p, am)
	}
}

func findFirstEntity(zm *zoneMap, handle uint16) (chan *player, chan *monster, chan *npc) {
	p := make(chan *player, 1)
	m := make(chan *monster, 1)
	n := make(chan *npc, 1)

	go func(p chan<- *player, zm *zoneMap, targetHandle uint16) {
		for ap := range zm.entities.players.all() {
			if ap.getHandle() == targetHandle {
				p <- ap
				return
			}
		}
		p <- nil
	}(p, zm, handle)

	go func(m chan<- *monster, zm *zoneMap, targetHandle uint16) {
		for am := range zm.entities.monsters.all() {
			if am.getHandle() == targetHandle {
				m <- am
				return
			}
		}
		m <- nil
	}(m, zm, handle)

	go func(n chan<- *npc, zm *zoneMap, targetHandle uint16) {
		for an := range zm.entities.npcs.all() {
			if an.getHandle() == targetHandle {
				n <- an
				return
			}
		}
		n <- nil
	}(n, zm, handle)
	return p, m, n
}

func showAllNPC(p *player, zm *zoneMap) {
	var npcs structs.NcBriefInfoMobCmd

	for n := range zm.entities.npcs.all() {
		info := n.ncBriefInfoRegenMobCmd()
		// if portal, FlagState = 1
		// FlagData, Destination Map Index
		if isPortal(n) {
			info.FlagState = 1
			info.FlagData.Data = n.npcData.ServerMapIndex
		}
		npcs.Mobs = append(npcs.Mobs, info)
	}

	npcs.MobNum = byte(len(npcs.Mobs))

	networking.Send(p.conn.outboundData, networking.NC_BRIEFINFO_MOB_CMD, &npcs)
}

func isPortal(n *npc) bool {
	if n.npcData == nil {
		return false
	}

	if n.npcData.ShinePortal == nil {
		return false
	}

	var loaded bool
	for _, m := range mapData {
		if m.MapInfoIndex == n.npcData.ShinePortal.ServerMapIndex {
			loaded = true
			break
		}
	}
	if loaded {
		return true
	}
	return false
}

func portalMatchesLocation(portal *world.ShinePortal, next *location) bool {
	var md *world.Map
	for i, m := range mapData {
		if m.MapInfoIndex == portal.ClientMapIndex {
			md = mapData[i]
			break
		}
	}

	if md == nil {
		return false
	}

	if md.ID == next.mapID {
		return true
	}

	return false
}
