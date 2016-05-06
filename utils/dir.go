package utils

import (
	"errors"
	"os"
	"strings"
)

// IsDirectory reports whether the named file is a directory.
func IsDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func IsFile(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

var (
	ErrTplNotFound = errors.New("Template file not found")
)

func GetExistingPathFor(file, currDir string) (string, error) {
	if IsFile(file) {
		//passed a full path
		return file, nil
	}
	localFile := strings.Join([]string{currDir, file}, "/")
	if IsFile(localFile) {
		return localFile, nil
	}
	return "", ErrTplNotFound
}
