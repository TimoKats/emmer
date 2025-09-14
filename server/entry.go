package server

import (
	"errors"
	"log"
	"strconv"
)

type EntryItem struct{}

// used to query on multi-keys. E.g. [1,2,3] returns map[1,2,3] > value
func findKey(data map[string]any, key []string) (any, error) {
	var current any = data
	if len(key) == 0 || key[0] == "" {
		return data, nil
	}
	for _, step := range key {
		switch typed := current.(type) {
		case map[string]any:
			val, ok := typed[step]
			if !ok {
				return nil, errors.New("key " + step + " not found in map")
			}
			current = val
		case []any:
			index, err := strconv.Atoi(step)
			if err != nil {
				return nil, errors.New("invalid index " + step + " for list")
			}
			if index < 0 || index >= len(typed) {
				return nil, errors.New("index " + step + " out of bounds")
			}
			current = typed[index]
		default:
			return nil, errors.New("cannot descend into type")
		}
	}
	return current, nil
}

// fetches path for table name, then removes key from JSON.
func (EntryItem) Del(request Request) Response {
	log.Printf("deleting key %s in %v", request.Key, request.Table)
	if _, err := config.fs.Fetch(request.Table); err != nil {
		return Response{Data: nil, Error: err}
	}
	err := config.fs.DeleteJSON(request.Table, request.Key)
	return Response{Data: "deleted key in " + request.Table, Error: err}
}

// parses entry payload and updates the corresponding table
func (EntryItem) Add(request Request) Response {
	log.Printf("adding value for %s in table %s", request.Key, request.Table)
	// if it doesn't exist, create it. still errors? return error.
	if _, err := config.fs.Fetch(request.Table); err != nil {
		if config.autoTable {
			err = config.fs.CreateJSON(request.Table, request.Value)
		}
		if err != nil {
			return Response{Data: nil, Error: err}
		}
	}
	// update json file with new values
	err := config.fs.UpdateJSON(request.Table, request.Key, request.Value, request.Mode)
	return Response{Data: "added key in " + request.Table, Error: err}
}

// query for an entry in a table. Returns query result.
func (EntryItem) Get(request Request) Response {
	log.Printf("querying table %s", request.Table)
	// get complete json data
	data, err := config.fs.ReadJSON(request.Table)
	if err != nil {
		return Response{Data: nil, Error: err}
	}
	// filter json data
	result, err := findKey(data, request.Key)
	return Response{Data: result, Error: err}
}
