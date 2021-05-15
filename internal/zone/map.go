package zone

import (
	"github.com/shine-o/shine.engine.emulator/internal/pkg/crypto"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	path "github.com/shine-o/shine.engine.emulator/internal/pkg/pathfinding"
	"github.com/spf13/viper"
	"sync"
)

type zoneMap struct {
	data                  *data.Map
	rawNodes              path.Grid
	presetNodes           path.Grid
	presetNodesWithMargin path.Grid
	entities              *entities
	entities2             *entities2
	events                *events
}

type entities struct {
	players *players
	npcs    *npcs
}

type entities2 struct {
	list entitiesmap
	sync.RWMutex
}


func (zm *zoneMap) run() {

	num := viper.GetInt("workers.num_map_workers")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		zm.spawnMobs()
	}()

	go func() {
		defer wg.Done()
		zm.spawnNPCs()
	}()

	wg.Wait()

	go zm.removeInactiveHandles()

	for i := 0; i <= num; i++ {
		go zm.mapHandles()

		go zm.playerActivity()
		go zm.monsterActivity()
		go zm.npcInteractions()
	}
}

func (zm *zoneMap) spawnNPCs() {
	npcs, ok := npcData.MapNPCs[zm.data.MapInfoIndex]

	if ok {
		var (
			sem = make(chan int, 100)
			wg  sync.WaitGroup
		)
		for _, npc := range npcs {
			wg.Add(1)
			sem <- 1
			go func(inxName string, sn *data.ShineNPC) {
				zm.spawnNPC(inxName, sn)
				wg.Done()
				<-sem
			}(npc.MobIndex, npc)
		}

		wg.Wait()
		return
	}
	log.Error(errors.Err{
		Code:    errors.ZoneMissingNpcData,
		Message: "",
		Details: errors.ErrDetails{
			"mapInx": zm.data.MapInfoIndex,
		},
	})
}

func (zm *zoneMap) spawnNPC(inxName string, sn *data.ShineNPC) {
	h, err := newHandle()

	if err != nil {
		log.Error(err)
		return
	}

	npc, err := loadBaseNpc(inxName, isNPC)

	if err != nil {
		log.Error(err)
		return
	}

	npc.nType = getNpcType(sn.Role, sn.RoleArg)
	npc.data.npcData = sn
	npc.baseEntity.handle = h
	npc.spawnLocation(zm)

	zm.entities.npcs.Lock()
	zm.entities.npcs.active[h] = npc
	zm.entities.npcs.Unlock()
}

// a mob is a collection of entities, in this case monsters
func (zm *zoneMap) spawnMobs() {
	spawnData, ok := monsterData.MapRegens[zm.data.MapInfoIndex]
	if ok {
		var (
			sem = make(chan int, 300)
			wg  sync.WaitGroup
		)

		for _, group := range spawnData.Groups {
			wg.Add(1)
			sem <- 1
			go func(g *data.RegenEntry) {
				defer wg.Done()
				zm.spawnMob(g)
				<-sem
			}(group)
		}
		wg.Wait()
	}
}

func (zm *zoneMap) spawnMob(re *data.RegenEntry) {

	var (
		wg  sync.WaitGroup
		sem = make(chan int, 100)
	)

	for _, mob := range re.Mobs {
		monsterNpc, err := loadBaseNpc(mob.Index, isMonster)
		if err != nil {
			log.Error(err)
			continue
		}

		for i := 0; i < int(mob.Num); i++ {
			wg.Add(1)
			sem <- 1
			go func() {
				defer wg.Done()
				zm.spawnMonster(monsterNpc, re)
				<-sem
			}()
		}
	}

	wg.Wait()
}

