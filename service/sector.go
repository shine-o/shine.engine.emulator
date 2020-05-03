package service

import (
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/shine-o/shine.engine.core/game-data/blocks"
	"github.com/shine-o/shine.engine.core/game/entities"
	"golang.org/x/image/bmp"
	"image"
	"os"
)

type sector struct {
	row             int
	column          int
	image           image.Image
	adjacentSectors []*sector

	// broadcast data to all entities within the sector, usually data from an adjacent sector

}

// Sectors are overkill.
func (s *sector) run() {
	// launch logic routines (movement, combat, appearance)
	// each nc handler takes the sector the player is at from
}



type sectorGrid map[int]map[int]*sector

func createSectorGrid(s *blocks.SHBD, sectorX, sectorY int, img *image.NRGBA) (sectorGrid, error) {
	sX := s.X * 8 / sectorX
	sY := s.Y / sectorY

	grid := make(sectorGrid)

	for i := 1; i <= sY; i++ {
		grid[i] = make(map[int]*sector)
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

			walkableX, walkableY, err := walkingPositions(&ss)

			if err != nil {
				return grid, err
			}

			// ignore black only sectors
			if walkableX.IsEmpty() && walkableY.IsEmpty() {
				continue columnsLoop
			}

			grid[i][j] = &sector{
				row:       i,
				column:    j,
				image:     subImg,
				walkableX: walkableX,
				walkableY: walkableY,
			}
		}
	}
	return grid, nil
}

func adjacentSectors(row, column int, grid map[int]map[int]*sector) []*sector {
	var adjacentSectors []*sector

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

// for each sector append adjacent sectors
func (sg *sectorGrid) adjacentSectorsMesh() {
	for i, row := range *sg {
		for j, column := range row {
			as := adjacentSectors(i, j, *sg)
			column.adjacentSectors = as
		}
	}
}

// SaveGridToBMPFiles for debugging purposes
func SaveGridToBMPFiles(grid *sectorGrid, fileName string) error {
	for i, row := range *grid {
		for j, column := range row {
			out, err := os.OpenFile(fileName+fmt.Sprintf("_sector_%v-%v_", i, j)+".bmp", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
			if err != nil {
				return err
			}
			err = bmp.Encode(out, column.image)
			if err != nil {
				return err
			}
			out.Close()
		}
	}
	return nil
}
