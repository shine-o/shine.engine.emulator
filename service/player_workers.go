package service

import (
	"reflect"
)

func (zm *zoneMap) playerActivity() {
	for {
		select {
		case e := <-zm.recv[playerData]:

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

			p.load(ev.playerName)
			//ncCharClientBaseCmd(ctx, &char) // todo: check if race condition
			//ncCharClientShapeCmd(ctx, char.Appearance)
			//
			//// todo: quest wrapper
			//ncCharClientQuestDoingCmd(ctx, &char)
			//ncCharClientQuestDoneCmd(ctx, &char)
			//ncCharClientQuestReadCmd(ctx, &char)
			//ncCharClientQuestRepeatCmd(ctx, &char)
			//
			//// todo: skills wrapper
			//ncCharClientPassiveCmd(ctx, &char)
			//ncCharClientSkillCmd(ctx, &char)
			//
			//ncCharClientItemCmd(ctx, char.AllEquippedItems(db))
			//ncCharClientItemCmd(ctx, char.InventoryItems(db))
			//ncCharClientItemCmd(ctx, char.MiniHouseItems(db))
			//ncCharClientItemCmd(ctx, char.PremiumActionItems(db))
			//
			//ncCharClientCharTitleCmd(ctx, &char)
			//
			//ncCharClientGameCmd(ctx)
			//ncCharClientChargedBuffCmd(ctx, &char)
			//ncCharClientCoinInfoCmd(ctx, &char)
			//ncQuestResetTimeClientCmd(ctx, &char)

		case e := <-zm.recv[playerAppeared]:
			// notify all nearby entities about it
			// players will get packet data
			// mobs will check if player is in range for attack
			ev, ok := e.(*playerAppearedEvent)
			if !ok {
				log.Errorf("expected event type %vEvent but got %v", playerAppeared, reflect.TypeOf(e).String())
			}
			zm.handles.mu.Lock()
			zm.handles.players[ev.player.handle] = ev.player
			zm.handles.mu.Unlock()

			go newPlayer(ev.player, zm.handles.players)
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
