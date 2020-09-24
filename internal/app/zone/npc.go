package zone

import (
	mobs "github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/monsters"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/shn"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"sync"
	"time"
)

type npc struct {
	baseEntity
	hp, sp        uint32
	mobInfo       *shn.MobInfo
	mobInfoServer *shn.MobInfoServer
	regenData     *mobs.RegenEntry
	tickers       []*time.Ticker
	status
	sync.RWMutex
}


// find a way to merge npc and monster structs
func (n *npc) ncBriefInfoRegenMobCmd() structs.NcBriefInfoRegenMobCmd {
	n.RLock()
	nc := structs.NcBriefInfoRegenMobCmd{
		Handle: n.handle,
		Mode:   byte(n.mobInfoServer.EnemyDetect),
		MobID:  n.mobInfo.ID,
		Coord: structs.ShineCoordType{
			XY: structs.ShineXYType{
				X: uint32(n.current.x),
				Y: uint32(n.current.y),
			},
			Direction: uint8(n.current.d),
		},
	}
	n.RUnlock()
	return nc
}