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
	case ".csv": // to enum ?
		return fs.AppendCSV(path, entry.Values)
	case ".json":
		return fs.AddJSON(path, entry.Key, entry.Value)
	default:
		return errors.New("file extension not supported")
	}
}

func (entry *EntryPayload) del() error {
	path, err := fs.GetFileByName("data/", entry.TableName)

	if err != nil {
		return err
	}
	switch filepath.Ext(path) {
	case ".csv":
		return errors.New("file extension not implemented")
	case ".json":
		return fs.DelJSON(path, entry.Key)
	default:
		return errors.New("file extension not supported")
	}
}
