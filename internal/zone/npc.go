package zone

import (
	"sync"
	"time"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
)

type npcType int

const (
	npcNoRole npcType = iota
	npcPortal
	npcItemMerchant
	npcSkillMerchant
	npcWeaponMerchant
	npcWeaponTitleMerchant
	npcSoulStoneMerchant
	npcDeposit
	npcQuest
	npcCasino
	npcBijouAnvil
	npcWeaponLicenceMerchant
	npcGuildManager
	npcGuildMerchant
	npcSlotMachine
	npcCoinExchangeMerchant
	npcTownChief
)

type npc struct {
	*baseEntity
	data      *npcStaticData
	ticks     *entityTicks
	state     *entityState
	stats     *npcStats
	targeting *targeting
	nType     npcType
	sync.RWMutex
}

func (n *npc) unselectedBy(e entity) {
	n.targeting.unselectedBy(e)
}

func (n *npc) removeSelection() {
	n.targeting.Lock()
	n.targeting.currentlySelected = nil
	n.targeting.Unlock()
}

func (n *npc) getTargetPacketData() *structs.NcBatTargetInfoCmd {
	n.targeting.RLock()
	order := n.targeting.selectionOrder
	n.targeting.RUnlock()
	nc := npcNcBatTargetInfo(n, order)
	return nc
}

func (n *npc) getNextTargetPacketData() *structs.NcBatTargetInfoCmd {
	return n.targeting.currentlySelected.getTargetPacketData()
}

func (n *npc) selects(e entity) {
	n.targeting.Lock()
	n.targeting.selectionOrder += 32
	n.targeting.currentlySelected = e
	n.targeting.Unlock()
}

func (n *npc) selectedBy(e entity) {
	n.targeting.selectedBy(e)
}

func (n *npc) selected() entity {
	n.targeting.RLock()
	defer n.targeting.RUnlock()
	return n.targeting.currentlySelected
}

func npcNcBatTargetInfo(n *npc, assignedOrder byte) *structs.NcBatTargetInfoCmd {
	nc := &structs.NcBatTargetInfoCmd{
		Order:       assignedOrder,
		Handle:      n.getHandle(),
		TargetMaxHP: n.data.mobInfo.MaxHP,               // todo: use the same player stat system for mobs and NPCs
		TargetMaxSP: uint32(n.data.mobInfoServer.MaxSP), // todo: use the same player stat system for mobs and NPCs
		TargetLevel: byte(n.data.mobInfo.Level),
		TargetHP:    n.data.mobInfo.MaxHP,               // for now we use the static value, later it should be switched to stats value
		TargetSP:    uint32(n.data.mobInfoServer.MaxSP), // for now we use the static value, later it should be switched to stats value
	}
	// TODO: dynamically add this data in the SHN files, NPC look ugly with no SP/HP or the star level
	if n.nType != npcNoRole {
		nc.TargetMaxHP = 1000
		nc.TargetHP = 1000
		nc.TargetMaxSP = 1000
		nc.TargetSP = 1000
		nc.TargetLevel = 200
	}

	return nc
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

func (n *npc) notifyAboutRemovedEntity(e entity) {
	// panic("implement me")
}

func (n *npc) alreadyNearbyEntity(e entity) bool {
	n.baseEntity.proximity.RLock()
	_, exists := n.baseEntity.proximity.entities[e.getHandle()]
	n.baseEntity.proximity.RUnlock()
	return exists
}

func (n *npc) newNearbyEntitiesTicker(zm *zoneMap) {
	log.Infof("[player_ticks] newNearbyEntitiesTicker for handle %v", n.getHandle())
	tick := time.NewTicker(200 * time.Millisecond)
	n.ticks.Lock()
	n.ticks.list = append(n.ticks.list, tick)
	n.ticks.Unlock()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			newEntities := addWithinRangeEntities(n, zm)
			for _, e := range newEntities {
				go n.notifyAboutNewEntity(e)
			}
		}
	}
}

func (n *npc) oldNearbyEntitiesTicker() {
	log.Infof("[player_ticks] oldNearbyEntitiesTicker for handle %v", n.getHandle())
	tick := time.NewTicker(200 * time.Millisecond)
	n.ticks.Lock()
	n.ticks.list = append(n.ticks.list, tick)
	n.ticks.Unlock()
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			removeOutOfRangeEntities(n)
		}
	}
}

func (n *npc) notifyAboutNewEntity(e entity) {
	log.Info("implement me")
}

func (n *npc) getNewEntityNearbyPacketData() interface{} {
	return npcNcBriefInfoRegenMobCmd(n)
}

// return a buffered channel with all nearby entities
func (n *npc) getNearbyEntities() <-chan entity {
	return getNearbyEntities(n.baseEntity.proximity)
}

func (n *npc) removeNearbyEntity(e entity) {
	n.baseEntity.proximity.Lock()
	delete(n.baseEntity.proximity.entities, e.getHandle())
	n.baseEntity.proximity.Unlock()
}

func (n *npc) addNearbyEntity(e entity) {
	h := e.getHandle()
	n.baseEntity.proximity.Lock()
	n.baseEntity.proximity.entities[h] = e
	n.baseEntity.proximity.Unlock()
}

func (n *npc) spawnLocation(zm *zoneMap) {
	n.Lock()
	var (
		shineD int
		sn     = n.data.npcData
	)

	if sn.D < 0 {
		shineD = (360 + sn.D) / 2
	} else {
		shineD = sn.D / 2
	}

	n.baseEntity.current.mapName = zm.data.MapInfoIndex
	n.baseEntity.current.mapID = zm.data.ID
	n.baseEntity.current.x = sn.X
	n.baseEntity.current.y = sn.Y
	n.baseEntity.current.d = shineD

	n.Unlock()
}

