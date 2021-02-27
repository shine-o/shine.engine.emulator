package zone

import (
	"fmt"
	"github.com/go-pg/pg/v9"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"reflect"
	"time"
)

func (z *zone) security() {
	log.Infof("[worker] security worker")
	for {
		select {
		case e := <-z.recv[playerSHN]:
			go playerSHNLogic(e)
		}
	}
}

func (z *zone) playerSession() {
	log.Infof("[zone_worker] playerSession worker")
	for {
		select {
		case e := <-z.recv[playerMapLogin]:
			go playerMapLoginLogic(e)
		case e := <-z.recv[playerData]:
			go playerDataLogic(e, z.worldDB)
		case e := <-z.recv[heartbeatUpdate]:
			go hearbeatUpdateLogic(e)
		case e := <-z.recv[playerLogoutStart]:
			go playerLogoutStartLogic(z, e)
		case e := <-z.recv[playerLogoutCancel]:
			go playerLogoutCancelLogic(z, e)
		case e := <-z.recv[playerLogoutConclude]:
			go playerLogoutConcludeLogic(z, e)
		case e := <-z.recv[changeMap]:
			go func() {
				log.Info(e)
				ev, ok := e.(*changeMapEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(changeMapEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				//todo: zone rpc method for external maps, for now, all maps are running in the same zone
				p := ev.p

				for _, v := range p.tickers {
					v.Stop()
				}

				handle := ev.p.getHandle()

				ev.prev.entities.players.remove(handle)
				ev.prev.entities.players.handler.remove(p.handle)

				ev.next.entities.players.add(p)

				ev.prev.send[playerDisappeared] <- &playerDisappearedEvent{
					handle: handle,
				}

				p.Lock()
				newLocation := *p.next

				p.fallback = newLocation
				p.current = newLocation

				p.next = nil

				p.players = make(map[uint16]*player)
				p.monsters = make(map[uint16]*monster)
				p.npcs = make(map[uint16]*npc)
				p.tickers = make([]*time.Ticker, 0)
				p.Unlock()

				nc := structs.NcMapLinkSameCmd{
					//MapID:    ev.next.data.Info.ID,
					MapID: ev.next.data.Info.ID,
					Location: structs.ShineXYType{
						X: uint32(newLocation.x),
						Y: uint32(newLocation.y),
					},
				}

				ev.s.mapID = ev.next.data.ID

				networking.Send(p.conn.outboundData, networking.NC_MAP_LINKSAME_CMD, &nc)

				ev.next.send[playerAppeared] <- &playerAppearedEvent{
					handle: handle,
				}
			}()
		}
	}
}

func (z *zone) playerGameData() {
	log.Infof("[zone_worker] playerGameData worker")
	for {
		select {
		case e := <-z.recv[persistPlayerPosition]:
			go persistPLayerPositionLogic(e, z)
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
					ev.err <- fmt.Errorf("map with id %v is not running on this zone", ev.id)
				}
				ev.zm <- zm
			}()
		}
	}
}

