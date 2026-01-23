package server

import (
	"log/slog"
	"strings"
)

type EntryItem struct{}

// fetches path for table name, then removes key from JSON.
func (EntryItem) Del(request Request) Response {
	slog.Debug("delete:", "key", request.Key, "table", request.Table)
	data, err := read(request.Table, request.Mode)
	if err != nil {
		return Response{Data: nil, Error: err}
	}
	if err = pop(data, request.Key); err != nil {
		return Response{Data: nil, Error: err}
	}
	if err = write(request.Table, data); err != nil {
		return Response{Data: nil, Error: err}
	}
	return Response{Data: "deleted key in " + request.Table, Error: err}
}

// parses entry payload and updates the corresponding table
func (EntryItem) Add(request Request) Response {
	slog.Debug("add:", "key", request.Key, "table", request.Table)
	// if it doesn't exist, create it. still errors? return error.
	data, err := read(request.Table, request.Mode)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			slog.Warn("autocreate table", "name", request.Table)
			if len(session.cache.tables) > 0 {
				session.cache.tables = append(session.cache.tables, request.Table)
			}
			err = write(request.Table, nil)
			data = make(map[string]any)
		}
		if err != nil {
			return Response{Data: nil, Error: err}
		}
	}
	// update json, and update cache
	if err = insert(data, request.Key, request.Value, request.Mode); err != nil {
		return Response{Data: nil, Error: err}
	}
	if err = write(request.Table, data); err != nil {
		return Response{Data: nil, Error: err}
	}
	return Response{Data: "added key in " + request.Table, Error: err}
}

// query for an entry in a table. Returns query result (and updates cache).
func (EntryItem) Get(request Request) Response {
	slog.Debug("query:", "table", request.Table)
	data, err := read(request.Table, request.Mode)
	if err != nil {
		return Response{Data: nil, Error: err}
	}
	result, err := query(data, request.Key)
	return Response{Data: result, Error: err}
}
