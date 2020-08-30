package zone

import "time"

func playerLogoutStartLogic(z *zone, ev *playerLogoutStartEvent) {
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

	z.dynamicEvents.Lock()

	sid := ev.sessionID

	z.dynamicEvents.add(sid, dLogoutCancel)
	z.dynamicEvents.add(sid, dLogoutConclude)

	z.dynamicEvents.Unlock()

	playerLogout(z, m, p, sid)
}

func playerLogoutCancelLogic(z *zone, ev *playerLogoutCancelEvent) {
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

func playerLogoutConcludeLogic(z *zone, ev *playerLogoutConcludeEvent) {
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

func playerLogout(z *zone, zm *zoneMap, p *player, sid string) {
	t := time.NewTicker(15 * time.Second)

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
		select {
		case <-z.dynamicEvents.events[sid].recv[dLogoutCancel]:
			t.Stop()
			return
		case <-z.dynamicEvents.events[sid].recv[dLogoutConclude]:
			finish()
			return
		case <-t.C:
			finish()
			return
		}
	}
}
