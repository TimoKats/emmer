package server

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// returns default config path (~/.emmer/) or value in EMPATH
func configPath() string {
	customPath := os.Getenv("EMPATH")
	dirname, _ := os.UserHomeDir()
	if len(customPath) > 0 {
		return customPath
	}
	return dirname + ".emmer/"
}

func createJSON(path string) error {
	f, err := os.Create(path)
	defer f.Close()
	f.WriteString("{}")
	return err
}

func getFile(path string, fileName string) (string, error) {
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
