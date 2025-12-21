package server

import (
	"log/slog"
)

type EntryItem struct{}

// fetches path for table name, then removes key from JSON.
func (EntryItem) Del(request Request) Response {
	slog.Debug("delete:", "key", request.Key, "table", request.Table)
	// read file from cache/fs
	data, err := read(request.Table, request.Mode)
	if err != nil {
		return Response{Data: nil, Error: err}
	}
	// update contents, and write to cache/fs
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
		return Response{Data: nil, Error: err}
	}
	// update json, and update cache
	err = insert(data, request.Key, request.Value, request.Mode)
	if err != nil {
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
