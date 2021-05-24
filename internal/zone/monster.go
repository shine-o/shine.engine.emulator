package zone

import (
	"time"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

type monster npc

func (m *monster) getTargetPacketData() *structs.NcBatTargetInfoCmd {
	m.targeting.RLock()
	order := m.targeting.selectionOrder
	m.targeting.RUnlock()
	nc := monsterNcBatTargetInfo(m, order)
	return nc
}

func (m *monster) getNextTargetPacketData() *structs.NcBatTargetInfoCmd {
	m.targeting.RLock()
	order := m.targeting.selectionOrder + 1
	m.targeting.RUnlock()
	nc := m.targeting.currentlySelected.getTargetPacketData()
	nc.Order = order
	return nc
}

func (m *monster) selects(e entity) {
	m.targeting.Lock()
	m.targeting.selectionOrder += 32
	m.targeting.currentlySelected = e
	m.targeting.Unlock()
}

func (m *monster) selectedBy(e entity) {
	m.targeting.selectedBy(e)
}

func (m *monster) currentlySelected() entity {
	m.targeting.RLock()
	defer m.targeting.RUnlock()
	return m.targeting.currentlySelected
}

func monsterNcBatTargetInfo(m *monster, assignedOrder byte) *structs.NcBatTargetInfoCmd {
	nc := &structs.NcBatTargetInfoCmd{
		Order:       assignedOrder,
		Handle:      m.getHandle(),
		TargetMaxHP: m.data.mobInfo.MaxHP,               // todo: use the same player stat system for mobs and NPCs
		TargetMaxSP: uint32(m.data.mobInfoServer.MaxSP), // todo: use the same player stat system for mobs and NPCs
		TargetLevel: byte(m.data.mobInfo.Level),
		TargetHP:    m.data.mobInfo.MaxHP,               // for now we use the static value, later it should be switched to stats value
		TargetSP:    uint32(m.data.mobInfoServer.MaxSP), // for now we use the static value, later it should be switched to stats value
	}
	return nc
}

func (m *monster) notifyAboutRemovedEntity(e entity) {
	// panic("implement me")
}

func (m *monster) alreadyNearbyEntity(e entity) bool {
	m.baseEntity.proximity.RLock()
	_, exists := m.baseEntity.proximity.entities[e.getHandle()]
	m.baseEntity.proximity.RUnlock()
	return exists
}

func (m *monster) newNearbyEntitiesTicker(zm *zoneMap) {
	log.Infof("[player_ticks] newNearbyEntitiesTicker for handle %v", m.getHandle())
	tick := time.NewTicker(200 * time.Millisecond)
	m.ticks.Lock()
	m.ticks.list = append(m.ticks.list, tick)
	m.ticks.Unlock()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			newEntities := addWithinRangeEntities(m, zm)
			for _, e := range newEntities {
				go m.notifyAboutNewEntity(e)
			}
		}
	}
}

func (m *monster) oldNearbyEntitiesTicker() {
	log.Infof("[player_ticks] oldNearbyEntitiesTicker for handle %v", m.getHandle())
	tick := time.NewTicker(200 * time.Millisecond)
	m.ticks.Lock()
	m.ticks.list = append(m.ticks.list, tick)
	m.ticks.Unlock()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			removeOutOfRangeEntities(m)
		}
	}
}

func (m *monster) notifyAboutNewEntity(e entity) {
	log.Info("implement me")
}

func (m *monster) getNewEntityNearbyPacketData() interface{} {
	return monsterNcBriefInfoRegenMobCmd(m)
}

// return a buffered channel with all nearby entities
func (m *monster) getNearbyEntities() <-chan entity {
	return getNearbyEntities(m.baseEntity.proximity)
}

func (m *monster) removeNearbyEntity(e entity) {
	m.baseEntity.proximity.Lock()
	delete(m.baseEntity.proximity.entities, e.getHandle())
	m.baseEntity.proximity.Unlock()
}

func (m *monster) addNearbyEntity(e entity) {
	h := e.getHandle()
	m.baseEntity.proximity.Lock()
	m.baseEntity.proximity.entities[h] = e
	m.baseEntity.proximity.Unlock()
}

func (m *monster) spawnLocation(zm *zoneMap) {
	m.Lock()
	var (
		shineD int
		sn     = m.data.npcData
	)

	if sn.D < 0 {
		shineD = (360 + sn.D) / 2
	} else {
		shineD = sn.D / 2
	}

	m.baseEntity.current.mapName = zm.data.MapInfoIndex
	m.baseEntity.current.mapID = zm.data.ID
	m.baseEntity.current.x = sn.X
	m.baseEntity.current.y = sn.Y
	m.baseEntity.current.d = shineD

	m.Unlock()
}
