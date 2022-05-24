package logger

import (
	"errors"
	"io/fs"
	"os"
)

func isExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}
