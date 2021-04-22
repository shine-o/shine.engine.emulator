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
			log.Info(e)
		case e := <-p.events.recv[eduUnEquipItem]:
			log.Info(e)

		}
	}
}

// combination of events that must be processed in order
func (p *player) eduPlayerEvents002() {
	log.Infof("[player_worker] eduPlayerEvents002 worker for player %v", p.getHandle())
	for {
		select {
		case e := <-p.events.recv[eduPosition]:
			updatePlayerPosition(e)
			//case e := <-p.events.recv[eduUseItem]:

		}
	}
}

func updatePlayerPosition(e event) {
	ev, ok := e.(*eduPositionEvent)
	if !ok {
		log.Error(errors.Err{
			Code:    errors.ZoneUnexpectedEvent,
			Details: errors.ErrDetails{
				"expected": reflect.TypeOf(eduPositionEvent{}).String(),
				"actual":  reflect.TypeOf(ev).String(),
			},
		})
		return
	}
	err := ev.player.move(ev.zm, ev.x, ev.y)
	ev.err <- err
}