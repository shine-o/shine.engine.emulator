package world

import (
	"context"
	zm "github.com/shine-o/shine.engine.emulator/internal/pkg/grpc/zone-master"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/persistence"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"reflect"
)

func (w *world) characterCRUD() {
	log.Info("[character_worker] characterCRUD worker")
	for {
		select {
		case e := <-w.recv[createCharacter]:
			go createCharacterLogic(e, w)
		case e := <-w.recv[deleteCharacter]:
			go deleteCharacterLogic(e, w)
		}
	}
}

func (w *world) characterSession() {
	log.Info("[character_worker] characterSession worker")
	for {
		select {
		case e := <-w.recv[characterLogin]:
			go characterLoginLogic(e, w)
		case e := <-w.recv[characterSettings]:
			go characterSettingsLogic(e)
		case e := <-w.recv[updateShortcuts]:
			go updateShortcutsLogic(w, e)
		case e := <-w.recv[updateGameSettings]:
			log.Info(e)
		case e := <-w.recv[updateKeymap]:
			log.Info(e)
		case e := <-w.recv[characterSelect]:
			go characterSelectLogic(e, w)
		}
	}
}

func characterSelectLogic(e event, w *world) {
	ev, ok := e.(*characterSelectEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&characterSelectEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	nc, err := userCharacters(ev.session)
	if err != nil {
		log.Error(err)
		return
	}

	//networking.Send(ev.np.OutboundSegments.Send, networking.NC_USER_LOGOUT_DB, &nc)
	networking.Send(ev.np.OutboundSegments.Send, networking.NC_USER_LOGINWORLD_ACK, &nc)
}

func characterLoginLogic(e event, w *world) {
	ev, ok := e.(*characterLoginEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&characterLoginEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	s, ok := ev.np.Session.(*session)
	if !ok {
		log.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
		return
	}

	char, err := persistence.GetBySlot(ev.nc.Slot, s.UserID)
	if err != nil {
		log.Error(err)
		return
	}

	nc, err := zoneConnectionInfo(char)
	if err != nil {
		log.Error(err)
		return
	}

	networking.Send(ev.np.OutboundSegments.Send, networking.NC_CHAR_LOGIN_ACK, &nc)

	go worldTimeNotification(ev.np)

	session, ok := ev.np.Session.(*session)

	if !ok {
		log.Error("no session available")
		return
	}

	session.Lock()
	session.characterID = char.ID
	session.Unlock()

	cs := characterSettingsEvent{
		char: &char,
		np:   ev.np,
	}

	worldEvents[characterSettings] <- &cs
}

func zoneConnectionInfo(char persistence.Character) (structs.NcCharLoginAck, error) {
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

func createCharacterLogic(e event, w *world) {
	ev, ok := e.(*createCharacterEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&createCharacterEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	s, ok := ev.np.Session.(*session)

	if !ok {
		log.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
	}

	err := persistence.Validate(s.UserID, ev.nc)

	if err != nil {
		log.Error(err)
		ncAvatarCreateFailAck(ev.np, 385)

		return
	}

	char, err := persistence.New(s.UserID, ev.nc)

	if err != nil {
		log.Error(err)
		ncAvatarCreateFailAck(ev.np, 385)
		return
	}

	nc := structs.NcAvatarCreateSuccAck{
		NumOfAvatar: 1,
		Avatar:      char.NcRepresentation(),
	}

	networking.Send(ev.np.OutboundSegments.Send, networking.NC_AVATAR_CREATESUCC_ACK, &nc)
}

func deleteCharacterLogic(e event, w *world) {
	ev, ok := e.(*deleteCharacterEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&deleteCharacterEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	s, ok := ev.np.Session.(*session)
	if !ok {
		log.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
	}

	err := persistence.Delete(s.UserID, ev.nc)
	if err != nil {
		log.Error(err)
		return
	}

	nc := &structs.NcAvatarEraseSuccAck{
		Slot: ev.nc.Slot,
	}
	networking.Send(ev.np.OutboundSegments.Send, networking.NC_AVATAR_ERASESUCC_ACK, nc)

}

func characterSettingsLogic(e event) {
	ev, ok := e.(*characterSettingsEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&characterSettingsEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	gameOptions, err := persistence.NcGameOptions(ev.char.Options.GameOptions)

	if err != nil {
		log.Error(err)
		return
	}

	keyMap, err := persistence.NcKeyMap(ev.char.Options.Keymap)

	if err != nil {
		log.Error(err)
		return
	}

	shortcuts, err := persistence.NcShortcutData(ev.char.Options.Shortcuts)

	if err != nil {
		log.Error(err)
		return
	}

	networking.Send(ev.np.OutboundSegments.Send, networking.NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD, &gameOptions)
	networking.Send(ev.np.OutboundSegments.Send, networking.NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD, &keyMap)
	networking.Send(ev.np.OutboundSegments.Send, networking.NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD, &shortcuts)

}

func updateShortcutsLogic(w *world, e event) {
	ev, ok := e.(*updateShortcutsEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&updateShortcutsEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	c, err := persistence.Get(ev.characterID)

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

	err = persistence.Update(&c)
	if err != nil {
		log.Error(err)
		return
	}

	nc := structs.NcCharOptionImproveShortcutDataAck{ErrCode: 8448}

	networking.Send(ev.np.OutboundSegments.Send, networking.NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_ACK, &nc)
}
