package service

import (
	"reflect"
)

func (z *zone) playerSession() {
	for {
		select {
		case e := <-z.recv[playerData]:
			ev, ok := e.(*playerDataEvent)
			if !ok {
				log.Errorf("expected event type %vEvent but got %v", playerData, reflect.TypeOf(e).String())
			}

				p := &player{
				conn: playerConnection{
					close:        ev.net.NetVars.CloseConnection,
					outboundData: ev.net.NetVars.OutboundSegments.Send,
				},
			}

			err := p.load(ev.playerName)
			if err != nil {
				ev.err <- err
				break
			}
			ev.player <- p
		}
	}
}

func (zm *zoneMap) playerActivity() {
	for {
		select {
		case e := <-zm.recv[playerAppeared]:
			// notify all nearby entities about it
			// players will get packet data
			// mobs will check if player is in range for attack
			ev, ok := e.(*playerAppearedEvent)
			if !ok {
				log.Errorf("expected event type %vEvent but got %v", playerAppeared, reflect.TypeOf(ev).String())
				break
			}
			zm.entities.players.Lock()
			zm.entities.players.active[ev.player.handle] = ev.player
			zm.entities.players.Unlock()

			go newPlayer(ev.player, zm.entities.players.active)
			go nearbyPlayers(ev.player, zm.entities.players.active)

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
