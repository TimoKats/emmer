package server

import (
	"errors"
)

// Fetches path on table name, then updates JSON.
func (entry *EntryPayload) add() error {
	if len(entry.TableName) == 0 && len(entry.Key) == 0 {
		return errors.New("no table/key supplied")
	}
	path, err := fs.Fetch(entry.TableName)
	if err != nil {
		return err
	}
	return fs.UpdateJSON(path, entry.Key, entry.Value)
}

// Fetches path for table name, then removes key from JSON.
func (entry *EntryPayload) del() error {
	if len(entry.TableName) == 0 && len(entry.Key) == 0 {
		return errors.New("no table/key supplied")
	}
	path, err := fs.Fetch(entry.TableName)
	if err != nil {
		return err
	}
	return fs.DeleteJson(path, entry.Key)
}