func (zm *zoneMap) spawnMonster(monsterNpc *npc, re *data.RegenEntry) {
	var (
		x, y, d       int
		staticMonster = false
		mi            = monsterNpc.data.mobInfo
		mis           = monsterNpc.data.mobInfoServer
	)

	if mi.WalkSpeed == 0 && mi.RunSpeed == 0 {
		staticMonster = true
	}

	bx, by := bitmapCoordinates(re.X, re.Y)

	node := zm.presetNodes.GetNearest(bx, by)

	if node != nil {
		x = node.X
		y = node.Y
		d = crypto.RandomIntBetween(1, 250)
	} else {
		log.Error("could not find spawn position, exiting")
		return
	}

	h, err := newHandle()

	if err != nil {
		log.Error(err)
		return
	}

	m := &npc{
		baseEntity: &baseEntity{
			handle: h,
			fallback: location{
				x: x,
				y: y,
				d: d,
			},
			current: location{
				mapID:     zm.data.ID,
				mapName:   zm.data.MapInfoIndex,
				x:         x,
				y:         y,
				d:         d,
				movements: []movement{},
			},
		},
		stats: &npcStats{
			hp: mi.MaxHP,
			sp: uint32(mis.MaxSP),
		},
		data: &npcStaticData{
			mobInfo:       mi,
			mobInfoServer: mis,
			regenData:     re,
		},
		state: &entityState{
			idling:   make(chan bool),
			fighting: make(chan bool),
			chasing:  make(chan bool),
			fleeing:  make(chan bool),
		},
		ticks: &entityTicks{},
	}

	zm.entities.npcs.Lock()
	zm.entities.npcs.active[h] = m
	zm.entities.npcs.Unlock()

	if !staticMonster {
		// for now its a weird thing, better not to use it
		//go m.roam(zm)
	}
}

func (zm *zoneMap) addEntity(e entity) {
	zm.entities2.Lock()
	zm.entities2.list[e.getHandle()] = e
	zm.entities2.Unlock()
}

func validateLocation(zm *zoneMap, x, y, numSteps, speed int) bool {
	var directions = make(map[string]bool)

	directions["left"] = false
	directions["right"] = false
	directions["up"] = false
	directions["down"] = false

	for k, _ := range directions {
		lx := x
		ly := y

		for i := 0; i < numSteps; i++ {
			switch k {
			case "left":
				lx -= speed
				break
			case "right":
				lx += speed
				break
			case "up":
				ly -= speed
				break
			case "down":
				ly += speed
				break
			}

			rx, ry := bitmapCoordinates(lx, ly)
			if path.CanWalk(zm.presetNodesWithMargin, rx, ry) {
				continue
			}

			return false
		}
		directions[k] = true
	}
	return true
}

// translates in game coordinates to SHBD coordinates
func bitmapCoordinates(x, y int) (int, int) {
	rX := (x * 8) / 50
	rY := (y * 8) / 50
	return rX, rY
}

// translates SHBD coordinates to in game coordinates
func gameCoordinates(bX, bY int) (int, int) {
	gx := (bX * 50) / 8
	gy := (bY * 50) / 8
	return gx, gy
}

func loadMap(mapID int) (*zoneMap, error) {
	md, ok := mapData.Maps[mapID]

	if !ok {
		return nil, errors.Err{
			Code:    errors.ZoneMissingMapData,
			Message: "",
			Details: errors.ErrDetails{
				"mapID": mapID,
			},
		}
	}

	var (
		rawNodes               = path.RawGrid(md.SHBD)
		presetNodes            = path.PresetNodesGrid(md.SHBD, path.PresetNodesMargin)
		presetNodesWallMargins = path.PresetNodesWithMargins(md.SHBD, rawNodes, path.PresetNodesMargin)
	)

	zm := &zoneMap{
		data:                  md,
		rawNodes:              rawNodes,
		presetNodes:           presetNodes,
		presetNodesWithMargin: presetNodesWallMargins,
		entities: &entities{
			players: &players{
				active: make(map[uint16]*player),
			},
			npcs: &npcs{
				active: make(map[uint16]*npc),
			},
		},
		entities2: &entities2{
			list:    make(entitiesmap),
		},
		events: &events{
			send: make(sendEvents),
			recv: make(recvEvents),
		},
	}

	for _, index := range mapEvents {
		c := make(chan event, 500)
		zm.events.recv[index] = c
		zm.events.send[index] = c
	}

	return zm, nil
}

func addWithinRangeEntities(e1 entity, zm *zoneMap) {
	for e2 := range zm.entities2.all() {
		if withinRange(e1, e2) {
			if e1.getHandle() == e2.getHandle() {
				continue
			}
			e1.addNearbyEntity(e2)
			// go e1.
		}
	}
}

func removeOutOfRangeEntities(e1 entity) {
	for e2 := range e1.getNearbyEntities() {
		if !withinRange(e1, e2) {
			if e1.getHandle() == e2.getHandle() {
				continue
			}
			e1.removeNearbyEntity(e2)
			// go e1.
		}
	}
}