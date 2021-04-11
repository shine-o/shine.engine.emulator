package data

import (
	"github.com/google/logger"
	"os"
)

var (
	log       *logger.Logger
	filesPath = "../../../files"
)

func init() {
	log = logger.Init("maps logger", true, false, os.Stdout)
}
