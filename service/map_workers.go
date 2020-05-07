package service

func (zm *zoneMap) playerActivity() {
	for {
		select {
		case e := <-zm.recv[playerAppeared]:
			// notify all nearby entities about it
			// players will get packet data
			// mobs will check if player is in range for attack
			ev := e.(playerAppearedEvent)
			zm.handles.mu.Lock()
			zm.handles.players[ev.player.handle] = ev.player
			zm.handles.mu.Unlock()
			go newPlayer(zm, ev)
			go nearbyPlayers(ev.player, zm.handles.players)

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
