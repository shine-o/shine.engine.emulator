package maps

import (
	"bytes"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/disintegration/imaging"
	"github.com/shine-o/shine.engine.core/structs"
	"github.com/shine-o/shine.engine.core/structs/blocks"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"math"
	"os"
)

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

func SHBDToImage(s *blocks.SHBD) (*image.NRGBA, error) {
	r := bytes.NewReader(s.Data)

	img := image.NewNRGBA(image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: s.X * 8,
			Y: s.Y,
		},
	})

	for y := s.Y; y != 0; y-- {
		for x := 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return img, err
			}
			for i := 0; i < 8; i++ {
				var (
					rX, rY int
					c      color.Color
				)

				rX = x*8 + i
				rY = y

				if b&byte(math.Pow(2, float64(i))) == 0 {
					c = color.White
				} else {
					c = color.Black
				}
				img.Set(rX, rY, c)
			}
		}
	}
	return img, nil
}

func ImageToSHBD(img *image.NRGBA) blocks.SHBD {
	img = imaging.FlipV(img)

	bounds := img.Bounds()

	var rs = blocks.SHBD{
		X:    bounds.Max.X / 8,
		Y:    bounds.Max.Y,
		Data: make([]byte, 0),
	}

	for y := 0; y < rs.Y; y++ {
		for x := 0; x < rs.X; x++ {
			var sb uint8 = 0

			for i := 0; i < 8; i++ {
				offset := img.PixOffset(x*8+i, y)

				b := img.Pix[offset]

				if b == 0 {
					sb |= 1 << i
				}
			}
			rs.Data = append(rs.Data, sb)
		}
	}
	return rs
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

			ss := ImageToSHBD(subImg)

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

// CanWalk translates in game coordinates to SHBD coordinates
func CanWalk(x, y *roaring.Bitmap, rX, rY int) bool {
	if x.ContainsInt(rX) && y.ContainsInt(rY) {
		return true
	}
	return false
}

// SaveBmpFile for debugging purposes
func SaveBmpFile(img *image.NRGBA, path, fileName string) error {
	out, err := os.OpenFile(path+fileName+".bmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	err = bmp.Encode(out, img)
	if err != nil {
		return err
	}
	out.Close()
	return nil
}

// SaveSHBDFile for debugging purposes
func SaveSHBDFile(s *blocks.SHBD, path, fileName string) error {
	data, err := structs.Pack(s)

	if err != nil {
		return err
	}

	out, err := os.OpenFile(path+fileName+".shbd", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
	if err != nil {
		return err
	}
	out.Write(data)
	out.Close()
	return nil
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
