package service

import "reflect"

func (z *zone) playerSession() {
	for {
		select {
		case e := <-z.recv[loadPlayerData]:
			ev, ok := e.(*playerDataEvent)
			if !ok {
				log.Errorf("expected event type %vEvent but got %v", loadPlayerData, reflect.TypeOf(e).String())
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
