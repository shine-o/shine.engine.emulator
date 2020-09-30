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

func (n *npc) ncBatTargetInfoCmd() *structs.NcBatTargetInfoCmd {
	var nc structs.NcBatTargetInfoCmd
	n.RLock()
	nc = structs.NcBatTargetInfoCmd{
		Order:         0,
		Handle:        n.handle,
		TargetHP:      n.hp,
		TargetMaxHP:   n.mobInfo.MaxHP, //todo: use the same player stat system for mobs and NPCs
		TargetSP:      n.sp,
		TargetMaxSP:   uint32(n.mobInfoServer.MaxSP), //todo: use the same player stat system for mobs and NPCs
		TargetLP:      0,
		TargetMaxLP:   0,
		TargetLevel:   byte(n.mobInfo.Level),
		HpChangeOrder: 0,
	}
	n.RUnlock()
	return &nc
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