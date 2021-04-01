package data

import (
	"github.com/google/logger"
	"os"
)

var (
	log * logger.Logger
	filePath = "../../files"
)

func init() {
	log = logger.Init("maps logger", true, false, os.Stdout)
}