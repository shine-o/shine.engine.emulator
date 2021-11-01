package zone

import (
	"reflect"
	"sync"

	"github.com/shine-o/shine.engine.emulator/internal/pkg/crypto"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/spf13/viper"
)

type zoneMap struct {
	data *data.Map
	//rawNodes              path.Grid
	validCoordinates      ValidCoordinates
	presetNodes           Grid
	presetNodesWithMargin Grid
	entities              *entities
	events                *events
}

type entities struct {
	players map[uint16]entity
	npc     map[uint16]entity
	monster map[uint16]entity
	sync.RWMutex
}

func (e *entities) getPlayer(handle uint16) (*player, error) {
	e.RLock()
	defer e.RUnlock()
	p, ok := e.players[handle].(*player)
	if !ok {
		return nil, errors.Err{
			Code: errors.ZoneMissingPlayer,
			Details: errors.Details{
				"handle": handle,
			},
		}
	}
	return p, nil
}

func (e *entities) getNpc(handle uint16) *npc {
	e.RLock()
	defer e.RUnlock()
	p, ok := e.npc[handle].(*npc)
	if ok {
		return p
	}
	return nil
}

func (e *entities) all() <-chan entity {
	e.RLock()
	ch := make(chan entity, len(e.players)+len(e.npc)+len(e.monster))
	e.RUnlock()

	go func(el *entities, send chan<- entity) {
		el.RLock()
		for _, e := range el.players {
			send <- e
		}
		for _, e := range el.npc {
			send <- e
		}
		for _, e := range el.monster {
			send <- e
		}
		el.RUnlock()
		close(send)
	}(e, ch)

	return ch
}

func (e *entities) allPlayers() <-chan *player {
	e.RLock()
	ch := make(chan *player, len(e.players))
	e.RUnlock()

	go func(el *entities, send chan<- *player) {
		el.RLock()
		for _, e := range el.players {
			p, ok := e.(*player)
			if ok {
				send <- p
			} else {
				log.Error(errors.Err{
					Code:    errors.ZoneBadEntityType,
					Message: "",
					Details: errors.Details{
						"expected": "player",
						"got":      reflect.TypeOf(e).String(),
					},
				})
			}
		}
		el.RUnlock()
		close(send)
	}(e, ch)

	return ch
}

func (e *entities) allNpc() <-chan *npc {
	e.RLock()
	ch := make(chan *npc, len(e.npc))
	e.RUnlock()

	go func(el *entities, send chan<- *npc) {
		el.RLock()
		for _, e := range el.npc {
			n, ok := e.(*npc)
			if ok {
				send <- n
			} else {
				log.Error(errors.Err{
					Code:    errors.ZoneBadEntityType,
					Message: "",
					Details: errors.Details{
						"expected": "npc",
						"got":      reflect.TypeOf(e).String(),
					},
				})
			}
		}
		el.RUnlock()
		close(send)
	}(e, ch)

	return ch
}

func (e *entities) allMonsters() <-chan *monster {
	e.RLock()
	ch := make(chan *monster, len(e.players))
	e.RUnlock()

	go func(el *entities, send chan<- *monster) {
		el.RLock()
		for _, e := range el.monster {
			m, ok := e.(*monster)
			if ok {
				send <- m
			} else {
				log.Error(errors.Err{
					Code:    errors.ZoneBadEntityType,
					Message: "",
					Details: errors.Details{
						"expected": "monster",
						"got":      reflect.TypeOf(e).String(),
					},
				})
			}
		}
		el.RUnlock()
		close(send)
	}(e, ch)

	return ch
}

func (e *entities) removePlayer(h uint16) {
	e.Lock()
	delete(e.players, h)
	e.Unlock()
}

