package service

func (zm *zoneMap) playerMovement() {
	for {
		select {
		case e := <-zm.recv[playerAppeared]:
			// notify all nearby entities about it
			// players will get packet data
			// mobs will check if player is in range for attack
			if e.eventType() != playerAppeared {
				log.Errorf("unexpected event %v", e.eventType())
				return
			}
			ev := e.(playerAppearedEvent)
			for _, p := range zm.handles.players {
				if p.getHandle() == p.getHandle() {
					return
				}
				go NcBriefInfoLoginCharacterCmd(p.conn, &ev.nc)
			}

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
