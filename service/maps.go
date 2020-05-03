package service

import (
	"bytes"
	"github.com/RoaringBitmap/roaring"
	"github.com/shine-o/shine.engine.core/game-data/blocks"
	"github.com/shine-o/shine.engine.core/game-data/utils"
	"github.com/shine-o/shine.engine.core/game-data/world"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/spf13/viper"
	bolt "go.etcd.io/bbolt"
	"math"
)

type zoneMap struct {
	data      world.MapData
	walkableX *roaring.Bitmap
	walkableY *roaring.Bitmap
	entities  map[uint32]entity
	send    map[uint32]chan <- event
	recv    map[uint32]<-chan event
}

func (zm *zoneMap) addEntity(handle int) {

}

func (zm *zoneMap) removeEntity(handle int) {

}

func (zm *zoneMap) run() {
	// load NPCs for this map
	// run logic routines
	go zm.entityMovement()
}

func (zm *zoneMap) entityMovement() {
	for {
		select {
		case e := <- zm.recv[entityAppeared]:
			log.Info(e)
			// notify all nearby entities about it
			// players will get packet data
			// mobs will check if player is in range for attack
		case e := <- zm.recv[entityDisappeared]:
			log.Info(e)
		case e := <- zm.recv[entityMoved]:
			log.Info(e)
		case e := <- zm.recv[entityStopped]:
			log.Info(e)
		case e := <- zm.recv[entityJumped]:
			log.Info(e)
		}
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
	}

	// zones defined in the service, load them and register the loaded ones to the master
	var zoneMaps []zoneMap
	for _, id := range normalMaps {
		var md world.MapData

		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("maps"))
			data := b.Get([]byte(string(id)))
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
			entities:  make(map[uint32]entity),
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
