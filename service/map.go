package service

import (
	"bytes"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/shine-o/shine.engine.core/game-data/blocks"
	"github.com/shine-o/shine.engine.core/game-data/utils"
	"github.com/shine-o/shine.engine.core/game-data/world"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
	"math"
	"sync"
)


type zoneMap struct {
	data      world.MapData
	walkableX *roaring.Bitmap
	walkableY *roaring.Bitmap
	handles   mapEntities
	send      sendEvents
	recv      recvEvents
}

type mapEntities struct {
	counter  uint16
	players  map[uint16]*player
	monsters map[uint16]*monster
	mu       sync.RWMutex
}

func (zm *zoneMap) run() {
	// load NPCs for this map
	// run logic routines
	// as many workers as needed can be launched
	go zm.playerActivity()
	go zm.playerActivity()
}

func (e *mapEntities) newHandle() uint16 {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.counter++

	if _, used := e.players[e.counter]; used {
		return e.newHandle()
	} else if _, used := e.monsters[e.counter]; used {
		return e.newHandle()
	} else {
		return e.counter
	}
}

func (e *mapEntities) addEntity(en entity) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if pen, ok := en.(*player); ok {
		e.players[en.getHandle()] = pen
	} else if men, ok := en.(*monster); ok {
		e.monsters[en.getHandle()] = men
	} else {
		log.Error("unknown entity")
	}
}

func (e *mapEntities) removeEntity(en entity) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if pen, ok := en.(player); ok {
		delete(e.players, pen.getHandle())
	} else if men, ok := en.(monster); ok {
		delete(e.monsters, men.getHandle())
	} else {
		log.Error("unknown entity")
	}
}

// load maps
func loadMaps() []zoneMap {
	reload := viper.GetBool("game_data.reload")
	dbPath := viper.GetString("game_data.database")
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
			handles: mapEntities{
				counter:  0,
				players:  make(map[uint16]*player),
				monsters: make(map[uint16]*monster),
			},
			send: nil,
			recv: nil,
		}
		zoneMaps = append(zoneMaps, zm)
	}

	return zoneMaps
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
