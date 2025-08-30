package server

import (
	"encoding/json"
)

type TableData struct{}

// Calles by engine. Check if table exists, if yes, remove table.
func (TableData) Del(payload []byte) error {
	var table TablePayload
	if err := json.Unmarshal(payload, &table); err != nil {
		return err
	}
	return fs.DeleteFile(table.Name)
}

// parses payload of table, and creates it if it doesn't exist.
func (TableData) Add(payload []byte) error {
	var table TablePayload
	if err := json.Unmarshal(payload, &table); err != nil {
		return err
	}
	return fs.CreateJSON(table.Name)
}

func (TableData) Query(payload []byte) error {
	return nil
}
