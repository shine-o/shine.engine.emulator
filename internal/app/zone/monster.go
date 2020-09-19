package zone

import (
	mobs "github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/monsters"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/shn"
)

type monster struct {
	baseEntity
	hp, sp uint32

	mobInfo * shn.MobInfo
	mobInfoServer * shn.MobInfoServer
	regenData  mobs.RegenEntry
}


func (m * monster) alive()  {

	// tick for roaming around, idle action, return to base

}

func (m * monster) dead() {

	// initial removal from monsters...
	// trigger monsterDied event
	// create ticker so it can respawn again
}

func (m * monster) roam() {

	// initial removal from monsters...
	// trigger monsterDied event
	// create ticker so it can respawn again
}