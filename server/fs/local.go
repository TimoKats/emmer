package server

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type LocalFS struct {
	Folder string
}

// takes filename and returns full path + extension to json file
func (io LocalFS) getPath(filename string) string {
	return filepath.Join(io.Folder, filename) + ".json"
}

// creates empty (or prefilled) JSON file at path
func (io LocalFS) CreateJSON(filename string, value any) error {
	path := io.getPath(filename)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck
	bytes, err := json.Marshal(value)
	if err != nil || value == nil {
		_, err = f.WriteString("{}")
		return err
	}
	_, err = f.Write(bytes)
	return err
}

// reads JSON file into map[string]any variable
func (io LocalFS) ReadJSON(filename string) (map[string]any, error) {
	// get raw data
	data := make(map[string]any)
	path := io.getPath(filename)
	file, err := os.ReadFile(path)
	if err != nil {
		return data, errors.New("table " + filename + " not found")
	}
	// put raw data into map object
	err = json.Unmarshal(file, &data)
	return data, err
}

// reads JSON file, updates key/value pair, writes to fs
func (io LocalFS) UpdateJSON(filename string, key []string, value any, mode string) error {
	// get json data
	path := io.getPath(filename)
	data, err := io.ReadJSON(filename)
	if err != nil {
		return err
	}
	// update json data
	err = insertNested(data, key, value, mode)
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
func (io LocalFS) DeleteJSON(filename string, key []string) error {
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
func (io LocalFS) List() ([]string, error) {
	// get all files in io folder
	files, err := os.ReadDir(io.Folder)
	result := []string{}
	if err != nil {
		return result, err
	}
	// iterate over json files
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			result = append(result, f.Name())
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
	if folder == "" {
		var baseFolder string
		if runtime.GOOS == "windows" {
			// Use %AppData% on Windows
			baseFolder = os.Getenv("AppData")
		} else {
			// Use XDG_DATA_HOME on linux
			baseFolder = os.Getenv("XDG_DATA_HOME")
		}
		if folder == "" {
			// if nothing is found, just use home
			baseFolder = os.Getenv("HOME")
		}
		folder = filepath.Join(baseFolder, "emmer")
	}
	// create selected folder if it doesn't exist
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		log.Printf("created folder: %s", folder)
		if err := os.Mkdir(folder, 0755); err != nil {
			log.Panic("can't setup emmer folder")
		}
	}
	log.Printf("selected local fs in: %s", folder)
	return &LocalFS{
		Folder: folder,
	}
}
