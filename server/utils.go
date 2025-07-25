package server

import (
	"encoding/csv"
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

func createCSV(path string, columns []string, sep rune) error {
	data := [][]string{columns}
	f, _ := os.Create(path)
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = sep
	return w.WriteAll(data)
}

func appendCSVRow(path string, values []string, sep rune) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	w.Comma = sep
	err = w.Write(values)
	if err != nil {
		return err
	}
	w.Flush()
	return w.Error()
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
