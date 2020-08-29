package world

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game/character"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

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
