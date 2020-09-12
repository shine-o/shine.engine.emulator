package world

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game/character"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

type createCharacterEvent struct {
	nc *structs.NcAvatarCreateReq
	np *networking.Parameters
}

type deleteCharacterEvent struct {
	nc *structs.NcAvatarEraseReq
	np *networking.Parameters
}

type characterLoginEvent struct {
	nc *structs.NcCharLoginReq
	np *networking.Parameters
}

type characterSettingsEvent struct {
	char *character.Character
	np   *networking.Parameters
}

type updateShortcutsEvent struct {
	np          *networking.Parameters
	nc          structs.NcCharOptionSetShortcutDataReq
	characterID uint64
}

type characterSelectEvent struct {
	np      *networking.Parameters
	session *session
}
