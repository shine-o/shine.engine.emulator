package zone

import (
	"fmt"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game/character"
	"reflect"
	"time"
)

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

	zm.entities.players.Lock()
	p, ok := zm.entities.players.active[ev.session.handle]
	zm.entities.players.Unlock()

	p.Lock()
	p.conn.lastHeartBeat = time.Now()
	log.Infof("updating heartbeat for player %v", p.view.name)
	p.Unlock()
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
		id:  p.location.mapID,
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
		ncCharClientBaseCmd(p)
		ncCharClientShapeCmd(p)
		// weird bug sometimes the client stucks in character select
		ncMapLoginAck(p)
	case err := <-phe.err:
		log.Error(err)
	}
}

func playerDataLogic(e event) {
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

	err := p.load(ev.playerName)

	if err != nil {
		log.Error(err)
		ev.err <- err
	}
	ev.player <- p
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

	m.entities.players.Lock()
	p, ok := m.entities.players.active[ev.handle]
	m.entities.players.Unlock()

	if !ok {
		log.Errorf("map with id %v not available", ev.mapID)
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
	z.dynamicEvents.Lock()
	defer z.dynamicEvents.Unlock()

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

	z.dynamicEvents.Lock()
	defer z.dynamicEvents.Unlock()

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
	}
	ev.p.Lock()
	c := ev.p.char
	c.Location.MapID = uint32(ev.p.mapID)
	c.Location.MapName = ev.p.mapName
	c.Location.X = ev.p.x
	c.Location.Y = ev.p.y
	c.Location.D = ev.p.d
	c.Location.IsKQ = false
	ev.p.Unlock()

	err := character.Update(z.worldDB, c)

	if err != nil {
		log.Error(err)
		return
	}
}