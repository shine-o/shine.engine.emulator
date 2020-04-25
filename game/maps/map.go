package maps

import "github.com/RoaringBitmap/roaring"

type Sector interface {}
type SectorEvent interface {}

func CanWalk(x, y *roaring.Bitmap, igX, igY uint32) bool {
	rX := (igX * 8) / 50
	rY := (igY * 8) / 50

	if x.Contains(rX) && y.Contains(rY) {
		return true
	}
	return false
}