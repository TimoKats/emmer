package server

// example...

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// returns default config path (~/.emmer/) or value in EMPATH
func EmmerPath() string {
	customPath := os.Getenv("EMPATH")
	dirname, _ := os.UserHomeDir()
	if len(customPath) > 0 {
		return customPath
	}
	return dirname + ".emmer/"
}

// checks if filename exists (exclusive of extension) (to fs export)
func GetFileByName(path string, fileName string) (string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return "", err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			if strings.TrimSuffix(name, filepath.Ext(name)) == fileName {
				return filepath.Join(path, name), nil
			}
		}
	}
	return "", errors.New("file '" + fileName + "' not found")
}

func ReadLocal(path string) ([]byte, error) {
	return os.ReadFile(path)
}
