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
	// check if table exists
	if _, err := session.fs.Fetch(request.Table); err != nil {
		return Response{Data: nil, Error: err}
	}
	// update json, update cache
	data, err := session.fs.DeleteJSON(request.Table, request.Key)
	if err == nil {
		session.cache.data[request.Table] = data
	}
	return Response{Data: "deleted key in " + request.Table, Error: err}
}

// parses entry payload and updates the corresponding table
func (EntryItem) Add(request Request) Response {
	log.Printf("adding value for %s in table %s", request.Key, request.Table)
	// if it doesn't exist, create it. still errors? return error.
	if _, err := session.fs.Fetch(request.Table); err != nil {
		if session.config.autoTable {
			err = session.fs.CreateJSON(request.Table, nil)
		}
		if err != nil {
			return Response{Data: nil, Error: err}
		}
	}
	// update json, and update cache
	data, err := session.fs.UpdateJSON(request.Table, request.Key, request.Value, request.Mode)
	if err == nil {
		session.cache.data[request.Table] = data
	}
	return Response{Data: "added key in " + request.Table, Error: err}
}

// query for an entry in a table. Returns query result (and updates cache).
func (EntryItem) Get(request Request) Response {
	log.Printf("querying table %s", request.Table)
	// read from cache
	if _, ok := session.cache.data[request.Table]; ok {
		log.Printf("reading %s from cache", request.Table)
		data := session.cache.data[request.Table]
		result, err := findKey(data, request.Key)
		return Response{Data: result, Error: err}
	}
	// read from fs
	data, err := session.fs.ReadJSON(request.Table)
	if err != nil {
		return Response{Data: nil, Error: err}
	}
	session.cache.data[request.Table] = data
	result, err := findKey(data, request.Key)
	return Response{Data: result, Error: err}
}
