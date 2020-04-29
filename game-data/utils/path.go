package utils

import (
	"os"
	"path/filepath"
)

func ValidPath(path string) (string, error) {
	var absPath string
	absPath, err := filepath.Abs(path)
	if err != nil {
		return absPath, err
	}
	if _, err := os.Stat(path); err == os.ErrNotExist {
		return absPath, err
	}
	return absPath, nil
}

