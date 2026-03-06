package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Json struct {
	Folder string
}

// takes filename and returns full path + extension to json file
func (fs Json) getPath(filename string) string {
	return filepath.Join(fs.Folder, filename) + ".json"
}

// creates empty (or prefilled) JSON file at path
func (fs Json) Put(filename string, value any) error {
	path := fs.getPath(filename)
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
func (fs Json) Get(filename string) (any, error) {
	// get raw data
	mapping := make(map[string]any)
	list := []any{}
	path := fs.getPath(filename)
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.New("table " + filename + " not found")
	}
	// put raw data into map object
	if err := json.Unmarshal(file, &mapping); err == nil {
		return mapping, nil
	} else if err := json.Unmarshal(file, &list); err == nil {
		return list, nil
	}
	return nil, errors.New("error reading file, is it json?")
}

// removes entire JSON file
func (fs Json) Del(filename string) error {
	path := fs.getPath(filename)
	return os.Remove(path)
}

// list json files in fs folder
func (fs Json) Ls() ([]string, error) {
	files, err := os.ReadDir(fs.Folder)
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
func SetupJSON() *Json {
	folder := selectFolder()
	// create selected folder if it doesn't exist
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		slog.Debug("created folder", "folder", folder)
		if err := os.Mkdir(folder, 0755); err != nil {
			slog.Error("can't create emmer folder:", "path", folder)
			os.Exit(1)
		}
	}
	slog.Info("selected JSON fs:", "folder", folder)
	return &Json{
		Folder: folder,
	}
}
