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
	entities  entities
	send      sendEvents
	recv      recvEvents
}

type entities struct {
	players
	monsters
}

type players struct {
	manager handleManager
	active  map[uint16]*player
	sync.RWMutex
}

type monsters struct {
	manager handleManager
	active  map[uint16]*monster
	sync.RWMutex
}

type handleManager struct {
	min         uint16
	max         uint16
	index       uint16
	maxAttempts uint16
	list        *roaring.Bitmap
}

const playerHandleMin = 8000
const playerHandleMax = 12000
const playerAttemptsMax = 50

const monsterHandleMin = 17000
const monsterHandleMax = 27000
const monsterAttemptsMax = 50

func (zm *zoneMap) run() {
	// load NPCs for this map
	// run logic routines
	// as many workers as needed can be launched
	go zm.mapHandles()
	go zm.playerActivity()
	go zm.playerQueries()
	go zm.monsterQueries()
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
			entities: entities{
				players: players{
					manager: handleManager{
						min:         playerHandleMin,
						max:         playerHandleMax,
						index:       playerHandleMin,
						maxAttempts: playerAttemptsMax,
						list:        roaring.NewBitmap(),
					},
				},
				monsters: monsters{
					manager: handleManager{
						min:         monsterHandleMin,
						max:         monsterHandleMax,
						index:       monsterHandleMin,
						maxAttempts: monsterAttemptsMax,
						list:        roaring.NewBitmap(),
					},
				},
			},
			send: make(sendEvents),
			recv: make(recvEvents),
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

func (h *handleManager) newHandle() error {
	var attempts uint16 = 0
	for h.index < h.max {
		if attempts == h.maxAttempts {
			return fmt.Errorf("\nmaximum number of attempts reached, no handle is available")
		}
		h.index++
		if h.index == h.max {
			h.index = h.min
		}
		if h.list.Contains(uint32(h.index)) {
			fmt.Printf("\nhandle %v already in use", h.index)
			attempts++
			continue
		} else {
			h.list.Add(uint32(h.index))
			return nil
		}
	}
	return nil
}

func (h *handleManager) removeHandle(index uint16) error {
	if h.list.Contains(uint32(index)) {
		h.list.Remove(uint32(index))
		return nil
	} else {
		return fmt.Errorf("\nhandle %v not in use, removal is not necessary", index)
	}
}
