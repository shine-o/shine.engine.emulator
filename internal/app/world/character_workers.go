package world

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game/character"
	zm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/zone-master"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
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
				s, ok := ev.np.Session.(*session)
				if !ok {
					log.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
				}
				err := character.Validate(w.db, s.UserID, ev.nc)
				if err != nil {
					log.Error(err)
					createCharErr(ev.np, err)
					return
				}

				char, err := character.New(w.db, s.UserID, ev.nc)

				if err != nil {
					log.Error(err)
					createCharErr(ev.np, err)
					return
				}

				nc := structs.NcAvatarCreateSuccAck{
					NumOfAvatar: 1,
					Avatar:      char.NcRepresentation(),
				}
				ncAvatarCreateSuccAck(ev.np, &nc)
			}()
		case e := <-w.recv[deleteCharacter]:
			go func() {
				ev, ok := e.(*deleteCharacterEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&deleteCharacterEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}

				s, ok := ev.np.Session.(*session)
				if !ok {
					log.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
				}

				err := character.Delete(w.db, s.UserID, ev.nc)
				if err != nil {
					log.Error(err)
					return
				}

				avatarEraseSuccAck(ev.np, &structs.NcAvatarEraseSuccAck{
					Slot: ev.nc.Slot,
				})
			}()
		}
	}
}

func (w *world) characterSession() {
	log.Info("[character_worker] characterSession worker")
	for {
		select {
		case e := <-w.recv[characterLogin]:
			ev, ok := e.(*characterLoginEvent)

			if !ok {
				log.Errorf("expected event type %v but got %v", reflect.TypeOf(&characterLoginEvent{}).String(), reflect.TypeOf(ev).String())
				return
			}

			go handleCharacterLogin(ev, w)

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
		case e := <-w.recv[updateShortcuts]:
			go func() {
				ev, ok := e.(*updateShortcutsEvent)
				if !ok {
					log.Errorf("expected event type %v but got %v", reflect.TypeOf(&updateShortcutsEvent{}).String(), reflect.TypeOf(ev).String())
					return
				}
				c, err := character.Get(w.db, ev.characterID)
				if err != nil {
					log.Error(err)
					return
				}

				// first load the existing data into a  structs.NcCharGetShortcutDataCmd struct

				// iterate over the shortcuts

				storageNC := structs.NcCharGetShortcutDataCmd{ // we need this struct for storage, since the client needs this struct for logging the character
					Count:     uint16(ev.nc.Count),
					Shortcuts: ev.nc.Shortcuts,
				}

				data, err := structs.Pack(&storageNC)
				if err != nil {
					log.Error(err)
					return
				}

				c.Options.Shortcuts = data

				err = character.Update(w.db, &c)
				if err != nil {
					log.Error(err)
					return
				}

				nc := structs.NcCharOptionImproveShortcutDataAck{ErrCode: 8448}
				ncCharOptionImproveSetShortcutDataAck(ev.np, &nc)

			}()
		//case e := <-w.recv[updateGameSettings]:
		//case e := <-w.recv[updateKeymap]:

		}
	}
}

func handleCharacterLogin(ev *characterLoginEvent, w *world) {
	s, ok := ev.np.Session.(*session)
	if !ok {
		log.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
		return
	}

	char, err := character.GetBySlot(w.db, ev.nc.Slot, s.UserID)
	if err != nil {
		log.Error(err)
		return
	}

	nc, err := zoneConnectionInfo(char)
	if err != nil {
		log.Error(err)
		return
	}

	ncCharLoginAck(ev.np, &nc)

	session, ok := ev.np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	session.characterID = char.ID

	cs := characterSettingsEvent{
		char: &char,
		np:   ev.np,
	}

	worldEvents[characterSettings] <- &cs
}

func zoneConnectionInfo(char character.Character) (structs.NcCharLoginAck, error) {
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
