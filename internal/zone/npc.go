package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"sync"
)

type npc struct {
	baseEntity
	data  *npcStaticData
	ticks *entityTicks
	state *entityState
	stats *npcStats
	sync.RWMutex
}

type npcStats struct {
	hp, sp uint32
}

type npcStaticData struct {
	mobInfo       *data.MobInfo
	mobInfoServer *data.MobInfoServer
	regenData     *data.RegenEntry
	npcData       *data.ShineNPC
}

func ncBatTargetInfoCmd(n *npc) *structs.NcBatTargetInfoCmd {
	var nc = structs.NcBatTargetInfoCmd{
		Handle:      n.getHandle(),
		TargetMaxHP: n.data.mobInfo.MaxHP,               //todo: use the same player stat system for mobs and NPCs
		TargetMaxSP: uint32(n.data.mobInfoServer.MaxSP), //todo: use the same player stat system for mobs and NPCs
		TargetLevel: byte(n.data.mobInfo.Level),
	}

	nc.TargetHP = n.stats.hp
	nc.TargetSP = n.stats.sp

	return &nc
}

// find a way to merge npc and monster structs
func ncBriefInfoRegenMobCmd(n *npc) structs.NcBriefInfoRegenMobCmd {
	var nc = structs.NcBriefInfoRegenMobCmd{
		Handle: n.getHandle(),
		Mode:   byte(n.data.mobInfoServer.EnemyDetect),
		MobID:  n.data.mobInfo.ID,
	}
	n.baseEntity.RLock()
	nc.Coord = structs.ShineCoordType{
		XY: structs.ShineXYType{
			X: uint32(n.current.x),
			Y: uint32(n.current.y),
		},
		Direction: uint8(n.current.d),
	}
	n.baseEntity.RUnlock()
	return nc
}
