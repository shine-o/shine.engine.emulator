package data

import (
	"bytes"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/structs"
	"golang.org/x/image/bmp"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"os"
)

type SHBD struct {
	X    int    `struct:"int32"`
	Y    int    `struct:"int32"`
	Data []byte `struct-size:"X * Y"`
}

func LoadSHBDFile(filesPath string) (*SHBD, error) {
	var s *SHBD

	data, err := ioutil.ReadFile(filesPath)
	if err != nil {
		return s, err
	}

	err = structs.Unpack(data, &s)
	if err != nil {
		return s, err
	}

	return s, nil
}

func SHBDToImage(s *SHBD) (*image.RGBA, error) {
	r := bytes.NewReader(s.Data)

	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: s.X * 8,
			Y: s.Y,
		},
	})

	for y := 0; y < s.Y; y++ {
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

func ImageToSHBD(img *image.RGBA) SHBD {
	//img = imaging.FlipV(img)

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

// SaveBmpFile for debugging purposes
func SaveBmpFile(img *image.RGBA, path, fileName string) error {
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
