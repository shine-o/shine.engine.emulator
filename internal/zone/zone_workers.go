package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
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
			go playerDataLogic(e)
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
				ev, ok := e.(*changeMapEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(changeMapEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				//todo: zone rpc method for external maps, for now, all maps are running in the same zone
				p := ev.p
				handle := p.getHandle()

				ev.prev.entities.players.remove(handle)
				ev.prev.entities.players.handler.remove(handle)

				ev.next.entities.players.add(p)

				ev.prev.send[playerDisappeared] <- &playerDisappearedEvent{
					handle: handle,
				}

				for _, v := range p.ticks.list {
					v.Stop()
				}

				newLocation := p.next

				p.Lock()
				p.baseEntity.fallback = newLocation
				p.baseEntity.current = newLocation

				p.ticks = &entityTicks{}
				p.Unlock()

				p.proximity.Lock()
				p.proximity.players = make(map[uint16]*player)
				p.proximity.npcs = make(map[uint16]*npc)
				p.proximity.Unlock()

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
			go persistPLayerPositionLogic(e)
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
		e1 playerSHNEvent
		e2 playerDataEvent
		e3 playerHandleEvent
	)

	e1 = playerSHNEvent{
		inboundNC: ev.nc,
		ok:        make(chan bool),
		err:       make(chan error),
	}

	zoneEvents[playerSHN] <- &e1

	e2 = playerDataEvent{
		player:     make(chan *player),
		net:        ev.np,
		playerName: ev.nc.CharData.CharID.Name,
		err:        make(chan error),
	}

	zoneEvents[playerData] <- &e2

	select {
	case <-e1.ok:
		break
	case err := <-e1.err:
		log.Error(err)
		// fail ack with failure code
		// drop connection
		return
	}

	var p *player
	select {
	case p = <-e2.player:
		break
	case err := <-e2.err:
		log.Error(err)
		// fail ack with failure code
		// drop connection
		return
	}

	zm, ok := maps.list[p.current.mapID]
	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"mapID": p.current.mapID,
			},
		})
		return
	}

	session, ok := ev.np.Session.(*session)

	if !ok {
		log.Errorf("no session available for player %v", p.view.name)
		return
	}

	e3 = playerHandleEvent{
		player:  p,
		session: session,
		done:    make(chan bool),
		err:     make(chan error),
	}

	zm.send[playerHandle] <- &e3

	select {
	case <-e3.done:
		networking.Send(p.conn.outboundData, networking.NC_CHAR_CLIENT_BASE_CMD, ncCharClientBase(p))
		networking.Send(p.conn.outboundData, networking.NC_CHAR_CLIENT_SHAPE_CMD, protoAvatarShapeInfo(p.view))
		networking.Send(p.conn.outboundData, networking.NC_MAP_LOGIN_ACK, ncMapLoginAck(p))
		networking.Send(p.conn.outboundData, networking.NC_CHAR_CLIENT_ITEM_CMD, ncCharClientItemCmd(p, persistence.BagInventory))
		networking.Send(p.conn.outboundData, networking.NC_CHAR_CLIENT_ITEM_CMD, ncCharClientItemCmd(p, persistence.EquippedInventory))
	case err := <-e3.err:
		log.Error(err)
	}
}

func playerDataLogic(e event) {
	ev, ok := e.(*playerDataEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerDataEvent{}).String(), reflect.TypeOf(ev).String())
	}

	p := &player{
		baseEntity: baseEntity{},
		conn: &playerConnection{
			lastHeartBeat: time.Now(),
			close:         ev.net.CloseConnection,
			outboundData:  ev.net.OutboundSegments.Send,
		},
	}

	err := p.load(ev.playerName)

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

	zm, ok := maps.list[ev.session.mapID]
	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneMapNotFound,
			Details: errors.ErrDetails{
				"session": ev.session,
			},
		})
		return
	}

	p, err := zm.entities.players.get(ev.session.handle)
	if err != nil {
		log.Error(err)
		return
	}

	p.conn.Lock()
	p.conn.lastHeartBeat = time.Now()
	p.conn.Unlock()

	log.Infof("updating heartbeat for player %v", p.view.name)
}

