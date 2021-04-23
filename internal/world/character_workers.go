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
			go createCharacterLogic(e)
		case e := <-w.recv[deleteCharacter]:
			go deleteCharacterLogic(e)
		}
	}
}

func (w *world) characterSession() {
	log.Info("[character_worker] characterSession worker")
	for {
		select {
		case e := <-w.recv[characterLogin]:
			go characterLoginLogic(e)
		case e := <-w.recv[characterSettings]:
			go characterSettingsLogic(e)
		case e := <-w.recv[updateShortcuts]:
			go updateShortcutsLogic(e)
		case e := <-w.recv[updateGameSettings]:
			log.Info(e)
		case e := <-w.recv[updateKeymap]:
			log.Info(e)
		case e := <-w.recv[characterSelect]:
			go characterSelectLogic(e)
		}
	}
}

func createCharacterLogic(e event) {
	ev, ok := e.(*createCharacterEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&createCharacterEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	s, ok := ev.np.Session.(*session)

	if !ok {
		log.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
	}

	err := persistence.ValidateCharacter(s.UserID, ev.nc)

	if err != nil {
		log.Error(err)
		ncAvatarCreateFailAck(ev.np, 385)

		return
	}

	char, err := persistence.NewCharacter(s.UserID, ev.nc, true)

	if err != nil {
		log.Error(err)
		ncAvatarCreateFailAck(ev.np, 385)
		return
	}

	nc := structs.NcAvatarCreateSuccAck{
		NumOfAvatar: 1,
		Avatar:      avatarInformation(char),
	}

	networking.Send(ev.np.OutboundSegments.Send, networking.NC_AVATAR_CREATESUCC_ACK, &nc)
}

func deleteCharacterLogic(e event) {
	ev, ok := e.(*deleteCharacterEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&deleteCharacterEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	s, ok := ev.np.Session.(*session)
	if !ok {
		log.Errorf("failed to cast given session %v to world session %v", reflect.TypeOf(ev.np.Session).String(), reflect.TypeOf(&session{}).String())
	}

	err := persistence.DeleteCharacter(s.UserID, int(ev.nc.Slot))
	if err != nil {
		log.Error(err)
		return
	}

	nc := &structs.NcAvatarEraseSuccAck{
		Slot: ev.nc.Slot,
	}
	networking.Send(ev.np.OutboundSegments.Send, networking.NC_AVATAR_ERASESUCC_ACK, nc)

}

func characterLoginLogic(e event) {
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

	char, err := persistence.GetCharacterBySlot(ev.nc.Slot, s.UserID)
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

func characterSettingsLogic(e event) {
	ev, ok := e.(*characterSettingsEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&characterSettingsEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	gameOptions, err := ncGameOptions(ev.char.Options.GameOptions)

	if err != nil {
		log.Error(err)
		return
	}

	keyMap, err := ncKeyMap(ev.char.Options.Keymap)

	if err != nil {
		log.Error(err)
		return
	}

	shortcuts, err := ncShortcutData(ev.char.Options.Shortcuts)

	if err != nil {
		log.Error(err)
		return
	}

	networking.Send(ev.np.OutboundSegments.Send, networking.NC_CHAR_OPTION_IMPROVE_GET_GAMEOPTION_CMD, &gameOptions)
	networking.Send(ev.np.OutboundSegments.Send, networking.NC_CHAR_OPTION_IMPROVE_GET_KEYMAP_CMD, &keyMap)
	networking.Send(ev.np.OutboundSegments.Send, networking.NC_CHAR_OPTION_IMPROVE_GET_SHORTCUTDATA_CMD, &shortcuts)

}

func updateShortcutsLogic(e event) {
	ev, ok := e.(*updateShortcutsEvent)

	if !ok {
		log.Errorf("expected event type %v but got %v", reflect.TypeOf(&updateShortcutsEvent{}).String(), reflect.TypeOf(ev).String())
		return
	}

	c, err := persistence.GetCharacter(ev.characterID)

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

	err = persistence.UpdateCharacter(&c)
	if err != nil {
		log.Error(err)
		return
	}

	nc := structs.NcCharOptionImproveShortcutDataAck{ErrCode: 8448}

	networking.Send(ev.np.OutboundSegments.Send, networking.NC_CHAR_OPTION_IMPROVE_SET_SHORTCUTDATA_ACK, &nc)
}

func characterSelectLogic(e event) {
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

// NcRepresentation returns a struct that can be serialized into bytes and can be sent to the client
func avatarInformation(c *persistence.Character) structs.AvatarInformation {
	nc := structs.AvatarInformation{
		ChrRegNum: uint32(c.ID),
		Name: structs.Name5{
			Name: c.Name,
		},
		Level: uint16(c.Attributes.Level),
		Slot:  c.Slot,
		LoginMap: structs.Name3{
			Name: c.Location.MapName,
		},
		DelInfo: structs.ProtoAvatarDeleteInfo{},
		Shape:   protoAvatarShapeInfo(c.Appearance),
		Equip:   protoEquipment(c.EquippedItems),
		TutorialInfo: structs.ProtoTutorialInfo{ // x(
			TutorialState: 2,
			TutorialStep:  byte(0),
		},
	}
	return nc
}

// NcRepresentation returns a struct that can be serialized into bytes and can be sent to the client
func protoEquipment(cei *persistence.EquippedItems) structs.ProtoEquipment {
	return structs.ProtoEquipment{
		EquHead:         cei.Head,
		EquMouth:        cei.ApparelFace,
		EquRightHand:    cei.RightHand,
		EquBody:         cei.Body,
		EquLeftHand:     cei.LeftHand,
		EquPant:         cei.Pants,
		EquBoot:         cei.Boots,
		EquAccBoot:      cei.ApparelBoots,
		EquAccPant:      cei.ApparelPants,
		EquAccBody:      cei.ApparelBody,
		EquAccHeadA:     cei.ApparelHead,
		EquMinimonR:     cei.RightMiniPet,
		EquEye:          cei.Face,
		EquAccLeftHand:  cei.ApparelLeftHand,
		EquAccRightHand: cei.ApparelRightHand,
		EquAccBack:      cei.ApparelBack,
		EquCosEff:       cei.ApparelAura,
		EquAccHip:       cei.ApparelTail,
		EquMinimon:      cei.LeftMiniPet,
		EquAccShield:    cei.ApparelShield,
		Upgrade:         structs.EquipmentUpgrade{},
	}
}

// NcRepresentation returns a struct that can be serialized into bytes and can be sent to the client
func protoAvatarShapeInfo(ca *persistence.Appearance) structs.ProtoAvatarShapeInfo {
	return structs.ProtoAvatarShapeInfo{
		BF:        1 | ca.Class<<2 | ca.Gender<<7,
		HairType:  ca.HairType,
		HairColor: ca.HairColor,
		FaceShape: ca.FaceType,
	}
}

func ncGameOptions(data []byte) (structs.NcCharOptionImproveGetGameOptionCmd, error) {
	nc := structs.NcCharOptionImproveGetGameOptionCmd{}
	err := structs.Unpack(data, &nc)
	if err != nil {
		return nc, err
	}
	return nc, nil
}

func ncKeyMap(data []byte) (structs.NcCharGetKeyMapCmd, error) {
	nc := structs.NcCharGetKeyMapCmd{}
	err := structs.Unpack(data, &nc)
	if err != nil {
		return nc, err
	}
	return nc, nil
}

func ncShortcutData(data []byte) (structs.NcCharGetShortcutDataCmd, error) {
	nc := structs.NcCharGetShortcutDataCmd{}
	err := structs.Unpack(data, &nc)
	if err != nil {
		return nc, err
	}
	return nc, nil
}
