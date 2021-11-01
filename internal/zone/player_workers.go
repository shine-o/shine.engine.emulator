package zone

import (
	"reflect"
)

// combination of events that must be processed in order
// the reasoning is to avoid certain events changing data at the same time
func (p *player) eduPlayerSyncEvents001() {
	log.Infof("[player_worker] eduPlayerSyncEvents001 worker for player %v", p.getHandle())
	for {
		select {
		case e := <-p.events.recv[eduStats]:
			log.Info(e)
		case e := <-p.events.recv[eduState]:
			log.Info(e)
		case e := <-p.events.recv[eduEquipItem]:
			eduEquipItemLogic(e, p)
		case e := <-p.events.recv[eduUnEquipItem]:
			eduUnEquipItemLogic(e, p)
		case e := <-p.events.recv[eduPosition]:
			eduPositionLogic(e, p)
		}
	}
}

func (p *player) eduPlayerSyncEvents002() {
	log.Infof("[player_worker] eduPlayerSyncEvents002 worker for player %v", p.getHandle())
	for {
		select {
		case e := <-p.events.recv[eduSelectEntity]:
			eduSelectEntityLogic(e, p)
		case e := <-p.events.recv[eduUnselectsEntity]:
			eduUnselectEntityLogic(e, p)
		}
	}
}

// combination of events that must be processed in order
// the reasoning is to avoid certain events changing data at the same time
func (p *player) eduPlayerAsyncEvents() {
	log.Infof("[player_worker] eduPlayerEvents003 worker for player %v", p.getHandle())
	//for {
	//	select {
	//	case e := <-p.events.recv[eduSelectEntityAsync]:
	//		go eduSelectEntityLogic(e, p)
	//	case e := <-p.events.recv[eduUnselectsEntityAsync]:
	//		go eduSelectEntityLogic(e, p)
	//	}
	//}
}

func eduSelectEntityLogic(e event, player *player) {
	ev, ok := e.(*eduSelectEntityEvent)
	if !ok {
		ev.err <- eventTypeCastError(reflect.TypeOf(eduSelectEntityEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	player.selects(ev.entity)
	ev.entity.selectedBy(player)

	ev.err <- nil
}

func eduUnselectEntityLogic(e event, player *player) {
	ev, ok := e.(*eduUnselectEntityEvent)
	if !ok {
		ev.err <- eventTypeCastError(reflect.TypeOf(eduUnselectEntityEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	player.removeSelection()
	ev.err <- nil
}

func eduUnEquipItemLogic(e event, player *player) {
	ev, ok := e.(*eduUnEquipItemEvent)
	if !ok {
		ev.err <- eventTypeCastError(reflect.TypeOf(eduUnEquipItemEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	change, err := player.unEquip(ev.from, ev.to)
	ev.change = change
	ev.err <- err
}

func eduEquipItemLogic(e event, player *player) {
	ev, ok := e.(*eduEquipItemEvent)
	if !ok {
		ev.err <- eventTypeCastError(reflect.TypeOf(eduEquipItemEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	change, err := player.equip(ev.slot)
	ev.change = change
	ev.err <- err
}

func eduPositionLogic(e event, player *player) {
	ev, ok := e.(*eduPositionEvent)
	if !ok {
		ev.err <- eventTypeCastError(reflect.TypeOf(eduPositionEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}
	err := player.move(ev.zm, ev.x, ev.y)
	ev.err <- err
}