func loadBaseNpc(inxName string, eType entityType) (*npc, error) {
	var (
		mi  *data.MobInfo
		mis *data.MobInfoServer
	)

	mi, mis = getNpcData(inxName)

	if mi == nil || mis == nil {
		return nil, errors.Err{
			Code: errors.ZoneMissingNpcData,
			Details: errors.Details{
				"mobIndex": inxName,
			},
		}
	}

	var nType npcType

	n := &npc{
		nType: nType,
		baseEntity: &baseEntity{
			eType: eType,
			proximity: &entityProximity{
				entities: make(map[uint16]entity),
			},
		},
		stats: &npcStats{
			hp: mi.MaxHP,
			sp: uint32(mis.MaxSP),
		},
		data: &npcStaticData{
			mobInfo:       mi,
			mobInfoServer: mis,
		},
		targeting: &targeting{
			players:  make(map[uint16]*player),
			monsters: make(map[uint16]*monster),
			npc:      make(map[uint16]*npc),
		},
		state: &entityState{
			idling:   make(chan bool),
			fighting: make(chan bool),
			chasing:  make(chan bool),
			fleeing:  make(chan bool),
		},
		ticks: &entityTicks{},
	}

	return n, nil
}

func getNpcType(role, arg string) npcType {
	if role == "Gate" || role == "IDGate" || role == "ModeIDGate" {
		return npcPortal
	} else if role == "StoreManager" {
		return npcDeposit
	} else if role == "RandomGate" {
		return npcCasino
	} else if role == "ClientMenu" {
		return npcTownChief
	} else {
		switch role + arg {
		case "MerchantSoulStone":
			return npcSoulStoneMerchant
		case "MerchantWeapon":
			return npcWeaponMerchant
		case "MerchantSkill":
			return npcSkillMerchant
		case "MerchantItem":
			return npcItemMerchant
		case "QuestNpcQuest":
			return npcQuest
		case "GuardQuest":
			return npcQuest
		case "NPCMenuRandomOption":
			return npcBijouAnvil
		case "MerchantWeaponTitle":
			return npcWeaponLicenceMerchant
		case "NPCMenuGuild":
			return npcGuildManager
		case "NPCMenuExchangeCoin":
			return npcCoinExchangeMerchant
		case "QuestNpcGBDice":
			return npcSlotMachine
		case "MerchantGuild":
			return npcGuildMerchant
		default:
			log.Error(errors.Err{
				Code: errors.ZoneUnknownNpcRole,
				Details: errors.Details{
					"role+arg": role + arg,
				},
			})
			return npcNoRole
		}
	}
}

// TODO: create a wrapping struct for this, as more monster shine files will be loaded
func getNpcData(mobIndex string) (*data.MobInfo, *data.MobInfoServer) {
	var (
		mi  *data.MobInfo
		mis *data.MobInfoServer
		// sn  *data.ShineNPC
	)

	for i, row := range monsterData.MobInfo.ShineRow {
		if row.InxName == mobIndex {
			mi = &monsterData.MobInfo.ShineRow[i]
		}
	}

	for i, row := range monsterData.MobInfoServer.ShineRow {
		if row.InxName == mobIndex {
			mis = &monsterData.MobInfoServer.ShineRow[i]
		}
	}

	return mi, mis
}

func npcNcBatTargetInfoCmd(n *npc) *structs.NcBatTargetInfoCmd {
	nc := structs.NcBatTargetInfoCmd{
		Handle:      n.getHandle(),
		TargetMaxHP: n.data.mobInfo.MaxHP,               // todo: use the same player stat system for mobs and NPCs
		TargetMaxSP: uint32(n.data.mobInfoServer.MaxSP), // todo: use the same player stat system for mobs and NPCs
		TargetLevel: byte(n.data.mobInfo.Level),
	}

	nc.TargetHP = n.stats.hp
	nc.TargetSP = n.stats.sp

	return &nc
}

func monsterNcBatTargetInfoCmd(m *monster) *structs.NcBatTargetInfoCmd {
	nc := structs.NcBatTargetInfoCmd{
		Handle:      m.getHandle(),
		TargetMaxHP: m.data.mobInfo.MaxHP,               // todo: use the same player stat system for mobs and NPCs
		TargetMaxSP: uint32(m.data.mobInfoServer.MaxSP), // todo: use the same player stat system for mobs and NPCs
		TargetLevel: byte(m.data.mobInfo.Level),
	}

	nc.TargetHP = m.stats.hp
	nc.TargetSP = m.stats.sp

	return &nc
}

// find a way to merge npc and monster structs
func npcNcBriefInfoRegenMobCmd(n *npc) *structs.NcBriefInfoRegenMobCmd {
	nc := &structs.NcBriefInfoRegenMobCmd{
		Handle: n.getHandle(),
		Mode:   byte(n.data.mobInfoServer.EnemyDetect),
		MobID:  n.data.mobInfo.ID,
		// AnimationLevel: 2,
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

func monsterNcBriefInfoRegenMobCmd(m *monster) *structs.NcBriefInfoRegenMobCmd {
	nc := &structs.NcBriefInfoRegenMobCmd{
		Handle: m.getHandle(),
		Mode:   byte(m.data.mobInfoServer.EnemyDetect),
		MobID:  m.data.mobInfo.ID,
		// AnimationLevel: 2,
	}
	m.baseEntity.RLock()
	nc.Coord = structs.ShineCoordType{
		XY: structs.ShineXYType{
			X: uint32(m.current.x),
			Y: uint32(m.current.y),
		},
		Direction: uint8(m.current.d),
	}
	m.baseEntity.RUnlock()
	return nc
}
