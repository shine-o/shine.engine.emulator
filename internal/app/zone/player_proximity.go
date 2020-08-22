package zone

import (
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
)

// notify every player in proximity about player that logged in
func newPlayer(p *player, nearbyPlayers *players) {
	nearbyPlayers.Lock()
	for _, np := range nearbyPlayers.active {
		if p.handle == np.handle {
			continue
		}
		p.Lock()
		nc := p.ncBriefInfoLoginCharacterCmd()
		p.Unlock()
		ncBriefInfoLoginCharacterCmd(np, &nc)
	}
	nearbyPlayers.Unlock()
}

// send info to player about nearby players
func nearbyPlayers(p *player, nearbyPlayers *players) {
	nearbyPlayers.Lock()
	var characters []structs.NcBriefInfoLoginCharacterCmd
	for _, np := range nearbyPlayers.active {
		if np.handle == p.handle {
			continue
		}
		p.Lock()
		nc := np.ncBriefInfoLoginCharacterCmd()
		p.Unlock()
		characters = append(characters, nc)
	}
	ncBriefInfoCharacterCmd(p, &structs.NcBriefInfoCharacterCmd{
		Number:     byte(len(characters)),
		Characters: characters,
	})
	nearbyPlayers.Unlock()
}