func playerSHNLogic(e event) {
	ev, ok := e.(*playerSHNEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerSHNEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	// u.u'
	ev.ok <- true
}

func playerMapLoginLogic(e event) {
	ev, ok := e.(*playerMapLoginEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerMapLoginEvent{}).String(), reflect.TypeOf(ev).String())
	}

	var (
		pse playerSHNEvent
		pde playerDataEvent
		qme queryMapEvent
		phe playerHandleEvent
	)

	pse = playerSHNEvent{
		inboundNC: ev.nc,
		ok:        make(chan bool),
		err:       make(chan error),
	}

	zoneEvents[playerSHN] <- &pse

	pde = playerDataEvent{
		player:     make(chan *player),
		net:        ev.np,
		playerName: ev.nc.CharData.CharID.Name,
		err:        make(chan error),
	}

	zoneEvents[playerData] <- &pde

	select {
	case <-pse.ok:
		break
	case err := <-pse.err:
		log.Error(err)
		// fail ack with failure code
		// drop connection
		return
	}

	var p *player
	select {
	case p = <-pde.player:
		break
	case err := <-pde.err:
		log.Error(err)
		// fail ack with failure code
		// drop connection
		return
	}

	qme = queryMapEvent{
		id:  p.current.mapID,
		zm:  make(chan *zoneMap),
		err: make(chan error),
	}

	zoneEvents[queryMap] <- &qme

	var zm *zoneMap
	select {
	case zm = <-qme.zm:
		break
	case err := <-qme.err:
		log.Error(err)
		return
	}

	session, ok := ev.np.Session.(*session)

	if !ok {
		log.Errorf("no session available for player %v", p.view.name)
		return
	}

	phe = playerHandleEvent{
		player:  p,
		session: session,
		done:    make(chan bool),
		err:     make(chan error),
	}

	zm.send[playerHandle] <- &phe

	select {
	case <-phe.done:

		nc := &structs.NcCharClientBaseCmd{
			ChrRegNum: uint32(p.char.ID),
			CharName: structs.Name5{
				Name: p.view.name,
			},
			Slot:       p.char.Slot,
			Level:      p.state.level,
			Experience: p.state.exp,
			PwrStone:   0,
			GrdStone:   0,
			HPStone:    p.stats.hpStones,
			SPStone:    p.stats.spStones,
			CurHP:      p.stats.hp,
			CurSP:      p.stats.sp,
			CurLP:      p.stats.lp,
			Unk:        1,
			Fame:       p.money.fame,
			Cen:        54983635, // ¿?
			LoginInfo: structs.NcCharBaseCmdLoginLocation{
				CurrentMap: structs.Name3{
					Name: p.current.mapName,
				},
				CurrentCoord: structs.ShineCoordType{
					XY: structs.ShineXYType{
						X: uint32(p.current.x),
						Y: uint32(p.current.y),
					},
					Direction: uint8(p.current.d),
				},
			},
			Stats: structs.CharStats{
				Strength:          p.stats.points.str,
				Constitute:        p.stats.points.end,
				Dexterity:         p.stats.points.dex,
				Intelligence:      p.stats.points.int,
				MentalPower:       p.stats.points.spr,
				RedistributePoint: p.stats.points.redistributionPoints,
			},
			IdleTime:   0,
			PkCount:    p.char.Attributes.KillPoints,
			PrisonMin:  0,
			AdminLevel: p.char.AdminLevel,
			Flag: structs.NcCharBaseCmdFlag{
				Val: 0,
			},
		}
		networking.Send(p.conn.outboundData, networking.NC_CHAR_CLIENT_BASE_CMD, nc)

		shape :=  p.view.protoAvatarShapeInfo()
		networking.Send(p.conn.outboundData, networking.NC_CHAR_CLIENT_SHAPE_CMD, shape)

		mapAck := &structs.NcMapLoginAck{
			Handle: p.handle, // id of the entity inside this map
			Params: p.charParameterData(),
			LoginCoord: structs.ShineXYType{
				X: uint32(p.current.x),
				Y: uint32(p.current.y),
			},
		}
		networking.Send(p.conn.outboundData, networking.NC_MAP_LOGIN_ACK, mapAck)

	case err := <-phe.err:
		log.Error(err)
	}
}

func playerDataLogic(e event, db *pg.DB) {
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

	err := p.load(ev.playerName, db)

	if err != nil {
		log.Error(err)
		ev.err <- err
	}
	ev.player <- p
}

