package shinelog

import (
	"fmt"
	"os"

	"github.com/google/logger"
	"github.com/sirupsen/logrus"
)

func NewLogger(name, outputFolder string, level logrus.Level) *logger.Logger {
	if _, err := os.Stat(outputFolder); os.IsNotExist(err) {
		err := os.Mkdir(outputFolder, 0o666)
		if err != nil {
			logger.Fatalf("Failed to create output folder: %v", err)
		}
	}
	path := fmt.Sprintf("%v/%v.log", outputFolder, name)
	lf, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o666)
	if err != nil {
		logger.Fatalf("Failed to create output file: %v", err)
	}

	log := logger.Init(name, true, false, lf)

	return log
}
