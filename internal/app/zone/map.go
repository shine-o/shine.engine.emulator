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
	entities * entities
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
			go spawnMobGroup(zm, group, &wg)
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
func spawnMobGroup(zm *zoneMap, re mobs.RegenEntry, wg *sync.WaitGroup) {
	defer wg.Done()



	var (
		mi  *shn.MobInfo
		mis *shn.MobInfoServer
	)

	for i, row := range monsterData.MobInfo.ShineRow {
		if row.InxName == re.MobIndex {
			mi = &monsterData.MobInfo.ShineRow[i]
		}
	}

	for i, row := range monsterData.MobInfoServer.ShineRow {
		if row.InxName == re.MobIndex {
			mis = &monsterData.MobInfoServer.ShineRow[i]
		}
	}

	if mi == nil {
		log.Errorf("no entry in MobInfo for %v", re.MobIndex)
		return
	}

	if mis == nil {
		log.Errorf("no entry in MobInfoServer for %v", re.MobIndex)
		return
	}

	for i := re.MobNum; i != 0; i-- {
		var (
			x, y     int
			maxTries = 20
			spawn    = false
		)

		for maxTries != 0 {

			if re.Width == 0 {
				x = re.X
			} else {
				x = networking.RandomIntBetween(re.X, re.X+re.Width)
			}

			if re.Height == 0 {
				y = re.Y
			} else {
				y = networking.RandomIntBetween(re.Y, re.Y+re.Height)
			}

			rX, rY := igCoordToBitmap(x, y)

			if canWalk(zm.walkableX, zm.walkableY, uint32(rX), uint32(rY)) {
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

			monster := &monster{
				baseEntity: baseEntity{
					handle: h,
					location: location{
						mapID:     zm.data.ID,
						mapName:   zm.data.MapInfoIndex,
						x:         uint32(x),
						y:         uint32(y),
						d:         0,
						movements: [15]movement{},
					},
					events: events{},
				},
				hp:            mi.MaxHP,
				sp:            uint32(mis.MaxSP),
				mobInfo:       mi,
				mobInfoServer: mis,
				regenData:     re,
			}

			zm.entities.monsters.Lock()
			zm.entities.monsters.active[h] = monster
			zm.entities.monsters.Unlock()
		}
	}

}

func igCoordToBitmap(x, y int) (int, int) {
	rX := (x * 8) / 50
	rY := (y * 8) / 50
	return rX, rY
}

// CanWalk translates in game coordinates to SHBD coordinates
func canWalk(x, y *roaring.Bitmap, rX, rY uint32) bool {
	if x.ContainsInt(int(rX)) && y.ContainsInt(int(rY)) {
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
