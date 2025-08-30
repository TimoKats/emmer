package server

import (
	"encoding/json"
	"errors"
)

type EntryData struct{}

// used to query on multi-keys. E.g. [1,2,3] returns map[1,2,3]
func (query *QueryPayload) filterEntry(data map[string]any) (any, error) {
	if len(query.Key) == 0 {
		return data, nil
	}
	current := data
	for _, step := range query.Key {
		match, ok := current[step].(map[string]any)
		if !ok {
			if _, ok := current[step]; !ok {
				return make(map[string]any), errors.New(step + " not found")
			}
			return current[step], nil
		}
		current = match
	}
	return current, nil

}

// fetches path for table name, then removes key from JSON.
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

// query for an entry in a table. Returns query result.
func (EntryData) Query(payload []byte) (Response, error) {
	var response Response
	var query QueryPayload
	if err := json.Unmarshal(payload, &query); err != nil {
		return Response{}, err
	}
	// fetching table contents
	if _, err := fs.Fetch(query.TableName); err != nil {
		return response, err
	}
	// filter contents on query
	data, err := fs.ReadJSON(query.TableName)
	if err == nil {
		response.Result, err = query.filterEntry(data)
	}
	return response, err
}
