package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LocalFS struct {
	Folder string
}

// takes filename and returns full path + extension to json file
func (io LocalFS) getPath(filename string) string {
	return filepath.Join(io.Folder, filename) + ".json"
}

// creates empty JSON file at path
func (io LocalFS) CreateJSON(filename string) error {
	path := io.getPath(filename)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck
	_, err = f.WriteString("{}")
	return err
}

// reads JSON file into map[string]any variable
func (io LocalFS) ReadJSON(filename string) (map[string]any, error) {
	path := io.getPath(filename)
	var data map[string]any
	file, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(file, &data)
	return data, err
}

// reads JSON file, updates key/value pair, writes to fs
func (io LocalFS) UpdateJSON(filename string, key []string, value any) error {
	// get json data
	path := io.getPath(filename)
	data, err := io.ReadJSON(filename)
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

// removes key from json file, writes to fs
func (io LocalFS) DeleteJson(filename string, key []string) error {
	// get json data
	path := io.getPath(filename)
	data, err := io.ReadJSON(filename)
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

// gets path based on table/file name. Returns error if not found
func (io LocalFS) Fetch(filename string) (string, error) {
	path := io.getPath(filename)
	if _, err := os.Stat(path); err != nil {
		return path, errors.New("table " + filename + " not found")
	}
	return path, nil
}

// removes entire JSON file
func (io LocalFS) DeleteFile(filename string) error {
	path := io.getPath(filename)
	return os.Remove(path)
}

// list json files in io folder, along with some statistics
func (io LocalFS) List() (map[string]any, error) {
	// get all files in io folder
	files, err := os.ReadDir(io.Folder)
	result := make(map[string]any)
	if err != nil {
		return result, err
	}
	// iterate over json files
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			info, err := f.Info()
			if err != nil {
				log.Printf("skipping %s due to read error", f.Name())
				continue
			}
			result[f.Name()] = map[string]any{
				"size":     fmt.Sprintf("%.2f KB", float64(info.Size())/1024),
				"last mod": info.ModTime().Format(time.RFC822),
			}
		}
	}
	return result, nil
}

// basic info function. Used for logging
func (io LocalFS) Info() string {
	return "local fs with root dir: " + io.Folder
}

// creates new localFS instance with settings applied
func SetupLocal() *LocalFS {
	folder := os.Getenv("EM_FOLDER")
	// default value is ~/.emmer
	if folder == "" {
		dirname, err := os.UserHomeDir()
		if err != nil {
			log.Panic("can't setup emmer folder")
		}
		folder = dirname + "/.emmer"
	}
	// create selected folder if it doesn't exist
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		log.Printf("created folder: %s", folder)
		if err := os.Mkdir(folder, 0755); err != nil {
			log.Panic("can't setup emmer folder")
		}
	}
	return &LocalFS{
		Folder: folder,
	}
}
