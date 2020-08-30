package world

import (
	"context"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game/character"
	zm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/zone-master"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"reflect"
)

func characterLoginLogic(ev *characterLoginEvent, w *world) {
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

func createCharacterLogic(ev *createCharacterEvent, w *world) {
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
}

func deleteCharacterLogic(ok bool, ev *deleteCharacterEvent, w *world) {
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
}

func characterSettingsLogic(ev *characterSettingsEvent) {
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
}

func updateShortcutsLogic(w *world, ev *updateShortcutsEvent) {
	c, err := character.Get(w.db, ev.characterID)

	if err != nil {
		log.Error(err)
		return
	}

	storedShortcuts := structs.NcCharGetShortcutDataCmd{}

	err = structs.Unpack(c.Options.Shortcuts, &storedShortcuts)
	if err != nil {
		log.Error(err)
		return
	}

	var newShortcuts []structs.ShortCutData

	for _, s1 := range ev.nc.Shortcuts {
		exists := false
		for j, s2 := range storedShortcuts.Shortcuts {
			if s2.SlotNo == s1.SlotNo {
				storedShortcuts.Shortcuts[j].CodeNo = s1.CodeNo
				storedShortcuts.Shortcuts[j].Value = s1.Value
				exists = true
			}
		}
		if !exists {
			newShortcuts = append(newShortcuts, s1)
		}
	}

	storedShortcuts.Shortcuts = append(storedShortcuts.Shortcuts, newShortcuts...)

	storedShortcuts.Count = uint16(len(storedShortcuts.Shortcuts))

	data, err := structs.Pack(&storedShortcuts)

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
}
