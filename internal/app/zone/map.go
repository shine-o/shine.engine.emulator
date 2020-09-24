package zone

import (
	"bytes"
	"github.com/RoaringBitmap/roaring"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/blocks"
	mobs "github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/monsters"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/shn"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/world"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/networking"
	"github.com/spf13/viper"
	"math"
	"sync"
)

type zoneMap struct {
	data      *world.Map
	walkableX *roaring.Bitmap
	walkableY *roaring.Bitmap
	entities  *entities
	events
}

type entities struct {
	*players
	*monsters
	*npcs
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
		zm.spawnNPC()
	}()
	wg.Wait()

	go zm.removeInactiveHandles()

	for i := 0; i <= num; i++ {
		go zm.mapHandles()

		go zm.playerActivity()
		go zm.monsterActivity()

		go zm.playerQueries()
		go zm.monsterQueries()
	}
}

func (zm * zoneMap) spawnNPC()  {
	npcs, ok := npcData.Data[zm.data.MapInfoIndex]

	if ok {
		var (
			sem = make(chan int, 100)
			wg  sync.WaitGroup
		)

		for _, data := range npcs {
			wg.Add(1)
			sem <- 1
			go func(sn * world.ShineNPC) {
				defer wg.Done()

				mi, mis := mobDataPointers(sn.MobIndex)

				if mi == nil {
					log.Errorf("no entry in MobInfo for %v", sn.MobIndex)
					return
				}

				if mis == nil {
					log.Errorf("no entry in MobInfoServer for %v", sn.MobIndex)
					return
				}

				h, err := zm.entities.npcs.new(npcHandleMin, npcHandleMax, npcAttemptsMax)
				if err != nil {
					log.Error(err)
					return
				}

				n := &npc{
					baseEntity: baseEntity{
						handle: h,
						fallback: location{
							x: sn.X,
							y: sn.Y,
							d: sn.D,
						},
						current: location{
							mapID:     zm.data.ID,
							mapName:   zm.data.MapInfoIndex,
							x: sn.X,
							y: sn.Y,
							d: sn.D,
							movements: [15]movement{},
						},
						events: events{},
					},
					hp:            mi.MaxHP,
					sp:            uint32(mis.MaxSP),
					mobInfo:       mi,
					mobInfoServer: mis,
					status: status{
						idling:   make(chan bool),
						fighting: make(chan bool),
						chasing:  make(chan bool),
						fleeing:  make(chan bool),
					},
				}

				zm.entities.npcs.Lock()
				zm.entities.npcs.active[h] = n
				zm.entities.npcs.Unlock()

				<-sem
			}(data)
		}
	}
}

// a mob is a collection of entities, in this case monsters
func (zm *zoneMap) spawnMobs() {
	spawnData, ok := monsterData.MapRegens[zm.data.MapInfoIndex]
	if ok {
		var (
			sem = make(chan int, 100)
			wg  sync.WaitGroup
		)

		for _, group := range spawnData.Groups {
			wg.Add(1)
			sem <- 1
			go func(g mobs.RegenEntry) {
				defer wg.Done()
				spawnMob(zm, g)
				<-sem
			}(group)
		}
		wg.Wait()
	}
}

func mobDataPointers(mobIndex string) (*shn.MobInfo, *shn.MobInfoServer) {
	var (
		mi  *shn.MobInfo
		mis *shn.MobInfoServer
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

func spawnMob(zm *zoneMap, re mobs.RegenEntry) {

	var iwg sync.WaitGroup

	var sem = make(chan int, 100)

	for _, mob := range re.Mobs {
		var (
			mi  *shn.MobInfo
			mis *shn.MobInfoServer
		)

		for i, row := range monsterData.MobInfo.ShineRow {
			if row.InxName == mob.Index {
				mi = &monsterData.MobInfo.ShineRow[i]
			}
		}

		for i, row := range monsterData.MobInfoServer.ShineRow {
			if row.InxName == mob.Index {
				mis = &monsterData.MobInfoServer.ShineRow[i]
			}
		}

		if mi == nil {
			log.Errorf("no entry in MobInfo for %v", mob.Index)
			break
		}

		if mis == nil {
			log.Errorf("no entry in MobInfoServer for %v", mob.Index)
			break
		}

		for i := 0; i < int(mob.Num); i++ {
			iwg.Add(1)
			sem <- 1
			go func() {
				iwg.Done()
				spawnMonster(zm, re, mi, mis)
				<-sem
			}()
		}
	}

	iwg.Wait()
}

func spawnMonster(zm *zoneMap, re mobs.RegenEntry, mi *shn.MobInfo, mis *shn.MobInfoServer) {
	var (
		x, y, d  int
		maxTries = 1000 // todo: use goroutines for this, with a maximum of 50 routines, first one to get a spot wins the spawn
		spawn    = false
		staticMonster = false
		numSteps = 5
		speed = 30
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
			re.Width = networking.RandomIntBetween(20, 250)
		}

		if re.Height == 0 {
			re.Height = networking.RandomIntBetween(20, 250)
		}

		x = networking.RandomIntBetween(re.X, re.X+re.Width)

		y = networking.RandomIntBetween(re.Y, re.Y+re.Height)

		d = networking.RandomIntBetween(1, 250)

		spawn = validateLocation(zm, x, y, numSteps, speed)

		maxTries--

	}

	if spawn {
		h, err := zm.entities.monsters.new(monsterHandleMin, monsterHandleMax, monsterAttemptsMax)
		if err != nil {
			log.Error(err)
			return
		}

		m := &monster{
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
					movements: [15]movement{},
				},
				events: events{},
			},
			hp:            mi.MaxHP,
			sp:            uint32(mis.MaxSP),
			mobInfo:       mi,
			mobInfoServer: mis,
			regenData:     &re,
			status: status{
				idling:   make(chan bool),
				fighting: make(chan bool),
				chasing:  make(chan bool),
				fleeing:  make(chan bool),
			},
		}

		zm.entities.monsters.Lock()
		zm.entities.monsters.active[h] = m
		zm.entities.monsters.Unlock()

		if !staticMonster {
			go m.roam(zm)
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

			rx, ry := igCoordToBitmap(lx, ly)
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
func igCoordToBitmap(x, y int) (int, int) {
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
func walkingPositions(s *blocks.SHBD) (*roaring.Bitmap, *roaring.Bitmap, error) {
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
