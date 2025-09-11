package server

import (
	"errors"
	"log"
)

type EntryItem struct{}

// used to query on multi-keys. E.g. [1,2,3] returns map[1,2,3] > value
func (request *Request) filterEntry(data map[string]any) (any, error) {
	if len(request.Key) == 0 || request.Key[0] == "" {
		return data, nil
	}
	current := data
	for _, step := range request.Key {
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
func (EntryItem) Del(request Request) Response {
	log.Printf("deleting key %s in %v", request.Key, request.Table)
	if _, err := config.fs.Fetch(request.Table); err != nil {
		return Response{Data: nil, Error: err}
	}
	err := config.fs.DeleteJson(request.Table, request.Key)
	return Response{Data: "deleted key in " + request.Table, Error: err}
}

// parses entry payload and updates the corresponding table
func (EntryItem) Add(request Request) Response {
	log.Printf("adding value for %s in table %s", request.Key, request.Table)
	if _, err := config.fs.Fetch(request.Table); err != nil {
		// if it doesn't exist, create it. still errors? return error.
		if config.autoTable {
			err = config.fs.CreateJSON(request.Table)
		}
		if err != nil {
			return Response{Data: nil, Error: err}
		}
	}
	err := config.fs.UpdateJSON(request.Table, request.Key, request.Value, request.Mode)
	return Response{Data: "added key in " + request.Table, Error: err}
}

// query for an entry in a table. Returns query result.
func (EntryItem) Query(request Request) Response {
	log.Printf("query-ing table %s", request.Table)
	if _, err := config.fs.Fetch(request.Table); err != nil {
		return Response{Data: nil, Error: err}
	}
	// filter contents on query
	data, err := config.fs.ReadJSON(request.Table)
	if err != nil {
		return Response{Data: nil, Error: err}
	}
	result, err := request.filterEntry(data)
	return Response{Data: result, Error: err}
}
