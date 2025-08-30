package server

import (
	"encoding/json"
	"errors"
)

type EntryData struct{}

// Fetches path for table name, then removes key from JSON.
func (EntryData) Del(payload []byte) error {
	var entry EntryPayload
	if err := json.Unmarshal(payload, &entry); err != nil {
		return err
	}
	if len(entry.TableName) == 0 && len(entry.Key) == 0 {
		return errors.New("no table/key supplied")
	}
	path, err := fs.Fetch(entry.TableName)
	if err != nil {
		return err
	}
	return fs.DeleteJson(path, entry.Key)
}

// parses entry payload and updates the corresponding table
func (EntryData) Add(payload []byte) error {
	var entry EntryPayload
	if err := json.Unmarshal(payload, &entry); err != nil {
		return err
	}
	if len(entry.TableName) == 0 && len(entry.Key) == 0 {
		return errors.New("no table/key supplied")
	}
	if _, err := fs.Fetch(entry.TableName); err != nil {
		return err
	}
	return fs.UpdateJSON(entry.TableName, entry.Key, entry.Value)
}

func (EntryData) Query(payload []byte) error {
	return nil
}
