package main

import (
	"fmt"
	"github.com/google/logger"
	"github.com/shine-o/shine.engine.emulator/internal/pkg/data"
)

func main() {
	shbdTest()
}

func shbdTest() {
	m := "Rou"
	var s *data.SHBD
	s, err := data.LoadSHBDFile(fmt.Sprintf("C:\\Users\\marbo\\go\\src\\github.com\\shine-o\\shine.engine.emulator\\files\\blocks\\%v.shbd", m))

	if err != nil {
		logger.Error(err)
	}

	img, err := data.SHBDToImage(s)
	if err != nil {
		logger.Error(err)
	}

	err = data.SaveBmpFile(img, "./", m)

	if err != nil {
		logger.Error(err)
	}

	rs := data.ImageToSHBD(img)

	err = data.SaveSHBDFile(&rs, "./", m)

	if err != nil {
		logger.Error(err)
	}
}
