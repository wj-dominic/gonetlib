package util

import (
	"path/filepath"
	"strings"
)

func GetFileNameWithoutExt(filePath string) string {
	filename := filepath.Base(filePath)
	return strings.TrimSuffix(filename, filepath.Ext(filename))
}
