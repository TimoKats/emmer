package server

import (
	"errors"
	"path/filepath"
)

func (entry *EntryPayload) add() error {
	path, err := fs.GetFileByName("data/", entry.TableName)

	if err != nil {
		return err
	}
	switch filepath.Ext(path) {
	case ".csv":
		return fs.AppendCSV(path, entry.Values)
	case ".json":
		return fs.WriteJSON(path, entry.Key, entry.Value)
	default:
		return errors.New("file extension not supported")
	}
}