func hearbeatUpdateLogic(e event) {
	ev, ok := e.(*heartbeatUpdateEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(heartbeatUpdateEvent{}).String(), reflect.TypeOf(ev).String())
	}

	var (
		mqe      queryMapEvent
		eventErr = make(chan error)
	)

	var (
		mapResult = make(chan *zoneMap)
		zm        *zoneMap
	)

	mqe = queryMapEvent{
		id:  ev.session.mapID,
		zm:  mapResult,
		err: eventErr,
	}

	zoneEvents[queryMap] <- &mqe

	select {
	case zm = <-mapResult:
		break
	case e := <-eventErr:
		log.Error(e)
		return
	}

	p := zm.entities.players.get(ev.session.handle)

	if p == nil {
		log.Errorf("nil player with handle %v", ev.session.handle)
		return
	}

	p.Lock()
	p.conn.lastHeartBeat = time.Now()
	p.Unlock()

	log.Infof("updating heartbeat for player %v", p.view.name)
}

func playerLogoutStartLogic(z *zone, e event) {
	ev, ok := e.(*playerLogoutStartEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerLogoutStartEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	m, ok := z.rm[ev.mapID]

	if !ok {
		log.Errorf("map with id %v not available", ev.mapID)
		return
	}

	p := m.entities.players.get(ev.handle)

	if p == nil {
		log.Errorf("player with handle %v not available", ev.handle)
		return
	}

	sid := ev.sessionID

	z.dynamicEvents.add(sid, dLogoutCancel)

	z.dynamicEvents.add(sid, dLogoutConclude)

	playerLogout(z, m, p, sid)
}

func playerLogoutCancelLogic(z *zone, e event) {
	ev, ok := e.(*playerLogoutCancelEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerLogoutCancelEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	z.dynamicEvents.RLock()
	defer z.dynamicEvents.RUnlock()

	sid := ev.sessionID

	select {
	case z.dynamicEvents.events[sid].send[dLogoutCancel] <- &emptyEvent{}:
		break
	default:
		log.Error("failed to send emptyEvent on dLogoutCancel")
		break
	}
}

func playerLogoutConcludeLogic(z *zone, e event) {
	ev, ok := e.(*playerLogoutConcludeEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerLogoutConcludeEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	z.dynamicEvents.RLock()
	defer z.dynamicEvents.RUnlock()

	sid := ev.sessionID

	select {
	case z.dynamicEvents.events[sid].send[dLogoutConclude] <- &emptyEvent{}:
		return
	default:
		log.Error("failed to send emptyEvent on dLogoutConclude")
		return
	}
}

func persistPLayerPositionLogic(e event, z *zone) {
	ev, ok := e.(*persistPlayerPositionEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(persistPlayerPositionEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	ev.p.Lock()
	c := ev.p.char
	c.Location.MapID = uint32(ev.p.current.mapID)
	c.Location.MapName = ev.p.current.mapName
	c.Location.X = ev.p.current.x
	c.Location.Y = ev.p.current.y
	c.Location.D = ev.p.current.d
	c.Location.IsKQ = false
	ev.p.Unlock()

	err := game.UpdateLocation(z.worldDB, c)

	if err != nil {
		log.Error(err)
		return
	}
}

// secondary workers that may be executed at runtime
func playerLogout(z *zone, zm *zoneMap, p *player, sid string) {
	t := time.NewTicker(15 * time.Second)
	defer t.Stop()
	finish := func() {
		t.Stop()
		select {
		case p.conn.close <- true:
			pde := &playerDisappearedEvent{
				handle: p.handle,
			}

			select {
			case zm.send[playerDisappeared] <- pde:
				break
			default:
				log.Error("unexpected error occurred while sending playerDisappeared event")
				break
			}

			break

		default:
			log.Error("unexpected error occurred while closing connection")
			return
		}
	}

	for {
		z.dynamicEvents.RLock()
		select {
		case <-z.dynamicEvents.events[sid].recv[dLogoutCancel]:
			z.dynamicEvents.RUnlock()
			return
		case <-z.dynamicEvents.events[sid].recv[dLogoutConclude]:
			z.dynamicEvents.RUnlock()
			finish()
			return
		case <-t.C:
			z.dynamicEvents.RUnlock()
			finish()
			return
		}
	}
}
