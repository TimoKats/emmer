package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// read json file at test location
func readJson(filename string) map[string]any {
	data := make(map[string]any)
	file, err := os.ReadFile(filename)
	if err != nil {
		return data
	}
	json.Unmarshal(file, &data) //nolint:errcheck
	return data
}

// test if two maps are equal (loosly/stringified)
func jsonEqual(a, b any) bool {
	aByte, _ := json.Marshal(a) //nolint:errcheck
	bByte, _ := json.Marshal(b) //nolint:errcheck
	return string(aByte) == string(bByte)
}

// create path to test file (GitHub actions uses EM_FOLDER)
func filename() string {
	var folder string
	if folder = os.Getenv("EM_FOLDER"); folder == "" {
		folder = filepath.Join(os.Getenv("HOME"), ".local", "share", "emmer")
	}
	return filepath.Join(folder, "test.json")
}
