package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"reflect"
)

// combination of events that must be processed in order
// the reasoning is to avoid certain events changing data at the same time
func (p *player) eduPlayerEvents001() {
	log.Infof("[player_worker] eduPlayerEvents001 worker for player %v", p.getHandle())
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
			//case e := <-p.events.recv[eduUseItem]:
		}
	}
}

// combination of events that must be processed in order
// the reasoning is to avoid certain events changing data at the same time
func (p *player) eduPlayerEvents002() {
	log.Infof("[player_worker] eduPlayerEvents002 worker for player %v", p.getHandle())
	for {
		select {
		case e := <-p.events.recv[eduPosition]:
			eduPositionLogic(e, p)

		}
	}
}

func eduUnEquipItemLogic(e event, player *player) {
	ev, ok := e.(*eduUnEquipItemEvent)
	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneUnexpectedEvent,
			Details: errors.ErrDetails{
				"expected": reflect.TypeOf(eduUnEquipItemEvent{}).String(),
				"actual":   reflect.TypeOf(ev).String(),
			},
		})
		return
	}
	change, err := player.unEquip(ev.from, ev.to)
	ev.change = change
	ev.err <- err
}

func eduEquipItemLogic(e event, player *player) {
	ev, ok := e.(*eduEquipItemEvent)
	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneUnexpectedEvent,
			Details: errors.ErrDetails{
				"expected": reflect.TypeOf(eduEquipItemEvent{}).String(),
				"actual":   reflect.TypeOf(ev).String(),
			},
		})
		return
	}
	change, err := player.equip(ev.slot)
	ev.change = change
	ev.err <- err
}

func eduPositionLogic(e event, player *player) {
	ev, ok := e.(*eduPositionEvent)
	if !ok {
		log.Error(errors.Err{
			Code: errors.ZoneUnexpectedEvent,
			Details: errors.ErrDetails{
				"expected": reflect.TypeOf(eduPositionEvent{}).String(),
				"actual":   reflect.TypeOf(ev).String(),
			},
		})
		return
	}
	err := player.move(ev.zm, ev.x, ev.y)
	ev.err <- err
}
