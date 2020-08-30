package world

import (
	"reflect"
)

func (w *world) characterCRUD() {
	log.Info("[character_worker] characterCRUD worker")
	for {
		select {
		case e := <-w.recv[createCharacter]:
			go func() {
				ev, ok := e.(*createCharacterEvent)

				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&createCharacterEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				createCharacterLogic(ev, w)
			}()
		case e := <-w.recv[deleteCharacter]:
			go func() {
				ev, ok := e.(*deleteCharacterEvent)

				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&deleteCharacterEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				deleteCharacterLogic(ok, ev, w)
			}()
		}
	}
}

func (w *world) characterSession() {
	log.Info("[character_worker] characterSession worker")
	for {
		select {
		case e := <-w.recv[characterLogin]:
			go func() {
				ev, ok := e.(*characterLoginEvent)

				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&characterLoginEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				characterLoginLogic(ev, w)
			}()

		case e := <-w.recv[characterSettings]:
			go func() {
				ev, ok := e.(*characterSettingsEvent)

				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&characterSettingsEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				characterSettingsLogic(ev)
			}()
		case e := <-w.recv[updateShortcuts]:
			go func() {
				ev, ok := e.(*updateShortcutsEvent)

				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&updateShortcutsEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				updateShortcutsLogic(w, ev)
			}()
		case e := <-w.recv[updateGameSettings]:
			log.Info(e)
		case e := <-w.recv[updateKeymap]:
			log.Info(e)
		}
	}
}
