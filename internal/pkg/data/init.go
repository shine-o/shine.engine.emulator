package data

import (
	"os"

	"github.com/google/logger"
)

var (
	log       *logger.Logger
	filesPath = "../../../files"
)

func init() {
	log = logger.Init("maps logger", true, false, os.Stdout)
}
