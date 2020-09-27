package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/RoaringBitmap/roaring"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/blocks"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/game-data/shn"
	"github.com/shine-o/shine.engine.emulator/pkg/structs"
	"math"
)

// used to run with debugger and check the rows are loaded correctly
func main() {
	//shn sample
	//shnTest()
	//packet data sample
	packetDataTest()
	//shbdTest()

}

//
func shnTest()  {
	var mi shn.ShineMobInfoServer
	err := shn.Load("assets/MobInfoServer.shn", &mi)
	if err != nil {
		logger.Error(err)
	}
}

func packetDataTest()  {
	var s structs.NcBatTargetInfoCmd

	hexS := "00ce1eab020000ab0200006400000064000000000000000000000014e807"

	data, err := hex.DecodeString(hexS)

	if err != nil {
		logger.Error(err)
	}

	err = structs.Unpack(data, &s)

	if err != nil {
		logger.Error(err)
	}
}

func shbdTest()  {
	m := "EldPri01"
	var s *blocks.SHBD
	s, err := blocks.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	img, err := blocks.SHBDToImage(s)
	if err != nil {
		logger.Error(err)
	}

	err = blocks.SaveBmpFile(img, "./", m)

	if err != nil {
		logger.Error(err)
	}

	rs := blocks.ImageToSHBD(img)

	blocks.SaveSHBDFile(&rs, "./", m)

	xbm, ybm, err := walkingPositions(s)

	if err != nil {
		logger.Error(err)
	}
	testWalk(xbm, ybm)
}

func canWalk(x, y *roaring.Bitmap, rX, rY uint32) bool {
	if x.ContainsInt(int(rX)) && y.ContainsInt(int(rY)) {
		return true
	}
	return false
}

func testWalk(walkableX *roaring.Bitmap, walkableY *roaring.Bitmap) {
	igX := 5868
	igY := 10462

	rX := (igX * 8.0) / 50.0

	//fmt.Printf("%.6f", float64(rX))
	// 1589 = (x * 8) / 50
	rstX := (rX * 50.0) / 8.0

	fmt.Printf("%.6f", float64(rstX))

	rY := (igY * 8.0) / 50

	if canWalk(walkableX, walkableY, uint32(rX), uint32(rY)) {
		fmt.Printf("\nrX: %v, rY: %v", rX, rY)
		fmt.Printf("\nigX: %v, igY: %v", igX, igY)
	}
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
