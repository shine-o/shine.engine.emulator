package maps

import (
	"bytes"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.core/game-data/blocks"
	"github.com/shine-o/shine.engine.core/game-data/shn"
	"github.com/shine-o/shine.engine.core/game-data/utils"
	"github.com/shine-o/shine.engine.core/game-data/world"
	"github.com/shine-o/shine.engine.core/structs"
	bolt "go.etcd.io/bbolt"
	"golang.org/x/image/bmp"
	"image"
	"io/ioutil"
	"math"
	"os"
)

var log *logger.Logger

func init() {
	log = logger.Init("maps logger", true, false, ioutil.Discard)
	log.Info("maps logger init()")
}

type MapError struct {
	Message string
	Code	string
}

func (e *MapError) Error() string {
	return e.Message
}

type Map struct {
	ID uint16
}

type SectorGrid map[int]map[int]*Sector

type Sector struct {
	Row             int
	Column          int
	Image           image.Image
	AdjacentSectors []*Sector
	WalkableX       *roaring.Bitmap
	WalkableY       *roaring.Bitmap
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

	for _, file := range mapFiles {
		mapsPath, err := utils.ValidPath(shineFolder + "/world/" + file)
		maps, err := world.LoadMaps(mapsPath)
		if err != nil {
			return err
		}

		for _, m := range maps {
			if attr, ok := attributes[m.Attributes.ID]; ok {
				m.Attributes = attr
			} else {
				return fmt.Errorf("unkown map attribute entry with ID %v", m.Attributes.ID)
			}

			for _, row := range mapInfo.Rows {
				if row.MapName.Name == m.Attributes.MapInfoIndex {
					m.Info = row
					break
				}
			}

			if m.Info == (shn.MapInfo{}) {
					log.Errorf("no MapInfo.shn entry found for normal map entry with ID %v, ignoring map", m.Attributes.ID)
				continue
			}

			// load shbd


			err = db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte("maps"))
				data, err := structs.Pack(&m)
				if err != nil {
					return err
				}
				err = b.Put([]byte(m.Attributes.MapInfoIndex), data)
				if err != nil {
					return err
				}
				return nil
			})
		}
	}
	return nil
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

func CreateSectorGrid(s *blocks.SHBD, sectorX, sectorY int, img *image.NRGBA) (SectorGrid, error) {
	sX := s.X * 8 / sectorX
	sY := s.Y / sectorY

	grid := make(SectorGrid)
	for i := 1; i <= sY; i++ {
		grid[i] = make(map[int]*Sector)
	columnsLoop:
		for j := 1; j <= sX; j++ {
			subImg := img.SubImage(image.Rectangle{
				Min: image.Point{
					X: sectorX * (j - 1),
					Y: sectorY * (i - 1),
				},
				Max: image.Point{
					X: sectorX * j,
					Y: sectorY * i,
				},
			}).(*image.NRGBA)

			ss := blocks.ImageToSHBD(subImg)

			walkableX, walkableY, err := WalkingPositions(&ss)

			if err != nil {
				return grid, err
			}

			// ignore black only sectors
			if walkableX.IsEmpty() && walkableY.IsEmpty() {
				continue columnsLoop
			}

			grid[i][j] = &Sector{
				Row:       i,
				Column:    j,
				Image:     subImg,
				WalkableX: walkableX,
				WalkableY: walkableY,
			}
		}
	}
	return grid, nil
}

func AdjacentSectors(row, column int, grid map[int]map[int]*Sector) []*Sector {
	var adjacentSectors []*Sector

	// top row
	if topLeft, ok := grid[row-1][column-1]; ok {
		adjacentSectors = append(adjacentSectors, topLeft)
	}

	if top, ok := grid[row-1][column]; ok {
		adjacentSectors = append(adjacentSectors, top)
	}

	if topRight, ok := grid[row-1][column+1]; ok {
		adjacentSectors = append(adjacentSectors, topRight)
	}

	// middle row
	if middleLeft, ok := grid[row][column-1]; ok {
		adjacentSectors = append(adjacentSectors, middleLeft)
	}

	// no need for the middle since its the one we're adding adjacent sectors to
	//if middle, ok := grid[row][column]; ok {
	//	adjacentSectors = append(adjacentSectors, middle)
	//}

	if middleRight, ok := grid[row][column+1]; ok {
		adjacentSectors = append(adjacentSectors, middleRight)
	}

	// bottom row
	if bottomLeft, ok := grid[row+1][column-1]; ok {
		adjacentSectors = append(adjacentSectors, bottomLeft)
	}

	if bottom, ok := grid[row+1][column]; ok {
		adjacentSectors = append(adjacentSectors, bottom)
	}

	if bottomRight, ok := grid[row+1][column+1]; ok {
		adjacentSectors = append(adjacentSectors, bottomRight)
	}

	return adjacentSectors
}

// AdjacentSectorsMesh for each sector append adjacent sectors
func (sg *SectorGrid) AdjacentSectorsMesh() {
	for i, row := range *sg {
		for j, column := range row {
			as := AdjacentSectors(i, j, *sg)
			column.AdjacentSectors = as
		}
	}
}

// SaveGridToBMPFiles for debugging purposes
func SaveGridToBMPFiles(grid *SectorGrid, fileName string) error {
	for i, row := range *grid {
		for j, column := range row {
			out, err := os.OpenFile(fileName+fmt.Sprintf("_sector_%v-%v_", i, j)+".bmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
			if err != nil {
				return err
			}
			err = bmp.Encode(out, column.Image)
			if err != nil {
				return err
			}
			out.Close()
		}
	}
	return nil
}
