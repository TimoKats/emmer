package server

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type LocalIO struct {
	Folder string
}

// creates empty JSON file at path
func (io LocalIO) CreateJSON(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString("{}")
	return err
}

// Reads JSON file into map[string]any variable
func (io LocalIO) ReadJSON(path string) (map[string]any, error) {
	var data map[string]any
	file, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(file, &data)
	return data, err
}

// Reads JSON file, updates key/value pair, writes to fs.
func (io LocalIO) UpdateJSON(path string, key []string, value any) error {
	// get json data
	data, err := io.ReadJSON(path)
	if err != nil {
		return err
	}
	// update json data
	err = insertNested(data, key, value)
	if err != nil {
		return err
	}
	// write json data to file
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

// Removes key from json file, writes to fs.
func (io LocalIO) DeleteJson(path string, key []string) error {
	// get json data
	data, err := io.ReadJSON(path)
	if err != nil {
		return err
	}
	// update json data
	err = deleteNested(data, key)
	if err != nil {
		return err
	}
	// write json data to file
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}

// Gets path based on table/file name. Returns error if not found.
func (io LocalIO) Fetch(filename string) (string, error) {
	path := filepath.Join(io.Folder, filename) + ".json"
	if _, err := os.Stat(path); err != nil {
		return path, errors.New("table not found")
	}
	return path, nil
}

// Removes entire JSON file.
func (io LocalIO) DeleteFile(filename string) error {
	path := filepath.Join(io.Folder, filename)
	return os.Remove(path)
}

// Basic info function. Used for logging.
func (io LocalIO) Info() string {
	return "local fs with root dir: " + io.Folder
}
