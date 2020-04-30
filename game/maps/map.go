package maps

import (
	"bytes"
	"github.com/RoaringBitmap/roaring"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.core/game-data/blocks"
	"github.com/shine-o/shine.engine.core/game-data/shn"
	"github.com/shine-o/shine.engine.core/game-data/utils"
	"github.com/shine-o/shine.engine.core/game-data/world"
	"github.com/shine-o/shine.engine.core/structs"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"math"
	"sync"
)

var log *logger.Logger

func init() {
	log = logger.Init("maps logger", true, false, ioutil.Discard)
	log.Info("maps logger init()")
}

type MapError struct {
	Message string
	Code    string
}

func (e *MapError) Error() string {
	return e.Message
}

type MapData struct {
	ID         int `struct:"int32"`
	Attributes world.MapAttributes
	Info       shn.MapInfo
	SHBD       blocks.SHBD
}

type mapData struct {
	attributes  map[int]world.MapAttributes
	mapInfo     *shn.ShineMapInfo
	shineFolder string
}

func LoadMapData(shineFolder string, db *bolt.DB) error {
	var attributes map[int]world.MapAttributes
	var mapInfo shn.ShineMapInfo

	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("maps"))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	mapFiles := []string{"NormalMaps.txt", "DungeonMaps.txt", "KingdomQuestMaps.txt", "GuildTournamentMaps.txt"}

	mapAttributesPath, err := utils.ValidPath(shineFolder + "/world/" + "MapAttributes.txt")

	attributes, err = world.LoadMapAttributes(mapAttributesPath)
	if err != nil {
		return err
	}

	mapInfoPath, err := utils.ValidPath(shineFolder + "/shn/client/" + "MapInfo.shn")
	if err != nil {
		return err
	}

	err = shn.Load(mapInfoPath, &mapInfo)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, file := range mapFiles {
		var maps []world.Map
		mapsPath, err := utils.ValidPath(shineFolder + "/world/" + file)
		maps, err = world.LoadMaps(mapsPath)
		if err != nil {
			return err
		}
		for _, m := range maps {
			wg.Add(1)
			md := mapData{
				shineFolder: shineFolder,
				attributes:  attributes,
				mapInfo:     &mapInfo,
			}
			go md.PersistMap(&wg, m, db)
		}
	}
	wg.Wait()
	return nil
}

func (md *mapData) PersistMap(wg *sync.WaitGroup, m world.Map, db *bolt.DB) {
	defer wg.Done()
	var lm MapData
	if attr, ok := md.attributes[m.MapAttributeID]; ok {
		lm.Attributes = attr
	} else {
		log.Errorf("unkown map attribute entry with ID %v", m.MapAttributeID)
		return
	}

	for _, row := range md.mapInfo.Rows {
		if row.MapName.Name == lm.Attributes.MapInfoIndex {
			lm.Info = row
		}
	}

	if lm.Info == (shn.MapInfo{}) {
		log.Errorf("no MapInfo.shn entry found for normal map entry with ID %v, ignoring map", lm.Attributes.ID)
		return
	}

	// load shbd
	var s *blocks.SHBD

	shbdPath, err := utils.ValidPath(md.shineFolder + "/blocks/" + lm.Info.MapName.Name + ".shbd")
	if err != nil {
		log.Errorf("shbd file found for normal map entry with ID %v, ignoring map", lm.Attributes.ID)
	}

	s, err = blocks.LoadSHBDFile(shbdPath)
	if err != nil {
		log.Errorf("failed to load shbd file for map entry with ID %v, ignoring map %v", lm.Attributes.ID, err)
	}

	lm.SHBD = *s

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("maps"))
		data, err := structs.Pack(&lm)
		if err != nil {
			return err
		}
		err = b.Put([]byte(lm.Attributes.MapInfoIndex), data)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		log.Errorf("failed to persist map with id %v %v", lm.Attributes.ID, err)
	}
}

// CanWalk translates in game coordinates to SHBD coordinates
func CanWalk(x, y *roaring.Bitmap, rX, rY int) bool {
	if x.ContainsInt(rX) && y.ContainsInt(rY) {
		return true
	}
	return false
}

// WalkingPositions creates two X,Y roaring bitmaps with walkable coordinates
func WalkingPositions(s *blocks.SHBD) (*roaring.Bitmap, *roaring.Bitmap, error) {
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
