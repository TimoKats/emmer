package server

import (
	"encoding/json"
	"os"
)

func createLocalJSON(path string) error {
	f, err := os.Create(path)
	defer f.Close()
	f.WriteString("{}")
	return err
}

func getLocalJson(path string) (map[string]any, error) {
	var data map[string]any
	file, err := os.ReadFile(path)
	if err != nil {
		return data, err
	}
	err = json.Unmarshal(file, &data)
	return data, err
}

func writeLocalJson(path string, data map[string]any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, bytes, 0644)
}
