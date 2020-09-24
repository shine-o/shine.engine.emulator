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
}

func (z *zone) addMap(mapId int) {
	md, ok := mapData[mapId]

	if !ok {
		log.Fatalf("no map data for map with id %v", mapId)
	}

	walkableX, walkableY, err := walkingPositions(md.SHBD)

	if err != nil {
		log.Fatal(err)
	}

	m := &zoneMap{
		data:      md,
		walkableX: walkableX,
		walkableY: walkableY,
		entities: &entities{
			players: &players{
				handleIndex: playerHandleMin,
				active:      make(map[uint16]*player),
			},
			monsters: &monsters{
				handleIndex: monsterHandleMin,
				active:      make(map[uint16]*monster),
			},
		},
		events: events{
			send: make(sendEvents),
			recv: make(recvEvents),
		},
	}

	events := []eventIndex{
		playerHandle,
		playerHandleMaintenance,
		queryPlayer, queryMonster,
		playerAppeared, playerDisappeared, playerJumped, playerWalks, playerRuns, playerStopped,
		unknownHandle, monsterAppeared, monsterDisappeared, monsterWalks, monsterRuns,
	}

	for _, index := range events {
		c := make(chan event, 500)
		m.recv[index] = c
		m.send[index] = c
	}

	z.Lock()
	z.rm[m.data.ID] = m
	z.Unlock()

	go m.run()

}

func (zm *zoneMap) run() {

	num := viper.GetInt("workers.num_map_workers")

	spawnData, ok := monsterData.MapRegens[zm.data.MapInfoIndex]

	if ok {
		var (
			sem = make(chan int, 100)
			wg sync.WaitGroup
		)

		for _, group := range spawnData.Groups {
			wg.Add(1)
			sem <- 1
			go func(g mobs.RegenEntry) {
				defer wg.Done()
				spawnMob(zm, g)
				<- sem
			}(group)
		}
		wg.Wait()
	}

	go zm.removeInactiveHandles()

	for i := 0; i <= num; i++ {
		go zm.mapHandles()

		go zm.playerActivity()
		go zm.monsterActivity()

		go zm.playerQueries()
		go zm.monsterQueries()
	}

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
				spawn(zm, re, mi, mis)
				<-sem
			}()
		}
	}

	iwg.Wait()
}

func spawn(zm *zoneMap, re mobs.RegenEntry, mi *shn.MobInfo, mis *shn.MobInfoServer) {
	var (
		x, y, d  int
		maxTries = 100 // todo: use goroutines for this, with a maximum of 50 routines, first one to get a spot wins the spawn
		spawn    = false
		staticMonster = false
	)

	if mi.WalkSpeed == 0 && mi.RunSpeed == 0 {
		staticMonster = true
	}

	for maxTries != 0 {

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

		spawn = validateLocation(zm, x, y, 5, 30)

		maxTries--

	}

	if spawn {
		h, err := zm.entities.monsters.newHandle()
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
func igCoordToBitmap(x, y int) (int, int) {
	rX := (x * 8) / 50
	rY := (y * 8) / 50
	return rX, rY
}

func bitmapCoordToIg(rX, rY int) (int, int) {
	igX := (rX * 50) / 8
	igY := (rY * 50) / 8
	return igX, igY
}

func calculateSpawnCoordinates(igx, igy, width, height, rotate int) (int, int) {
	//var cos, sin float64
	v := float64(rotate) * 0.01745329
	cos := int(math.Cos(v) * 1024.0)
	sin := int(math.Sin(v) * 1024.0)

	//  v2 = 0;
	//  v3 = this;
	//  if ( pLoc )
	//  {
	//    if ( this->mrlr_Width > 0 )
	//    {
	//      v4 = this->mrlr_Width;
	//      v2 = rand() % (2 * v4) - v4;
	//      v5 = v3->mrlr_Height;
	//      ody = rand() % (2 * v5) - v5;
	//    }
	//    else
	//    {
	//      ody = 0;
	//    }
	//    v6 = ody * v3->mrlr_CosD1024 - v2 * v3->mrlr_SinD1024;
	//    pLoc->x = v3->mrlr_X + (v2 * v3->mrlr_CosD1024 + ody * v3->mrlr_SinD1024) / 1024;
	//    pLoc->y = v6 / 1024 + v3->mrlr_Y;
	//  }
	var ody, v2 int
	v2 = 0
	if width > 0 {
		v2 =  networking.RandomIntBetween(0, width)
		ody = networking.RandomIntBetween(0, height)
	} else {
		ody = 0
	}

	v6 := (ody * cos) - (height  * sin)
	x := igx + ( (v2 * cos) + (ody * sin) ) / 1024
	y := (v6 / 1024) + igy

	return x, y
}

// CanWalk translates in game coordinates to SHBD coordinates
func canWalk(x, y *roaring.Bitmap, rX, rY int) bool {
	if x.ContainsInt(rX) && y.ContainsInt(rY) {
		return true
	}
	return false
}

// WalkingPositions creates two X,Y roaring bitmaps with walkable coordinates
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
