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

func (z *zone) addMap(mapId int, wg *sync.WaitGroup) {
	defer wg.Done()
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
	// load NPCs for this map
	// run logic routines
	// as many workers as needed can be launched
	num := viper.GetInt("workers.num_map_workers")

	// load mobs

	//zm.entities.monsters

	spawnData, ok := monsterData.MapRegens[zm.data.MapInfoIndex]

	if ok {
		// for each groupIndex in spawnData
		// 		for each number of mobs in groupIndex
		//		populate a monster object using monster data
		//		launch monster go routines
		//		keep pointers to monster data
		var wg sync.WaitGroup

		for _, group := range spawnData.Groups {
			wg.Add(1)
			go spawnMob(zm, group, &wg)
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
func spawnMob(zm *zoneMap, re mobs.RegenEntry, wg *sync.WaitGroup) {
	defer wg.Done()

	for _, mob := range re.Mobs {

		for i := 0; i <= int(mob.Num); i++ {

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
				return
			}

			if mis == nil {
				log.Errorf("no entry in MobInfoServer for %v", mob.Index)
				return
			}

			var (
				x, y, d  int
				maxTries = 500 // todo: use goroutines for this, with a maximum of 50 routines, first one to get a spot wins the spawn
				spawn    = false
			)

			for maxTries != 0 {

				if spawn {
					break
				}

				if re.Width == 0 {
					re.Width = networking.RandomIntBetween(100, 150)
				}

				if re.Height == 0 {
					re.Height = networking.RandomIntBetween(100, 150)
				}

				x = networking.RandomIntBetween(re.X, re.X+re.Width)
				y = networking.RandomIntBetween(re.Y, re.Y+re.Height)

				d = networking.RandomIntBetween(1, 250)

				rX, rY := igCoordToBitmap(x, y)

				if canWalk(zm.walkableX, zm.walkableY, rX, rY) {
					spawn = true
				}

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

				if m.mobInfo.RunSpeed > 0 && m.mobInfo.WalkSpeed > 0 {
					go m.roam(zm)
				}
			}
		}
	}

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
