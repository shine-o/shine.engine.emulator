package zone

import (
	"bytes"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/blocks"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/utils"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/world"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
	"math"
	"os"
)

type zoneMap struct {
	data      world.MapData
	mobRegens interface{}
	walkableX *roaring.Bitmap
	walkableY *roaring.Bitmap
	entities  entities
	events
}

type entities struct {
	*players
	*monsters
}

func (zm *zoneMap) run() {
	// load NPCs for this map
	// run logic routines
	// as many workers as needed can be launched
	num := viper.GetInt("workers.num_zone_workers")

	go zm.removeInactiveHandles()

	for i := 0; i <= num; i++ {
		go zm.mapHandles()
		go zm.playerActivity()
		go zm.playerQueries()
		go zm.monsterQueries()
	}
}

// load maps
func loadMaps() []zoneMap {
	var reload bool

	dbPath := viper.GetString("game_data.database")

	_, err := os.Stat(dbPath)
	if os.IsNotExist(err) {
		reload = true
	} else {
		reload = viper.GetBool("game_data.reload")
	}

	shinePath := viper.GetString("shine_folder")
	normalMaps := viper.GetIntSlice("normal_maps")

	db, err := bolt.Open(dbPath, 0600, nil)

	if err != nil {
		log.Fatal(err)
	}

	if reload {
		shineAbsPath, err := utils.ValidPath(shinePath)
		if err != nil {
			log.Fatal(err)
		}
		err = world.LoadMapData(shineAbsPath, db)
		if err != nil {
			log.Fatal(err)
		}
	}

	// zones defined in the service, load them and register the loaded ones to the master
	var zoneMaps []zoneMap
	for _, id := range normalMaps {
		var md world.MapData

		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("maps"))
			data := b.Get([]byte(fmt.Sprintf("%v", id)))
			err = structs.Unpack(data, &md)
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			log.Fatalf("failed to get stored map data with id %v %v", id, err)
		}

		walkableX, walkableY, err := walkingPositions(&md.SHBD)
		if err != nil {
			log.Fatal(err)
		}
		zm := zoneMap{
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
		zoneMaps = append(zoneMaps, zm)
	}

	return zoneMaps
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
