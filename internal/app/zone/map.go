package zone

import (
	"bytes"
	"github.com/RoaringBitmap/roaring"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/blocks"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/world"
	"github.com/spf13/viper"
	"math"
)

type zoneMap struct {
	data      * world.Map
	walkableX *roaring.Bitmap
	walkableY *roaring.Bitmap
	entities  entities
	events
}

type entities struct {
	*players
	*monsters
}

func (z * zone) addMap(mapId int) {

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
		entities: entities{
			players: &players{
				handleIndex: playerHandleMin,
				active:      make(map[uint16]*player),
			},
			monsters: &monsters{
				handleIndex: playerHandleMin,
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
		unknownHandle, 	monsterAppeared, monsterDisappeared, monsterWalks,monsterRuns,
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

	monsters, ok := monsterData.MapRegens[zm.data.MapInfoIndex]

	if !ok {
		log.Warningf("No mobs available for map %v", zm.data.MapInfoIndex)
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