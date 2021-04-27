package zone

import (
	"bytes"
	"github.com/RoaringBitmap/roaring"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/crypto"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/errors"
	"github.com/spf13/viper"
	"math"
	"sync"
)

type zoneMap struct {
	data      *data.Map
	walkableX *roaring.Bitmap
	walkableY *roaring.Bitmap
	entities  *entities
	events    *events
}

type entities struct {
	players *players
	npcs    *npcs
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
			go func(inxName string) {
				zm.spawnNPC(inxName, npcs)
				wg.Done()
				<-sem
			}(npc.MobIndex)
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

func (zm *zoneMap) spawnNPC(inxName string, shineNPCS []*data.ShineNPC) {
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

	var sn * data.ShineNPC
	for _, npc := range shineNPCS {
		if npc.MobIndex == inxName {
			sn = npc
			break
		}
	}

	if sn == nil {
		log.Error(errors.Err{
			Code:    errors.ZoneMissingNpcData,
			Message: "missing shineNPC data",
			Details: errors.ErrDetails{
				"mobInx": inxName,
				"mapInx": zm.data.MapInfoIndex,
			},
		})
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
		wg sync.WaitGroup
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
		maxTries      = 1000 // todo: use goroutines for this, with a maximum of 50 routines, first one to get a spot wins the spawn
		spawn         = false
		staticMonster = false
		numSteps      = 5
		speed         = 30
		mi            = monsterNpc.data.mobInfo
		mis           = monsterNpc.data.mobInfoServer
	)

	if mi.WalkSpeed == 0 && mi.RunSpeed == 0 {
		staticMonster = true
		numSteps = 20
		speed = 20
	}

	for maxTries != 0 {
		// todo use go routines for this too, with a semaphore
		if spawn {
			break
		}

		if re.Width == 0 {
			re.Width = crypto.RandomIntBetween(20, 250)
		}

		if re.Height == 0 {
			re.Height = crypto.RandomIntBetween(20, 250)
		}

		x = crypto.RandomIntBetween(re.X, re.X+re.Width)

		y = crypto.RandomIntBetween(re.Y, re.Y+re.Height)

		d = crypto.RandomIntBetween(1, 250)

		spawn = validateLocation(zm, x, y, numSteps, speed)

		maxTries--

	}

	if spawn {
		h, err := newHandle()
		if err != nil {
			log.Error(err)
			return
		}

		m := &npc{
			baseEntity: baseEntity{
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
			if canWalk(zm.walkableX, zm.walkableY, rx, ry) {
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
func bitmapCoordToIg(rX, rY int) (int, int) {
	igX := (rX * 50) / 8
	igY := (rY * 50) / 8
	return igX, igY
}

func canWalk(x, y *roaring.Bitmap, rX, rY int) bool {
	if x.ContainsInt(rX) && y.ContainsInt(rY) {
		return true
	}
	return false
}

// creates two X,Y roaring bitmaps with walkable coordinates
func walkingPositions(s *data.SHBD) (*roaring.Bitmap, *roaring.Bitmap, error) {
	walkableX := roaring.BitmapOf()
	walkableY := roaring.BitmapOf()

	r := bytes.NewReader(s.Data)

	for y := 0; y < s.Y; y++ {
		for x := 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return walkableX, walkableY, err
			}
			for i := 0; i < 8; i++ {
				if b&byte(math.Pow(2, float64(i))) == 0 {
					rX := x*8 + i
					rY := y
					walkableX.Add(uint32(rX))
					walkableY.Add(uint32(rY))
				}
			}
		}
	}
	return walkableX, walkableY, nil
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

	walkableX, walkableY, err := walkingPositions(md.SHBD)

	if err != nil {
		return nil, err
	}

	zm := &zoneMap{
		data:      md,
		walkableX: walkableX,
		walkableY: walkableY,
		entities: &entities{
			players: &players{
				active: make(map[uint16]*player),
			},
			npcs: &npcs{
				active: make(map[uint16]*npc),
			},
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
