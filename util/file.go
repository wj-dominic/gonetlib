package util

import (
	"os"
	"path/filepath"
	"strings"
)

func GetFileNameWithoutExt(filePath string) string {
	filename := filepath.Base(filePath)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}

func IsExistPath(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}
