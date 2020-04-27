package blocks

import (
	"bytes"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/disintegration/imaging"
	"github.com/shine-o/shine.engine.core/structs"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"math"
	"os"
)

type SHBD struct {
	X    int    `struct:"int32"`
	Y    int    `struct:"int32"`
	Data []byte `struct-size:"X * Y"`
}

type SectorGrid map[int]map[int]Sector

type Sector struct {
	Row       int
	Column    int
	Image     image.Image
	WalkableX * roaring.Bitmap
	WalkableY * roaring.Bitmap
}

func SaveGridToBMPFiles(grid *SectorGrid, fileName string) error {
	for i, row := range *grid {
		for j, column := range row {
			out, err := os.OpenFile(fileName + fmt.Sprintf("_sector_%v-%v_", i, j)+".bmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
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

func CreateSectorGrid(s * SHBD, sectorX int, sectorY int, img *image.NRGBA) (*SectorGrid, error) {
	sX := s.X * 8 / sectorX
	sY := s.Y / sectorY

	grid := make(SectorGrid)
	for i := 1; i <= sY; i++ {
		grid[i] = make(map[int]Sector)
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

			ss := BMPToSHBD(subImg)

			walkableX, walkableY, err := WalkingPositions(&ss)

			if err != nil {
				return &grid, err
			}

			// ignore black only sectors
			if walkableX.IsEmpty() && walkableY.IsEmpty() {
				continue
			}

			grid[i][j] = Sector{
				Row:       i,
				Column:    j,
				Image:     subImg,
				WalkableX: walkableX,
				WalkableY: walkableY,
			}
		}
	}
	return &grid, nil
}

func WalkingPositions(s * SHBD) (*roaring.Bitmap, *roaring.Bitmap,  error){
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
					//log.Infof("rX: %v, rY: %v", rX, rY)
				}
			}
		}
	}

	return walkableX, walkableY, nil
}

func TestWalk(walkableX *roaring.Bitmap, walkableY *roaring.Bitmap) {
	igX := 14987
	igY := 13415
	rX := (igX * 8) / 50
	rY := (igY * 8) / 50

	if CanWalk(walkableX, walkableY, rX, rY) {
		fmt.Printf("\nrX: %v, rY: %v", rX, rY)
		fmt.Printf("\nigX: %v, igY: %v", igX, igY)
	}
}

func CanWalk(x, y *roaring.Bitmap, rX, rY int) bool {
	if x.ContainsInt(rX) && y.ContainsInt(rY) {
		return true
	}
	return false
}


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

func SaveSHBDFile(s *SHBD, path, fileName string) error {
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

func SHBDtoImage(s *SHBD) (*image.NRGBA, error) {
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

func BMPToSHBD(img *image.NRGBA) SHBD {
	img = imaging.FlipV(img)

	bounds := img.Bounds()

	var rs = SHBD{
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

func AdjacentSectors(row int, column int, grid map[int]map[int]Sector) []Sector {
	var adjacentSectors []Sector

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

	if middle, ok := grid[row][column]; ok {
		adjacentSectors = append(adjacentSectors, middle)
	}

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
