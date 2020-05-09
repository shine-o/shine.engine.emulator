package service

import (
	"github.com/shine-o/shine.engine.core/structs"
)

// notify every player in proximity about player that logged in
func newPlayer(p *player, nearbyPlayers map[uint16]*player) {
	for _, np := range nearbyPlayers {
		if p.handle == np.handle {
			continue
		}
		//ncBriefInfoLoginCharacterCmd(p, &ev.nc)
	}
}

// send info to player about nearby players
func nearbyPlayers(p *player, nearbyPlayers map[uint16]*player) {
	var characters []structs.NcBriefInfoLoginCharacterCmd
	for _, np := range nearbyPlayers {
		if np.handle == p.handle {
			continue
		}
		nc := np.ncBriefInfoLoginCharacterCmd()
		characters = append(characters, nc)
	}
	ncBriefInfoCharacterCmd(p, &structs.NcBriefInfoCharacterCmd{
		Number:     byte(len(characters)),
		Characters: characters,
	})
}
