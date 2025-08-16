package server

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
)

type LocalIO struct {
	Folder string
}

func (io LocalIO) CreateJSON(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("{}")
	return err
}

func (io LocalIO) ReadJSON(path string) (map[string]any, error) {
	var data map[string]any
	file, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(file, &data)
	return data, err
}

func (io LocalIO) UpdateJSON(path string, key string, value any) error {
	data, err := io.ReadJSON(path)
	if err != nil {
		return err
	}
	data[key] = value
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

func (io LocalIO) DelJSON(path string, key string) error {
	data, err := io.ReadJSON(path)
	if err != nil {
		return err
	}
	delete(data, key)
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

func (io LocalIO) Fetch(filename string) (string, error) {
	log.Println(filename, io.Folder)
	path := filepath.Join(io.Folder, filename) + ".json"
	if _, err := os.Stat(path); err != nil {
		return path, errors.New("table not found")
	}
	return path, nil // NOTE: it feel strange
}

func (io LocalIO) DeleteFile(filename string) error {
	path := filepath.Join(io.Folder, filename)
	return os.Remove(path)
}

func (io LocalIO) Info() string {
	return "local fs with root dir: " + io.Folder
}
