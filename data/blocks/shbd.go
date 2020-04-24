package blocks

import (
	"bytes"
	"encoding/binary"
	"github.com/RoaringBitmap/roaring"
	"github.com/shine-o/shine.engine.core/structs/blocks"
	"gopkg.in/restruct.v1"
	"io/ioutil"
	"math"
)

func LoadSHBD(filePath string) (*roaring.Bitmap, *roaring.Bitmap, error) {
	walkableX := roaring.NewBitmap()
	walkableY := roaring.NewBitmap()
	data, err := ioutil.ReadFile(filePath)

	if err != nil {
		return walkableX, walkableY, err
	}
	var s blocks.SHBD

	err = restruct.Unpack(data, binary.LittleEndian, &s)
	if err != nil {
		return walkableX, walkableY, err
	}

	r := bytes.NewReader(s.Data)

	var y, x uint32

	for y = 0; y < s.Y; y++ {
		for x = 0; x < s.X; x++ {
			b, err := r.ReadByte()
			if err != nil {
				return walkableX, walkableY, err
			}
			for i := 0; i < 8; i++ {
				if b & byte(math.Pow(2, float64(i))) == 0 {
					rX := int(x) * 8 + i
					rY := int(y)
					walkableX.Add(uint32(rX))
					walkableY.Add(uint32(rY))
				}
			}
		}
	}
	return walkableX, walkableY, nil
}