func (e *entities) getEntity(handle uint16) (entity, error) {
	for en := range e.all() {
		if en.getHandle() == handle {
			return en, nil
		}
	}
	return nil, errors.Err{
		Code:    errors.ZoneEntityNotFound,
		Message: "",
		Details: errors.Details{
			"handle": handle,
		},
	}
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
		Details: errors.Details{
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

	zm.addEntity(npc)
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

func (zm *zoneMap) spawnMonster(baseNpc *npc, re *data.RegenEntry) {
	var x, y, d int

	staticMonster := false
	mi := baseNpc.data.mobInfo
	mis := baseNpc.data.mobInfoServer

	if mi.WalkSpeed == 0 && mi.RunSpeed == 0 {
		staticMonster = true
	}

	bx, by := bitmapCoordinates(re.X, re.Y)

	node := zm.presetNodes.GetNearest(bx, by)

	if node != nil {
		x, y = gameCoordinates(node.X, node.Y)
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

	m := &monster{
		baseEntity: &baseEntity{
			handle: h,
			fallback: location{
				x: x,
				y: y,
				d: d,
			},
			proximity: &entityProximity{
				entities: make(map[uint16]entity),
			},
			current: location{
				mapID:   zm.data.ID,
				mapName: zm.data.MapInfoIndex,
				x:       x,
				y:       y,
				d:       d,
			},
		},
		targeting: &targeting{
			players:  make(map[uint16]*player),
			monsters: make(map[uint16]*monster),
			npc:      make(map[uint16]*npc),
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

	zm.addEntity(m)

	if !staticMonster {
		// for now its a weird thing, better not to use it
		// go m.roam(zm)
	}
}

func (zm *zoneMap) addEntity(e entity) {
	zm.entities.Lock()
	h := e.getHandle()
	switch e.(type) {
	case *player:
		zm.entities.players[h] = e
		break
	case *npc:
		zm.entities.npc[h] = e
		break
	case *monster:
		zm.entities.monster[h] = e
		break
	default:
		log.Infof("unkown entity type %v", reflect.TypeOf(e).String())
	}
	zm.entities.Unlock()
}

func (zm *zoneMap) removeEntity(e entity) {
	zm.entities.Lock()
	switch e.(type) {
	case *player:
		delete(zm.entities.players, e.getHandle())
		break
	case *npc:
		delete(zm.entities.npc, e.getHandle())
		break
	}
	zm.entities.Unlock()
}

func validateLocation(zm *zoneMap, x, y, numSteps, speed int) bool {
	directions := make(map[string]bool)

	directions["left"] = false
	directions["right"] = false
	directions["up"] = false
	directions["down"] = false

	for k := range directions {
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
			if CanWalk(zm.presetNodes, rx, ry) {
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
			Details: errors.Details{
				"mapID": mapID,
			},
		}
	}

	var (
		//rawNodes               = path.RawGrid(md.SHBD)
		validCoordinates       = GetValidCoordinates(md.SHBD)
		presetNodes            = PresetNodesGrid(md.SHBD, PresetNodesMargin)
		presetNodesWallMargins = PresetNodesWithMargins(md.SHBD, presetNodes, PresetNodesMargin)
	)

	zm := &zoneMap{
		data: md,
		//rawNodes:              rawNodes,
		validCoordinates:      validCoordinates,
		presetNodes:           presetNodes,
		presetNodesWithMargin: presetNodesWallMargins,
		entities: &entities{
			players: make(map[uint16]entity),
			npc:     make(map[uint16]entity),
			monster: make(map[uint16]entity),
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

// add entities within range and return a slice of the added entities
func addWithinRangeEntities(e1 entity, zm *zoneMap) []entity {
	var newEntities []entity

	for e2 := range zm.entities.all() {
		if e1.getHandle() == e2.getHandle() {
			continue
		}
		if withinRange(e1, e2) {
			if e1.alreadyNearbyEntity(e2) {
				continue
			}
			e1.addNearbyEntity(e2)
			newEntities = append(newEntities, e2)
		}
	}
	return newEntities
}

// remove entities out of range and return a slice of the removed entities
func removeOutOfRangeEntities(e1 entity) []entity {
	var removedEntities []entity
	for e2 := range e1.getNearbyEntities() {
		if e1.getHandle() == e2.getHandle() {
			continue
		}
		if !withinRange(e1, e2) {
			e1.removeNearbyEntity(e2)
			removedEntities = append(removedEntities, e2)
		}
		p, ok := e2.(*player)
		if ok {
			err := validatePlayerEntity(p)
			if err != nil {
				continue
			}
			if lastHeartbeat(p) > playerHeartbeatLimit {
				e1.removeNearbyEntity(e2)
				removedEntities = append(removedEntities, e2)
			}
		}
	}
	return removedEntities
}
