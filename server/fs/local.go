package server

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type LocalFS struct {
	Folder string
}

// takes filename and returns full path + extension to json file
func (io LocalFS) getPath(filename string) string {
	return filepath.Join(io.Folder, filename) + ".json"
}

// selects folder based on OS and env variables of the user
func selectFolder() string {
	// user select
	if folder := os.Getenv("EM_FOLDER"); folder != "" {
		return folder
	}
	// defaults
	if runtime.GOOS == "windows" {
		// Use %AppData% on Windows
		return filepath.Join(os.Getenv("AppData"), "emmer")
	} else {
		// Use XDG_DATA_HOME on linux (if exists)
		xdgData := os.Getenv("XDG_DATA_HOME")
		if xdgData == "" {
			return filepath.Join(os.Getenv("HOME"), ".local", "share", "emmer")
		}
		return filepath.Join(xdgData, "emmer")
	}
}

// creates empty (or prefilled) JSON file at path
func (io LocalFS) Put(filename string, value any) error {
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
func (io LocalFS) Get(filename string) (map[string]any, error) {
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

// removes entire JSON file
func (io LocalFS) Del(filename string) error {
	path := io.getPath(filename)
	return os.Remove(path)
}

// list json files in io folder
func (io LocalFS) Ls() ([]string, error) {
	files, err := os.ReadDir(io.Folder)
	result := []string{}
	if err != nil {
		return result, err
	}
	// iterate over json files
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			filename := strings.TrimSuffix(f.Name(), ".json")
			result = append(result, filename)
		}
	}
	return result, nil
}

// creates new localFS instance with settings applied
func SetupLocal() *LocalFS {
	folder := selectFolder()
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
