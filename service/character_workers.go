package service

import (
	"context"
	"fmt"
	zm "github.com/shine-o/shine.engine.core/grpc/zone-master"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/shine-o/shine.engine.world/service/character"
	"reflect"
)

func (w *world) characterCRUD() {
	log.Info("[character_worker] characterCRUD worker)")
	for {
		select {
		case e := <- w.recv[createCharacter]:
			go func() {
				ev, ok := e.(*createCharacterEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&createCharacterEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				s, ok := ev.np.Session.(*session)
				if !ok {
					ev.err <- fmt.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
				}
				err := character.Validate(w.db, s.UserID, ev.nc)
				if err != nil {
					ev.err <- err
					return
				}
				char, err := character.New(db, s.UserID, ev.nc)
				if err != nil {
					ev.err <- err
					return
				}
				ev.char <- char
			}()
		case e := <- w.recv[deleteCharacter]:
			go func() {
				ev, ok := e.(*deleteCharacterEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&deleteCharacterEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				s, ok := ev.np.Session.(*session)
				if !ok {
					ev.err <- fmt.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
				}
				err := character.Delete(db, s.UserID, ev.nc)
				if err != nil {
					ev.err <- err
					return
				}
				ev.done <- true
			}()
		}
	}
}

func (w *world) characterSession() {
	log.Info("[character_worker] characterSession worker)")
	for {
		select {
		case e := <- w.recv[characterLogin]:
			go func() {
				var cs characterSettingsEvent
				ev, ok := e.(*characterLoginEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&characterLoginEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				s, ok := ev.np.Session.(*session)
				if !ok {
					ev.err <- fmt.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
				}

				char, err := character.GetBySlot(w.db, ev.nc.Slot, s.UserID)

				if err != nil {
					ev.err <- err
					return
				}

				nc, err := zoneConnectionInfo(err, ev, char)
				if err != nil {
					ev.err <- err
					return
				}
				cs = characterSettingsEvent{
					char: &char,
					np:   ev.np,
				}
				ev.zoneInfo <- &nc
				worldEvents[characterSettings] <- &cs
			}()
		case e := <-w.recv[characterSettings]:
			go func() {
				ev, ok := e.(*characterSettingsEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&characterSettingsEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				gameOptions, err := character.NcGameOptions(ev.char.Options.GameOptions)
				if err != nil {
					log.Error(err)
					return
				}
				keyMap, err := character.NcKeyMap(ev.char.Options.Keymap)
				if err != nil {
					log.Error(err)
					return
				}
				shortcuts, err := character.NcShortcutData(ev.char.Options.Shortcuts)
				if err != nil {
					log.Error(err)
					return
				}
				ncCharOptionImproveGetGameOptionCmd(ev.np, &gameOptions)
				ncCharOptionImproveGetKeymapCmd(ev.np, &keyMap)
				ncCharOptionImproveGetShortcutDataCmd(ev.np, &shortcuts)
			}()
		}
	}
}

func zoneConnectionInfo(err error, ev *characterLoginEvent, char character.Character) (structs.NcCharLoginAck, error) {
	var nc structs.NcCharLoginAck
	conn, err := newRPCClient("zone_master")
	if err != nil {
		return nc, err
	}

	c := zm.NewMasterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), gRPCTimeout)
	defer cancel()
	ci, err := c.WhereIsMap(ctx, &zm.MapQuery{
		ID: int32(char.Location.MapID),
	})

	if err != nil {
		ev.err <- err
		return nc, err
	}

	nc = structs.NcCharLoginAck{
		ZoneIP: structs.Name4{
			Name: ci.IP,
		},
		ZonePort: uint16(ci.Port),
	}
	return nc, nil
}