func playerLogoutStartLogic(z *zone, e event) {
	ev, ok := e.(*playerLogoutStartEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(playerLogoutStartEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	zm, ok := z.rm.list[ev.mapID]

	if !ok {
		log.Errorf("map with id %v not available", ev.mapID)
		return
	}

	p, err := zm.entities.players.get(ev.handle)
	if err != nil {
		log.Error(err)
		return
	}

	sid := ev.sessionID

	z.dynamicEvents.add(sid, dLogoutCancel)

	z.dynamicEvents.add(sid, dLogoutConclude)

	playerLogout(z, zm, p, sid)
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

func persistPLayerPositionLogic(e event) {
	ev, ok := e.(*persistPlayerPositionEvent)
	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(persistPlayerPositionEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	ev.p.persistence.Lock()
	c := ev.p.persistence.char
	ev.p.baseEntity.RLock()
	c.Location.MapID = uint32(ev.p.current.mapID)
	c.Location.MapName = ev.p.current.mapName
	c.Location.X = ev.p.current.x
	c.Location.Y = ev.p.current.y
	c.Location.D = ev.p.current.d
	ev.p.baseEntity.RUnlock()
	ev.p.persistence.Unlock()

	c.Location.IsKQ = false

	err := persistence.UpdateLocation(c)

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
				handle: p.getHandle(),
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

func ncMapLoginAck(p *player) *structs.NcMapLoginAck {
	nc := &structs.NcMapLoginAck{
		Handle: p.getHandle(), // id of the entity inside this map
		Params: p.charParameterData(),
	}
	p.baseEntity.RLock()
	nc.LoginCoord = structs.ShineXYType{
		X: uint32(p.current.x),
		Y: uint32(p.current.y),
	}
	p.baseEntity.RUnlock()
	return nc
}

func ncCharClientBase(p *player) *structs.NcCharClientBaseCmd {
	p.persistence.RLock()
	p.view.RLock()
	p.state.RLock()
	p.stats.RLock()
	p.money.RLock()
	p.baseEntity.RLock()
	nc := &structs.NcCharClientBaseCmd{
		ChrRegNum: uint32(p.persistence.char.ID),
		CharName: structs.Name5{
			Name: p.view.name,
		},
		Slot:       p.persistence.char.Slot,
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
		PkCount:    p.persistence.char.Attributes.KillPoints,
		PrisonMin:  0,
		AdminLevel: p.persistence.char.AdminLevel,
		Flag: structs.NcCharBaseCmdFlag{
			Val: 0,
		},
	}
	p.persistence.RUnlock()
	p.view.RUnlock()
	p.state.RUnlock()
	p.stats.RUnlock()
	p.money.RUnlock()
	p.baseEntity.RUnlock()
	return nc
}

// NC_CHAR_CLIENT_ITEM_CMD
func ncCharClientItemCmd(p *player, inventoryType persistence.InventoryType) *structs.NcCharClientItemCmd {
	p.inventories.RLock()
	defer p.inventories.RUnlock()

	nc := &structs.NcCharClientItemCmd{
		Box: byte(inventoryType),
		Flag: structs.ProtoNcCharClientItemCmdFlag{
			BF0: 211,
		},
	}

	switch inventoryType {
	case persistence.EquippedInventory:
		//p.inventories.inventory.
		nc.NumOfItem = byte(len(p.inventories.equipped.items))
		for _, item := range p.inventories.equipped.items {
			//           v7->location.Inven = ((_WORD)box << 10) ^ slot & 0x3FF;
			inc, err := protoItemPacketInformation(item)
			if err != nil {
				log.Error(err)
			}
			nc.Items = append(nc.Items, *inc)
		}
		break
	case persistence.BagInventory:
		nc.NumOfItem = byte(len(p.inventories.inventory.items))
		for _, item := range p.inventories.inventory.items {
			//           v7->location.Inven = ((_WORD)box << 10) ^ slot & 0x3FF;
			inc, err := protoItemPacketInformation(item)
			if err != nil {
				log.Error(err)
			}
			nc.Items = append(nc.Items, *inc)
		}
		break
	}

	return nc
}
